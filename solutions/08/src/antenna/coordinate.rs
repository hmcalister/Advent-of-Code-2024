
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Coordinate {
    pub x: i32,
    pub y: i32,
}

impl Coordinate {
    pub fn add_coordinate(&self, other: Coordinate)-> Coordinate {
        Coordinate { x: self.x+other.x, y: self.y+other.y }
    }

    pub fn subtract_coordinate(&self, other: Coordinate)-> Coordinate {
        Coordinate { x: self.x-other.x, y: self.y-other.y }
    }
}