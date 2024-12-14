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

impl AntennaMap {
    fn check_coordinate_inbounds(&self, coord: Coordinate) -> bool {
        coord.x >= 0 && coord.x < self.map_width as i32 && coord.y >= 0 && coord.y < self.map_height as i32
    }

    pub fn count_first_order_antinodes(&self) -> i64 {
        let mut valid_antinodes = HashSet::new();

        for freq in &self.antenna_frequency_list {
            let _span = span!(Level::DEBUG, "finding first order antinodes", "antenna frequency"=?freq).entered();
            let antenna_positions = self.antenna_positions_by_frequency.get(&freq).unwrap();

            for first_antenna_index in 0..antenna_positions.len() {
                let first_antenna_position = antenna_positions[first_antenna_index];
                for second_antenna_index in first_antenna_index+1..antenna_positions.len() {
                    let second_antenna_position = antenna_positions[second_antenna_index];

                    // Ordering (first, second)
                    let position_delta = first_antenna_position.subtract_coordinate(second_antenna_position);
                    let potential_antinode = first_antenna_position.add_coordinate(position_delta);
                    if self.check_coordinate_inbounds(potential_antinode) {
                        debug!("antinode coordinate" = ?potential_antinode, "found antinode");
                        valid_antinodes.insert(potential_antinode);
                    }
                    
                    // Ordering (second, first)
                    let position_delta = second_antenna_position.subtract_coordinate(first_antenna_position);
                    let potential_antinode = second_antenna_position.add_coordinate(position_delta);
                    if self.check_coordinate_inbounds(potential_antinode) {
                        debug!("antinode coordinate" = ?potential_antinode, "found antinode");
                        valid_antinodes.insert(potential_antinode);
                    }
                }
            }
        }

        valid_antinodes.len() as i64
    }

    pub fn count_all_antinodes(&self) -> i64 {
        let mut valid_antinodes = HashSet::new();

        for freq in &self.antenna_frequency_list {
            let _span = span!(Level::DEBUG, "finding first order antinodes", "antenna frequency"=?freq).entered();
            let antenna_positions = self.antenna_positions_by_frequency.get(&freq).unwrap();

            for first_antenna_index in 0..antenna_positions.len() {
                let first_antenna_position = antenna_positions[first_antenna_index];
                
                // Insert antinode that exists on current antenna
                valid_antinodes.insert(first_antenna_position);
                for second_antenna_index in first_antenna_index+1..antenna_positions.len() {
                    let second_antenna_position = antenna_positions[second_antenna_index];

                    // Ordering (first, second)
                    let position_delta = first_antenna_position.subtract_coordinate(second_antenna_position);
                    let mut potential_antinode = first_antenna_position.add_coordinate(position_delta);
                    while self.check_coordinate_inbounds(potential_antinode) {
                        debug!("antinode coordinate" = ?potential_antinode, "found antinode");
                        valid_antinodes.insert(potential_antinode);
                        potential_antinode = potential_antinode.add_coordinate(position_delta);
                    }

                    let position_delta = second_antenna_position.subtract_coordinate(first_antenna_position);
                    let mut potential_antinode = second_antenna_position.add_coordinate(position_delta);
                    while self.check_coordinate_inbounds(potential_antinode) {
                        debug!("antinode coordinate" = ?potential_antinode, "found antinode");
                        valid_antinodes.insert(potential_antinode);
                        potential_antinode = potential_antinode.add_coordinate(position_delta);
                    }
                }
            }
        }

        valid_antinodes.len() as i64
    }

}
