use std::error;
use nalgebra::{Matrix2, Vector2};

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};


pub const PRIZE_POSITION_OFFSET: Vector2<f64> = Vector2::new(10000000000000.0, 10000000000000.0);
const BUTTON_COSTS: Vector2<f64> = Vector2::new(3.0, 1.0);
const FLOAT_CLOSE_EPSILON: f64 = 0.0001;

#[derive(Debug)]
pub struct ClawMachine {
    button_increments: Matrix2<f64>,
    prize_position: Vector2<f64>,
}

pub fn new(button_increments: Matrix2<f64>, prize_position: Vector2<f64>) -> ClawMachine {
    ClawMachine {
        button_increments,
        prize_position,
    }
}

fn check_float_close_to_int(f: f64) -> bool {
    (f.round() - f).abs() < FLOAT_CLOSE_EPSILON
}

impl ClawMachine {
    pub fn solve_machine(&self) -> Result<i64, Box<dyn error::Error>> {
        let Some(mut button_pushes) = self.button_increments.lu().solve(&self.prize_position) else {
            return Err(Box::from("machine button increments not invertible"));
        };
        
        trace!(?self, ?button_pushes, "machine solved");
        if !(check_float_close_to_int(button_pushes[0])
            && check_float_close_to_int(button_pushes[1]))
        {
            return Err(Box::from("machine button pushes not integers!"));
        }
        button_pushes[0] = button_pushes[0].round();
        button_pushes[1] = button_pushes[1].round();
        // assert_eq!(self.button_increments*button_pushes, self.prize_position);
        
        if !(button_pushes[0] >= 0.0 && button_pushes[1] >= 0.0) {
            return Err(Box::from("machine button pushes not positive!"));
        }


        Ok(button_pushes.dot(&BUTTON_COSTS) as i64)
    }

    pub fn update_prize_position(&mut self, prize_position_increment: Vector2<f64>) {
        self.prize_position = self.prize_position + prize_position_increment;
    }
}
