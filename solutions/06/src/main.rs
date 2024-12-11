use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::Parser;
use std::time::SystemTime;
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};
use tracing::{span, Level};

mod guard_state;
mod logging;
mod obstacle_map;

use guard_state::coordinate;
use guard_state::coordinate::Coordinate;
use guard_state::direction::Direction;
use guard_state::GuardState;
use obstacle_map::ObstacleMap;

const GUARD_INITIAL_DIRECTION_UP: u8 = b'^';
const GUARD_INITIAL_DIRECTION_RIGHT: u8 = b'>';
const GUARD_INITIAL_DIRECTION_DOWN: u8 = b'v';
const GUARD_INITIAL_DIRECTION_LEFT: u8 = b'<';

/// Program to solve Advent of Code puzzles
#[derive(Parser, Debug)]
#[command(version, about)]
struct CommandLineArgs {
    /// File to read input from
    #[arg(long, default_value = "puzzleInput")]
    input_file: String,

    /// Selected part of puzzle to solve
    #[arg(long, value_parser=clap::value_parser!(u8).range(1..=2))]
    part: u8,
}

fn main() {
    logging::set_logging();
    let args = CommandLineArgs::parse();
    trace!("command line args" = ?args, "parsed command line args");

    let input_file_handle = File::open(args.input_file).expect("could not open input file");
    let input_file_reader = BufReader::new(input_file_handle);

    let start_time = SystemTime::now();
    let computation_result = match args.part {
        1 => part01(input_file_reader),
        2 => part02(input_file_reader),
        _ => unreachable!(), // clap has filtered out all other possibilities.
    }
    .expect("computation did not produce a value");
    let elapsed_time = start_time.elapsed().unwrap();

    info!(
        "computation_result" = computation_result,
        "elapsed_time" = elapsed_time.as_nanos(),
        "computation complete"
    );
}

fn parse_input_to_obstacle_map_and_state(
    input_file_reader: BufReader<File>,
) -> (ObstacleMap, Option<GuardState>) {
    let all_lines: Vec<_> = input_file_reader
        .lines()
        .collect::<Result<Vec<_>, _>>()
        .unwrap();
    let mut obs_map = obstacle_map::new(all_lines[0].len() as i32, all_lines[0].len() as i32);
    let mut guard: Option<GuardState> = None;

    for (y, line) in all_lines.iter().enumerate() {
        debug!("line" = ?line, "read line from input file");
        for (x, r) in line.bytes().enumerate() {
            let c = coordinate::new_from_usize(x, y);
            let _span = span!(Level::DEBUG, "parsed byte", "coordinate"=?c, "byte"= r).entered();
            match r {
                obstacle_map::OBSTACLE_BYTE => {
                    trace!("found obstacle");
                    obs_map.add_obstacle(c);
                }
                GUARD_INITIAL_DIRECTION_UP => {
                    trace!("found guard facing up");
                    guard = Some(guard_state::new(c, Direction::Up));
                }
                GUARD_INITIAL_DIRECTION_RIGHT => {
                    trace!("found guard facing right");
                    guard = Some(guard_state::new(c, Direction::Right));
                }
                GUARD_INITIAL_DIRECTION_DOWN => {
                    trace!("found guard facing down");
                    guard = Some(guard_state::new(c, Direction::Down));
                }
                GUARD_INITIAL_DIRECTION_LEFT => {
                    trace!("found guard facing left");
                    guard = Some(guard_state::new(c, Direction::Left));
                }
                _ => {
                    trace!("found other byte");
                }
            }
        }
    }

    (obs_map, guard)
}

/// Compute the ordered sequence of states required to leave the bounds of the map.
///
/// Returns None if the guard loops.
fn compute_path(obs_map: &ObstacleMap, initial_state: GuardState) -> Option<Vec<GuardState>> {
    let _span = span!(Level::DEBUG, "computing path", "initial_state"=?initial_state,).entered();
    let mut seen_states: HashSet<GuardState> = HashSet::new();
    let mut ordered_states: Vec<GuardState> = Vec::new();
    let mut current_state = initial_state;
    let mut total_steps = 0;

    loop {
        if seen_states.contains(&current_state) {
            return None;
        }
        seen_states.insert(current_state);
        ordered_states.push(current_state);

        let next_state = current_state.step();
        match obs_map.is_obstacle(next_state.get_coordinate()) {
            Some(false) => {
                total_steps += 1;
                current_state = next_state;
                trace!("current state"=?current_state, "total steps"=total_steps, "took step");
            }
            Some(true) => {
                trace!("current state"=?current_state, "total steps"=total_steps, "encountered obstacle");
                current_state = current_state.encounter_obstacle();
            }
            None => {
                trace!("current state"=?current_state,"left map");
                return Some(ordered_states);
            }
        }
    }
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let (obs_map, guard) = match parse_input_to_obstacle_map_and_state(input_file_reader) {
        (obs_map, Some(guard)) => (obs_map, guard),
        (_, None) => {
            error!("failed to initialize guard state from input");
            return None;
        }
    };
    debug!("obstacle map"=?obs_map, "initial state"=?guard, "parsed input");

    match compute_path(&obs_map, guard) {
        Some(ordered_states) => {
            let unique_coordinates = ordered_states
                .into_iter()
                .map(|state| state.get_coordinate())
                .collect::<HashSet<_>>();
            return Some(unique_coordinates.len() as i64);
        }
        None => {
            error!("did not find a valid path for part 01");
            return None;
        }
    }
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    let (mut obs_map, guard) = match parse_input_to_obstacle_map_and_state(input_file_reader) {
        (obs_map, Some(guard)) => (obs_map, guard),
        (_, None) => {
            error!("failed to initialize guard state from input");
            return None;
        }
    };
    debug!("obstacle map"=?obs_map, "initial state"=?guard, "parsed input");

    let original_path = match compute_path(&obs_map, guard) {
        Some(ordered_states) => ordered_states,
        None => {
            error!("did not find a valid path for unobstructed input");
            return None;
        }
    };

    let mut unique_loop_creating_obstruction_locations: HashSet<Coordinate> = HashSet::new();
    let mut previously_visited_coordinates: HashSet<Coordinate> = HashSet::new();
    for state in original_path {
        let next_state_coordinate = state.step().get_coordinate();
        if previously_visited_coordinates.contains(&next_state_coordinate) || !obs_map.in_bounds(next_state_coordinate){
             continue;
        }
        if let Some(false) = obs_map.is_obstacle(next_state_coordinate) {
            obs_map.add_obstacle(next_state_coordinate);
            match compute_path(&obs_map, state) {
                Some(_) => {
                    trace!("initial state"=?state, "inserted obstacle"=?next_state_coordinate, "obstacle did not create loop");
                }
                None => {
                    unique_loop_creating_obstruction_locations.insert(next_state_coordinate);
                    debug!("initial state"=?state, "inserted obstacle"=?next_state_coordinate, "obstacle created loop");
                }
            }
            obs_map.remove_obstacle(next_state_coordinate);
        }
        previously_visited_coordinates.insert(state.get_coordinate());
    }

    return Some(unique_loop_creating_obstruction_locations.len() as i64);
}
