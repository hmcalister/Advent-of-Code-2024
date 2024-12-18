use std::collections::{HashMap, HashSet};
use rayon::prelude::*;

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

// use crate::gridutils::direction::*;
use crate::gridutils::vector::*;

#[allow(unused)]
#[derive(Debug)]
pub struct TopographicMap {
    trailheads: Vec<LatticeVector2D>,
    coordinate_heights: HashMap<LatticeVector2D, i32>,
    map_width: usize,
    map_height: usize,
    max_height: i32,
}

pub fn new_topographic_map(map_data: Vec<String>) -> TopographicMap {
    let mut m = TopographicMap {
        trailheads: Vec::new(),
        coordinate_heights: HashMap::new(),
        map_width: map_data[0].len(),
        map_height: map_data.len(),
        max_height: 0,
    };

    debug!("got here");

    for (y, line) in map_data.iter().enumerate() {
        for (x, cell) in line.as_bytes().iter().enumerate() {
            let coordinate_height = (cell - b'0') as i32;
            let coordinate = LatticeVector2D {
                x: x as i32,
                y: y as i32,
            };
            debug!(?coordinate, ?coordinate_height, "next coordinate parsed");
            m.coordinate_heights.insert(coordinate, coordinate_height);
            if coordinate_height == 0 {
                debug!(?coordinate, "found trailhead");
                m.trailheads.push(coordinate);
            }
            if coordinate_height > m.max_height {
                debug!(?coordinate_height, "found new max height");
                m.max_height = coordinate_height
            }
        }
    }

    m
}

impl TopographicMap {
    /// Find and return the unique max height coordinates reachable from a given coordinate.
    /// Note the given coordinate may not have a height of zero to allow for a recursive solution.
    fn find_reachable_trailends(&self, coordinate: LatticeVector2D) -> HashSet<LatticeVector2D> {
        let mut reachable_trailends = HashSet::new();

        // If we do not have a height, then we cannot have a score
        let Some(&current_height) = self.coordinate_heights.get(&coordinate) else {
            return reachable_trailends;
        };

        // If we are already max height we are done
        if current_height == self.max_height {
            reachable_trailends.insert(coordinate);
            return reachable_trailends;
        };

        // We are not max height, so walk over neighbors and sum their heights
        let _: Vec<_> = coordinate
            .iter_neighbors()
            .filter(|neighbor| {
                let Some(&neighbor_height) = self.coordinate_heights.get(&neighbor) else {
                    return false;
                };
                neighbor_height == current_height + 1
            })
            .map(|neighbor| {
                let neighbor_reachable_trailheads = self.find_reachable_trailends(neighbor);
                reachable_trailends.extend(neighbor_reachable_trailheads);
            })
            .collect();

        reachable_trailends
    }

    /// Find and return the number of paths to max height coordinates from a given coordinate.
    /// Note the given coordinate may not have a height of zero to allow for a recursive solution.
    fn find_coordinate_rating(&self, coordinate: LatticeVector2D) -> i64 {
        // If we do not have a height, then we cannot have a score
        let Some(&current_height) = self.coordinate_heights.get(&coordinate) else {
            return 0;
        };

        // If we are already max height we are done
        if current_height == self.max_height {
            return 1;
        };

        // We are not max height, so walk over neighbors and sum their heights
        coordinate
            .iter_neighbors()
            .filter(|neighbor| {
                let Some(&neighbor_height) = self.coordinate_heights.get(&neighbor) else {
                    return false;
                };
                neighbor_height == current_height+1
            })
            .map(|neighbor| self.find_coordinate_rating(neighbor))
            .sum()
    }

    pub fn count_trailhead_scores(&self) -> i64 {
        self.trailheads
            .par_iter()
            .map(|trailhead| {
                let reachable_trailheads = self.find_reachable_trailends(*trailhead);
                let trailhead_score = reachable_trailheads.len() as i64;
                debug!(?trailhead, ?trailhead_score, "found score for trailhead");
                trailhead_score
            })
            .sum()
    }

    pub fn count_trailhead_ratings(&self) -> i64 {
        self.trailheads
            .par_iter()
            .map(|trailhead| {
                let trailhead_rating = self.find_coordinate_rating(*trailhead);
                debug!(?trailhead, ?trailhead_rating, "found rating for trailhead");
                trailhead_rating
            })
            .sum()
    }
}
