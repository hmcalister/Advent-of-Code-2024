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

/// Compute the number of steps required to leave the bounds of the map.
/// 
/// Returns None if the guard loops.
fn compute_path(obs_map: &ObstacleMap, initial_state: GuardState) -> Option<Vec<Coordinate>> {
    let _span = span!(Level::DEBUG, "computing path", "initial_state"=?initial_state,).entered();
    let mut seen_states: HashSet<GuardState> = HashSet::new();
    let mut seen_coordinates: Vec<Coordinate> = Vec::new();
    let mut current_state = initial_state;
    let mut total_steps = 0;

    loop {
        if seen_states.contains(&current_state) {
            return None
        }
        seen_states.insert(current_state);
        seen_coordinates.push(current_state.get_coordinate());
        
        let next_state = current_state.step();
        match obs_map.is_obstacle(next_state.get_coordinate()) {
            Some(false) => {
                total_steps += 1;
                current_state = next_state;
                debug!("current state"=?current_state, "total steps"=total_steps, "took step");
            },
            Some(true) => {
                debug!("current state"=?current_state, "total steps"=total_steps, "encountered obstacle");
                current_state = current_state.encounter_obstacle();
            },
            None => {
                debug!("current state"=?current_state,"left map");
                return Some(seen_coordinates);
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
        Some(seen_coordinates) => {
            let unique_coordinates = seen_coordinates.into_iter().collect::<HashSet<_>>();
            return Some(unique_coordinates.len() as i64)
        },
        None => {
            error!("did not find a vlid path for part 01");
            return None;
        }
    }

}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
}
