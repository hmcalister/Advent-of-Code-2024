use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::{BufRead, BufReader};

use tracing::{span, Level};
#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

mod coordinate;
use coordinate::Coordinate;

#[derive(Debug)]
pub struct AntennaMap {
    map_width: usize,
    map_height: usize,

    // A list of all frequencies present in the map
    antenna_frequency_list: HashSet<u8>,

    // A map from frequency to a list of the antennas of that frequency
    antenna_positions_by_frequency: HashMap<u8, Vec<Coordinate>>,
}

pub fn new_antenna_map(input_file_reader: BufReader<File>) -> Option<AntennaMap> {
    let mut antenna_map = AntennaMap {
        map_width: 0,
        map_height: 0,
        antenna_frequency_list: HashSet::new(),
        antenna_positions_by_frequency: HashMap::new(),
    };

    for (y, line_result) in input_file_reader.lines().enumerate() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");

        antenna_map.map_width = line.len();
        antenna_map.map_height = y+1;

        let _ = line
            .as_bytes()
            .iter()
            .enumerate()
            .filter(|(_, c)| **c != b'.')
            .map(|(x, c)| {
                let coord = Coordinate { x: x as i32, y: y as i32 };
                debug!(?coord, "frequency"=?c, "found antenna");
                antenna_map.antenna_frequency_list.insert(*c);
                antenna_map
                    .antenna_positions_by_frequency
                    .entry(*c)
                    .or_default()
                    .push(coord);
            })
            .collect::<Vec<_>>();
    }

    return Some(antenna_map);
}

