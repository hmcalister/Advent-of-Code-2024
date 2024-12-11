use std::collections::HashSet;

use crate::guard_state::coordinate::Coordinate;

pub const OBSTACLE_BYTE: u8 = b'#';

#[derive(Debug)]
pub struct ObstacleMap {
    obstacle_map: HashSet<Coordinate>,
    width: i32,
    height: i32,
}

pub fn new(width: i32, height: i32) -> ObstacleMap {
    ObstacleMap {
        obstacle_map: HashSet::new(),
        width: width,
        height: height,
    }
}

impl ObstacleMap {
    pub fn is_obstacle(&self, coord: Coordinate) -> Option<bool> {
        if !coord.in_bounds(0, self.width, 0, self.height) {
            return None;
        } else {
            return Some(self.obstacle_map.contains(&coord));
        }
    }

    pub fn add_obstacle(&mut self, coord: Coordinate) {
        self.obstacle_map.insert(coord);
    }

    pub fn remove_obstacle(&mut self, coord: Coordinate) -> bool {
        self.obstacle_map.remove(&coord)
    }
}
