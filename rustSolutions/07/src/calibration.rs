#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

pub mod operation;

#[derive(Debug)]
pub struct Calibration {
    target_value: i64,
    calibration_data: Vec<i64>,
}

pub fn parse_line_to_calibration(line: &str) -> Option<Calibration> {
    let [target_value_str, calibration_data_str] = line.split(": ").collect::<Vec<_>>()[..] else {
        error!(
            ?line,
            "cannot split line into target value and calibration data on ': '"
        );
        return None;
    };

    let target_value = match target_value_str.parse::<i64>() {
        Ok(num) => num,
        Err(e) => {
            error!(
                ?target_value_str,
                ?line,
                "error"=?e,
                "cannot parse target value to integer"
            );
            return None;
        }
    };

    let calibration_data = match calibration_data_str
        .split(" ")
        .map(|item| item.parse::<i64>())
        .collect::<Result<Vec<_>, _>>()
    {
        Ok(data) => data,
        Err(e) => {
            error!(
                ?calibration_data_str,
                ?line,
                "error"=?e,
                "cannot parse calibration data to integer"
            );
            return None;
        }
    };

    return Some(Calibration {
        target_value: target_value,
        calibration_data: calibration_data,
    });
}

impl Calibration {
    pub fn is_valid(&self, possible_operations: &Vec<operation::Operation>) -> bool {
        self.check_is_valid_recursive(possible_operations, 0, 0)
    }

    fn check_is_valid_recursive(
        &self,
        possible_operations: &Vec<operation::Operation>,
        current_index: usize,
        current_value: i64,
    ) -> bool {
        trace!(?current_index, ?current_value, "intermediate computation");
        if current_index == self.calibration_data.len() && current_value == self.target_value {
            return true;
        }
        if current_index >= self.calibration_data.len() || current_value > self.target_value {
            return false;
        }

        possible_operations.iter().any(|op| {
            self.check_is_valid_recursive(
                possible_operations,
                current_index + 1,
                op(current_value, self.calibration_data[current_index]),
            )
        })
    }

    pub fn get_target_value(&self) -> i64 {
        self.target_value
    }
}
