#[allow(unused_imports)]
use tracing::{info, debug, error};
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
    debug!("command line args" = ?args, "parsed command line args");

    let result_option = match args.part {
        1 => part01(),
        2 => part02(),
        _ => unreachable!() // clap has filtered out all other possibilities.
    };

    match result_option {
        Some(result_value) => info!("result_value"=result_value, "computation complete"),
        None => error!("computation did not produce a value")
    }

}

fn part01() -> Option<i64>{
    None
}

fn part02() -> Option<i64>{
    None
}