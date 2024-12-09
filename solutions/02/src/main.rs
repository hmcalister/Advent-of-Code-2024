use std::fs::File;
use std::io::{BufRead, BufReader};
use std::time::SystemTime;

use clap::Parser;
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

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

fn is_safe(levels: Vec<i32>) -> bool {
    if levels.len() < 2 {
        return true;
    }

    let mut previous_value = levels[0];
    let current_value = levels[1];
    let is_increasing = (current_value - previous_value) > 0;

    for level in levels.into_iter().skip(1) {
        let mut difference = level - previous_value;

        if !is_increasing {
            difference *= -1;
        }
        if !(1..=3).contains(&difference) {
            return false;
        }
        previous_value = level;
    }

    true
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut total_safe = 0;
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");

        let parsed_line_result: Result<Vec<_>, _> = line
            .split_ascii_whitespace()
            .map(|item| {
                item.parse::<i32>().map_err(|e| {
                    error!(
                        "value" = item,
                        "error" = format!("{:?}", e),
                        "could not parse value in line"
                    );
                    e
                })
            })
            .collect();

        let parsed_line = match parsed_line_result {
            Ok(parsed_line) => parsed_line,
            Err(_) => continue,
        };

        if is_safe(parsed_line) {
            total_safe += 1;
        }
    }

    Some(total_safe)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
    }

    None
}
