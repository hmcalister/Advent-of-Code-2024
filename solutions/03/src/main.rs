use regex::Regex;
use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::Parser;
use std::time::SystemTime;
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

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut mul_operation_total = 0;
    let mul_operation_regex = Regex::new(r"mul\((\d+),(\d+)\)").unwrap();
    let all_lines: String = input_file_reader
        .lines()
        .map(|s| s.unwrap())
        .collect::<Vec<_>>()
        .join(" ");

    for (_, [mul_left, mul_right]) in mul_operation_regex
        .captures_iter(&all_lines)
        .map(|capture| capture.extract())
    {
        debug!(
            "mul_left" = mul_left,
            "mul_right" = mul_right,
            "found multiplication"
        );
        let mul_left = match mul_left.parse::<i64>() {
            Ok(num) => num,
            Err(_) => {
                error!("mul_left" = mul_left, "failed to parse integer mul_left");
                continue;
            }
        };
        let mul_right = match mul_right.parse::<i64>() {
            Ok(num) => num,
            Err(_) => {
                error!("mul_right" = mul_right, "failed to parse integer mul_right");
                continue;
            }
        };
        mul_operation_total += mul_left * mul_right
    }

    Some(mul_operation_total)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
    }

    None
}
