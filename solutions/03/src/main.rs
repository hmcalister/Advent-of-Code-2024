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

fn compute_all_multiplications(program_string: &str) -> i64 {
    let mut mul_operation_total = 0;
    let mul_operation_regex = Regex::new(r"mul\((\d+),(\d+)\)").unwrap();

    for (_, [mul_left, mul_right]) in mul_operation_regex
        .captures_iter(program_string)
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

    mul_operation_total
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let all_lines: String = input_file_reader
        .lines()
        .map(|s| s.unwrap())
        .collect::<Vec<_>>()
        .join(" ");

    let mul_operation_total = compute_all_multiplications(&all_lines);

    Some(mul_operation_total)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut mul_operation_total = 0;
    let enable_regex = Regex::new(r"do\(\)").unwrap();
    let disable_regex = Regex::new(r"don't\(\)").unwrap();

    let all_input = input_file_reader
        .lines()
        .map(|s| s.unwrap())
        .collect::<Vec<_>>()
        .join(" ");
    let mut remaining_input: &str = all_input.as_str();
    let mut current_enabled_segment;

    while remaining_input.len() > 0 {
        // Find next disable index

        let next_disable_index = match disable_regex.find(&remaining_input) {
            Some(m) => {
                debug!("match"=?m, "remaining input length"=remaining_input.len(), "found next disable command");
                m.end()
            }
            None => {
                debug!(
                    "remaining input length" = remaining_input.len(),
                    "no disable command found"
                );
                remaining_input.len()
            }
        };

        (current_enabled_segment, remaining_input) = remaining_input.split_at(next_disable_index);
        trace!("current enabled segment"=current_enabled_segment, "enabled segment found");
        mul_operation_total += compute_all_multiplications(current_enabled_segment);

        let next_enable_index = match enable_regex.find(&remaining_input) {
            Some(m) => {
                debug!("match"=?m, "remaining input length"=remaining_input.len(), "found next enable command");
                m.start()
            }
            None => {
                debug!(
                    "remaining input length" = remaining_input.len(),
                    "no enable command found"
                );
                remaining_input.len()
            }
        };
        (_, remaining_input) = remaining_input.split_at(next_enable_index);
    }

    Some(mul_operation_total)
}
