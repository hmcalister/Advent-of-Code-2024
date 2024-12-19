#![allow(dead_code)]

use super::direction::{LatticeDirection2D, ALL_DIRECTIONS};

#[derive(PartialEq, Eq, Hash)]
#[derive(Clone, Copy)]
#[derive(Debug)]
/// A representation of the lattice points of a grid (integer coordinates only).
/// 
/// Note that for consistency, the lattice is orientated such that the y-axis points *down*.
/// This means that larger y values are *lower* on the grid.
pub struct LatticeVector2D {
    pub x: i32,
    pub y: i32,
}

impl LatticeVector2D {
    /// Add two lattice vectors together, returning the result.
    pub fn add(&self, other: &LatticeVector2D) -> LatticeVector2D {
        LatticeVector2D {
            x: self.x + other.x,
            y: self.y + other.y,
        }
    }
    
    /// Subtract two lattice vectors. The vector passed as a parameter is subtracted from the method receiver.
    /// i.e. `a.subtract(b)` is equivalent to `a-b`.
    pub fn subtract(&self, other: &LatticeVector2D) -> LatticeVector2D {
        LatticeVector2D {
            x: self.x - other.x,
            y: self.y - other.y,
        }
    }

    /// Step a lattice point in the specific direction.
    /// Example:
    /// ```
    /// let c = LatticeVector2D{x:0, y: 0};
    /// let d = LatticeDirection2D::Up;
    /// let c_step = c.step(d);
    /// assert_eq!(c_step, LatticeVector2D{x:0, y:-1});
    /// ```
    pub fn step(&self, direction: &LatticeDirection2D) -> LatticeVector2D {
        let direction_coordinate = direction.to_unit_vector();
        self.add(&direction_coordinate)
    }

    /// Iterate over the cardinal neighbors of this coordinate
    pub fn iter_neighbors<'a>(&'a self) -> impl Iterator<Item = LatticeVector2D> + use<'a> {
        ALL_DIRECTIONS.iter().map(|d| {
            self.step(d)
        })
    }
}
