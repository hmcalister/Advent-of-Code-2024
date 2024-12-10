use std::fs::File;
use std::io::{BufRead, BufReader};

use clap::Parser;
use std::time::SystemTime;
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

mod dependency_graph;
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

fn parse_dependency_graph_and_orders(
    input_file_reader: BufReader<File>,
) -> (dependency_graph::DependencyGraph, Vec<Vec<i32>>) {
    let mut dep_graph = dependency_graph::new_dependency_graph();

    let mut all_lines = input_file_reader.lines().into_iter();

    for line_result in &mut all_lines {
        let line = line_result.unwrap();
        if line.len() == 0 {
            break;
        }
        debug!("line" = line, "read line from input file");
        let split: Vec<&str> = line.split("|").collect();
        if split.len() != 2 {
            error!(line, ?split, "dependency does not have correct form");
            continue;
        }
        let (prior_item_str, posterior_item_str) = (split[0], split[1]);
        let prior_item = match prior_item_str.parse() {
            Ok(num) => num,
            Err(e) => {
                error!("error"=?e, prior_item_str, "failed to parse prior item");
                continue;
            }
        };
        let posterior_item = match posterior_item_str.parse() {
            Ok(num) => num,
            Err(e) => {
                error!("error"=?e, posterior_item_str, "failed to parse prior item");
                continue;
            }
        };

        dep_graph.add_dependency(prior_item, posterior_item);
    }

    let mut orders = Vec::new();
    for line_result in &mut all_lines {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");

        match line
            .split(",")
            .into_iter()
            .map(|item_str| item_str.parse::<i32>())
            .collect::<Result<Vec<_>, _>>()
        {
            Ok(order) => {
                debug!(?order, "parsed order");
                orders.push(order);
            }
            Err(e) => {
                error!("error"=?e, ?line, "order line not parsable");
                continue;
            }
        };
    }

    (dep_graph, orders)
}

fn part01(input_file_reader: BufReader<File>) -> Option<i64> {
    let (dep_graph, orders) = parse_dependency_graph_and_orders(input_file_reader);
    debug!(?dep_graph, "dependency graph parsed");

    let mut middle_number_sum: i64 = 0;
    for order in orders {
        if let Some(topologically_sorted_order) = dep_graph.topological_sort(&order) {
            if topologically_sorted_order
                .iter()
                .zip(order.iter())
                .all(|(a, b)| *a == *b)
            {
                let middle_item = order[order.len() / 2];
                middle_number_sum += middle_item as i64;
                info!(
                    ?order,
                    ?middle_item,
                    "updated total" = middle_number_sum,
                    "found valid order"
                );
            }
        };
    }
    Some(middle_number_sum)
}

fn part02(input_file_reader: BufReader<File>) -> Option<i64> {
    let (dep_graph, orders) = parse_dependency_graph_and_orders(input_file_reader);
    debug!(?dep_graph, "dependency graph parsed");

    let mut middle_number_sum: i64 = 0;
    for order in orders {
        if let Some(topologically_sorted_order) = dep_graph.topological_sort(&order) {
            if !topologically_sorted_order
                .iter()
                .zip(order.iter())
                .all(|(a, b)| *a == *b)
            {
                let middle_item = topologically_sorted_order[topologically_sorted_order.len() / 2];
                middle_number_sum += middle_item as i64;
                info!(
                    ?topologically_sorted_order,
                    ?middle_item,
                    "updated total" = middle_number_sum,
                    "found invalid order"
                );
            }
        };
    }
    Some(middle_number_sum)
}
