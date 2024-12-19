#[derive(Debug, Clone, Copy, Hash, Eq, PartialEq)]
pub enum Direction {
    Up,
    Right,
    Left,
    Down,
}

#[allow(dead_code)]
impl Direction {
    pub fn rotate_right(&self) -> Direction {
        match self {
            Direction::Up => Direction::Right,
            Direction::Right => Direction::Down,
            Direction::Down => Direction::Left,
            Direction::Left => Direction::Up,
        }
    }

    pub fn rotate_left(&self) -> Direction {
        match self {
            Direction::Up => Direction::Left,
            Direction::Right => Direction::Up,
            Direction::Down => Direction::Right,
            Direction::Left => Direction::Down,
        }
    }
}