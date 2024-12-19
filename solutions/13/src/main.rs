use std::error;
use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::Parser;
use claw_machine::ClawMachine;
use nalgebra::{Matrix2, Vector2};
use regex::Regex;
use std::time::SystemTime;
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

mod claw_machine;
mod logging;

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

fn parse_input_to_claw_machines(
    mut input_file_reader: BufReader<File>,
) -> Result<Vec<claw_machine::ClawMachine>, Box<dyn error::Error>> {
    let mut claw_machines: Vec<ClawMachine> = Vec::new();
    let mut line = String::new();
    let button_a_regex = Regex::new(r"Button A: X([+-]?\d+), Y([+-]?\d+)")?;
    let button_b_regex = Regex::new(r"Button B: X([+-]?\d+), Y([+-]?\d+)")?;
    let prize_regex = Regex::new(r"Prize: X=([+-]?\d+), Y=([+-]?\d+)")?;

    loop {
        line.clear();
        input_file_reader.read_line(&mut line)?;
        let Some(button_a_captures) = button_a_regex.captures(&line) else {
            error!(?line, "could not match button a line");
            return Err(Box::from("could not parse button a line"));
        };
        let (_, [button_a_x_increment, button_a_y_increment]) = button_a_captures.extract();
        let button_a_x_increment = button_a_x_increment.parse::<f64>()?;
        let button_a_y_increment = button_a_y_increment.parse::<f64>()?;
        
        line.clear();
        input_file_reader.read_line(&mut line)?;
        let Some(button_b_captures) = button_b_regex.captures(&line) else {
            error!(?line, "could not match button b line");
            return Err(Box::from("could not parse button b line"));
        };
        let (_, [button_b_x_increment, button_b_y_increment]) = button_b_captures.extract();
        let button_b_x_increment = button_b_x_increment.parse::<f64>()?;
        let button_b_y_increment = button_b_y_increment.parse::<f64>()?;
        
        line.clear();
        input_file_reader.read_line(&mut line)?;
        let Some(prize_captures) = prize_regex.captures(&line) else {
            error!(?line, "could not match prize line");
            return Err(Box::from("could not parse prize line"));
        };
        let (_, [prize_x, prize_y]) = prize_captures.extract();
        let prize_x = prize_x.parse::<f64>()?;
        let prize_y = prize_y.parse::<f64>()?;

        let new_claw_machine = claw_machine::new(
            Matrix2::new(
                button_a_x_increment,
                button_b_x_increment,
                button_a_y_increment,
                button_b_y_increment,
            ),
            Vector2::new(prize_x, prize_y),
        );
        trace!("machine"=?new_claw_machine, "parsed claw machine");
        claw_machines.push(new_claw_machine);

        input_file_reader.read_line(&mut line)?;
        if input_file_reader.fill_buf()?.is_empty() {
            break;
        }
    }

    return Ok(claw_machines)
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let claw_machines = match parse_input_to_claw_machines(input_file_reader) {
        Ok(machines) => machines,
        Err(e) => {
            error!("error"=?e, "error occurred during input parsing");
            return None;
        }
    };
    trace!(?claw_machines, "parsed claw machines");

    let total_cost: i64 = claw_machines
        .iter()
        .filter_map(|machine| match machine.solve_machine() {
            Ok(machine_result) => {
                debug!(?machine, ?machine_result, "solved machine");
                Some(machine_result)
            },
            Err(e) => {
                trace!(?machine, "error"=?e, "could not solve machine");
                None
            }
        }).sum();

    Some(total_cost)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
    }

    None
}
