use super::direction::Direction;

#[derive(Debug, Clone, Copy, Hash, Eq, PartialEq)]
pub struct Coordinate {
    x: i32,
    y: i32,
}

pub fn new_from_usize(x: usize, y: usize) -> Coordinate {
    Coordinate {
        x: x as i32,
        y: y as i32,
    }
}

impl Coordinate {
    pub fn in_bounds(&self, lower_x: i32, upper_x: i32, lower_y: i32, upper_y: i32) -> bool {
        lower_x <= self.x && self.x < upper_x && lower_y <= self.y && self.y < upper_y
    }

    pub fn move_in_direction(&self, d: Direction) -> Coordinate {
        match d {
            Direction::Up => Coordinate {
                x: self.x,
                y: self.y - 1,
            },
            Direction::Right => Coordinate {
                x: self.x + 1,
                y: self.y,
            },
            Direction::Down => Coordinate {
                x: self.x,
                y: self.y + 1,
            },
            Direction::Left => Coordinate {
                x: self.x - 1,
                y: self.y,
            },
        }
    }
}
