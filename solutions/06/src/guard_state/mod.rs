#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

pub mod coordinate;
pub mod direction;

use coordinate::Coordinate;
use direction::Direction;

#[derive(Debug, Clone, Copy, Hash, Eq, PartialEq)]
pub struct GuardState {
    coordinate: Coordinate,
    direction: Direction,
}

pub fn new(c: Coordinate, d: Direction) -> GuardState {
    GuardState {
        coordinate: c,
        direction: d,
    }
}

impl GuardState {
    pub fn get_coordinate(&self) -> Coordinate {
        self.coordinate
    }

    pub fn step(&self) -> GuardState {
        GuardState {
            coordinate: self.coordinate.move_in_direction(self.direction),
            direction: self.direction,
        }
    }

    pub fn encounter_obstacle(&self) -> GuardState {
        GuardState {
            coordinate: self.coordinate,
            direction: self.direction.rotate_right(),
        }
    }
}
