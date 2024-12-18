use rustc_hash::FxHashMap;

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};

#[derive(Debug)]
pub struct StoneData {
    stone_counts: FxHashMap<u64, i64>,
}

pub fn new(initial_stones: &Vec<u64>) -> StoneData {
    let mut stone_data = StoneData {
        stone_counts: FxHashMap::default(),
    };

    for stone in initial_stones.iter() {
        *stone_data.stone_counts.entry(*stone).or_default() += 1;
    }

    stone_data
}

impl StoneData {
    pub fn blink(&mut self) {
        let mut new_stone_counts: FxHashMap<u64, i64> = FxHashMap::default();
        for (&stone_value, &count) in &self.stone_counts {
            trace!(?stone_value, ?count, "computing stone blink update");
            if stone_value == 0 {
                *new_stone_counts.entry(1).or_default() += count;
            } else if stone_value.to_string().len() % 2 == 0 {
                let stone_value_str = stone_value.to_string();
                let stone_value_left = &stone_value_str[..stone_value_str.len()/2].parse::<u64>().unwrap();
                let stone_value_right = &stone_value_str[stone_value_str.len()/2..].parse::<u64>().unwrap();

                *new_stone_counts.entry(*stone_value_left).or_default() += count;
                *new_stone_counts.entry(*stone_value_right).or_default() += count;
            } else {
                *new_stone_counts.entry(2024*stone_value).or_default() += count;
            }
        }
        self.stone_counts = new_stone_counts;
    }

    pub fn total_stones(&self) -> i64 {
        self.stone_counts.values().sum()
    }
}
