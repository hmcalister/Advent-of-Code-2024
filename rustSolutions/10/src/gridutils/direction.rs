#![allow(dead_code)]

use super::vector::LatticeVector2D;

pub const ALL_DIRECTIONS: [LatticeDirection2D; 4] = [
    LatticeDirection2D::Up,
    LatticeDirection2D::Right,
    LatticeDirection2D::Down,
    LatticeDirection2D::Left,
];

#[derive(PartialEq, Eq, Hash)]
#[derive(Clone, Copy)]
#[derive(Debug)]
/// Representation of a direction on a 2D lattice.
/// 
/// These are the cardinal directions Up, Right, Down, and Left.
/// 
/// Note that for consistency, the lattice is orientated such that the y-axis points *down*.
/// This means that larger y values are *lower* on the grid.
pub enum LatticeDirection2D {
    Up,
    Right,
    Down,
    Left,
}

impl LatticeDirection2D {
    /// Convert the direction to a unit-vector pointing in that direction.
    /// Example:
    /// ```
    /// let d = LatticeDirection2D::Up;
    /// assert_eq!(d.to_unit_vector(), LatticeVector2D{x:0, y:-1});
    /// ```
    pub fn to_unit_vector(&self) -> LatticeVector2D {
        match self {
            LatticeDirection2D::Up => LatticeVector2D { x: 0, y: -1 },
            LatticeDirection2D::Right => LatticeVector2D { x: 1, y: 0 },
            LatticeDirection2D::Down => LatticeVector2D { x: 0, y: 1 },
            LatticeDirection2D::Left => LatticeVector2D { x: -1, y: 0 },
        }
    }

    /// Rotate the direction right one step.
    pub fn rotate_right(&self) -> LatticeDirection2D {
        match self {
            LatticeDirection2D::Up => LatticeDirection2D::Right,
            LatticeDirection2D::Right => LatticeDirection2D::Down,
            LatticeDirection2D::Down => LatticeDirection2D::Left,
            LatticeDirection2D::Left => LatticeDirection2D::Up,
        }
    }

    /// Rotate the direction left one step.
    pub fn rotate_left(&self) -> LatticeDirection2D {
        match self {
            LatticeDirection2D::Up => LatticeDirection2D::Left,
            LatticeDirection2D::Right => LatticeDirection2D::Up,
            LatticeDirection2D::Down => LatticeDirection2D::Right,
            LatticeDirection2D::Left => LatticeDirection2D::Down,
        }
    }
}
