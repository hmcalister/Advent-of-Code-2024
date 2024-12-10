use std::{
    fs::File,
    io::{BufRead, BufReader},
};

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

const XMAS_BYTES: [u8; 4] = [b'X', b'M', b'A', b'S'];
const CROSSED_MAS_BYTES: [u8; 4] = [b'M', b'M', b'S', b'S'];

#[derive(Debug)]
pub struct WordSearch {
    grid_data: Vec<u8>,
    height: usize,
    width: usize,
}

pub fn new_word_search_from_input(input_file_reader: BufReader<File>) -> WordSearch {
    let mut search = WordSearch {
        width: 0,
        height: 0,
        grid_data: Vec::new(),
    };

    for line_result in input_file_reader.lines() {
        let line = line_result.unwrap();
        debug!("line" = line, "read line from input file");
        search.grid_data.extend_from_slice(line.as_bytes());
        search.width = line.len();
        search.height += 1;
    }

    debug!("word search"=?search, "parsed word search");
    search
}

impl WordSearch {
    fn cartesian_to_linear_coordinate(&self, x: usize, y: usize) -> usize {
        y * self.width + x
    }

    fn linear_to_cartesian_coordinate(&self, i: usize) -> (usize, usize) {
        (i % self.width, i / self.width)
    }

    fn count_xmas_at_coordinate(&self, x: usize, y: usize) -> i64 {
        let mut total_xmas = 0;

        // Row Forward
        if x + 4 <= self.width
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 1, y)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 2, y)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 3, y)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas row forward");
            total_xmas += 1;
        }

        // Row Backward
        if x >= 3
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 1, y)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 2, y)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 3, y)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas row backward");
            total_xmas += 1;
        }

        // Column Down
        if y + 4 <= self.height
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y + 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y + 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y + 3)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas column down");
            total_xmas += 1;
        }

        // Column Up
        if y >= 3
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y - 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y - 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y - 3)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas column up");
            total_xmas += 1;
        }

        // Diagonal Up Right
        if x + 4 <= self.width
            && y >= 3
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 1, y - 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 2, y - 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 3, y - 3)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas diagonal up right");
            total_xmas += 1;
        }

        // Diagonal Down Right
        if x + 4 <= self.width
            && y + 4 <= self.height
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 1, y + 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 2, y + 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x + 3, y + 3)] == XMAS_BYTES[3]
        {
            debug!(
                "start x" = x,
                "start y" = y,
                "found xmas diagonal down right"
            );
            total_xmas += 1;
        }

        // Diagonal Down Left
        if x >= 3
            && y + 4 <= self.height
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 1, y + 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 2, y + 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 3, y + 3)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas diagonal down left");
            total_xmas += 1;
        }

        // Diagonal Up Left
        if x >= 3
            && y >= 3
            && self.grid_data[self.cartesian_to_linear_coordinate(x, y)] == XMAS_BYTES[0]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 1, y - 1)] == XMAS_BYTES[1]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 2, y - 2)] == XMAS_BYTES[2]
            && self.grid_data[self.cartesian_to_linear_coordinate(x - 3, y - 3)] == XMAS_BYTES[3]
        {
            debug!(x, y, "found xmas diagonal up left");
            total_xmas += 1;
        }

        total_xmas
    }

    pub fn find_all_xmas(&self) -> i64 {
        // let mut total_xmas = 0;
        // for y in 0..self.height {
        //     for x in 0..self.width {
        //         debug!(
        //             x,
        //             y,
        //             "linear coordinate" = self.cartesian_to_linear_coordinate(x, y),
        //             "attempting coordinate"
        //         );
        //         total_xmas += self.count_xmas_at_coordinate(x, y)
        //     }
        // }
        // total_xmas

        (0..self.grid_data.len())
            .map(|linear_coordinate| self.linear_to_cartesian_coordinate(linear_coordinate))
            .map(|(x, y)| self.count_xmas_at_coordinate(x, y))
            .sum::<i64>()
    }
}
