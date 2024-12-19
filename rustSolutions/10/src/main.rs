use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::Parser;
use std::time::SystemTime;
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

mod logging;
mod gridutils;
mod topographic_map;

use topographic_map::*;

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
    let mut map_data = Vec::new();
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
        map_data.push(line);
    }
    let topomap = new_topographic_map(map_data);
    let total_score = topomap.count_trailhead_scores();

    Some(total_score)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    let mut map_data = Vec::new();
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
        map_data.push(line);
    }
    let topomap = new_topographic_map(map_data);
    let total_rating = topomap.count_trailhead_ratings();

    Some(total_rating)
}
