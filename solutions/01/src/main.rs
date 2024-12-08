use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};

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

    let computation_result = match args.part {
        1 => part01(input_file_reader),
        2 => part02(input_file_reader),
        _ => unreachable!(), // clap has filtered out all other possibilities.
    }
    .expect("computation did not produce a value");

    info!(
        "computation_result" = computation_result,
        "computation complete"
    );
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut left_list = Vec::new();
    let mut right_list = Vec::new();
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        let line_parts: Vec<_> = line.split_ascii_whitespace().collect();

        debug!(
            "line" = line,
            "line parts" = format!("{:?}", line_parts),
            "read line from input file"
        );

        let parsed_numbers: Vec<_> = line_parts
            .iter()
            .filter_map(|item| item.parse::<i64>().ok())
            .collect();

        if line_parts.len() != parsed_numbers.len() {
            error!(
                "line" = line,
                "line parts" = format!("{:?}", line_parts),
                "parsed numbers" = format!("{:?}", parsed_numbers),
                "line failed to parse successfully to integer"
            );
            continue;
        }

        left_list.push(parsed_numbers[0]);
        right_list.push(parsed_numbers[1]);
    }

    left_list.sort();
    right_list.sort();

    let total_difference = left_list.iter().zip(right_list.iter())
        .map(|(a,b)| (a-b).abs())
        .sum();
    Some(total_difference)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut left_list = Vec::new();
    let mut right_list_count:HashMap<i64, i64> = HashMap::new();
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        let line_parts: Vec<_> = line.split_ascii_whitespace().collect();

        debug!(
            "line" = line,
            "line parts" = format!("{:?}", line_parts),
            "read line from input file"
        );

        let parsed_numbers: Vec<_> = line_parts
            .iter()
            .filter_map(|item| item.parse::<i64>().ok())
            .collect();

        if line_parts.len() != parsed_numbers.len() {
            error!(
                "line" = line,
                "line parts" = format!("{:?}", line_parts),
                "parsed numbers" = format!("{:?}", parsed_numbers),
                "line failed to parse successfully to integer"
            );
            continue;
        }

        left_list.push(parsed_numbers[0]);
        *right_list_count.entry(parsed_numbers[1]).or_default() += 1;
    }

    let score = left_list.iter()
        .map(|item| right_list_count.get(&item).unwrap_or(&0) * item)
        .sum();
    Some(score)
}
