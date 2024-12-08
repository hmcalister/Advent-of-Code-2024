use std::{fs::File, io::{BufRead, BufReader}};

#[allow(unused_imports)]
use tracing::{error, info, debug, trace};
use clap::Parser;


mod logging;

/// Program to solve Advent of Code puzzles
#[derive(Parser, Debug)]
#[command(version, about)]
struct CommandLineArgs {
    /// File to read input from
    #[arg(long, default_value="puzzleInput")]
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
        _ => unreachable!() // clap has filtered out all other possibilities.
    }.expect("computation did not produce a value");

    info!("computation_result"=computation_result, "computation complete");
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64>{
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line"=line, "read line from input file");
    }
    
    None
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64>{
    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line"=line, "read line from input file");
    }

    None
}