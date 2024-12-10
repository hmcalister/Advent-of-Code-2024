use std::collections::{HashMap, HashSet};

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};
use tracing::{span, Level};

#[derive(Debug)]
pub struct DependencyGraph {
    // Dependencies encoded as key item depends on (must come after) all items in the value vector
    dependencies: HashMap<i32, Vec<i32>>,
}

pub fn new_dependency_graph() -> DependencyGraph {
    DependencyGraph {
        dependencies: HashMap::new(),
    }
}

impl DependencyGraph {
    pub fn add_dependency(&mut self, prior_item: i32, posterior_item: i32) {
        self.dependencies
            .entry(prior_item)
            .or_default()
            .push(posterior_item);
    }

    pub fn topological_sort(&self, order: &Vec<i32>) -> Option<Vec<i32>> {
        let _span = span!(Level::DEBUG, "order", ?order).entered();
        let mut order_items = HashSet::new();
        let mut item_indegrees: HashMap<i32, usize> = HashMap::new();
        for item in order {
            order_items.insert(*item);
            item_indegrees.insert(*item, 0);
        }
        debug!(?order_items, "order items found");

        let mut order_specific_dependencies: HashMap<i32, Vec<i32>> = HashMap::new();
        for item in order {
            let mut order_specific_item_dependencies: Vec<i32> = Vec::new();
            let all_dependencies = match self.dependencies.get(item) {
                Some(dep) => dep,
                None => &Vec::new(),
            };

            let _ = all_dependencies
                .into_iter()
                .filter(|dependency| order_items.contains(dependency))
                .map(|posterior_item| {
                    order_specific_item_dependencies.push(*posterior_item);
                    item_indegrees
                        .entry(*posterior_item)
                        .and_modify(|indegree| *indegree += 1);
                })
                .collect::<Vec<_>>();

            order_specific_dependencies.insert(*item, order_specific_item_dependencies);
        }

        let mut zero_indegree_items: Vec<i32> = Vec::new();
        let _ = item_indegrees
            .iter()
            .filter(|(_, indegree)| **indegree == 0)
            .map(|(item, _)| zero_indegree_items.push(*item))
            .collect::<Vec<_>>();
        debug!(
            ?item_indegrees,
            ?order_specific_dependencies,
            ?zero_indegree_items,
            "order specific dependencies found for order"
        );

        let mut topologically_sorted_order = Vec::new();
        while zero_indegree_items.len() > 0 {
            let next_zero_indegree_item = zero_indegree_items.pop().unwrap();
            topologically_sorted_order.push(next_zero_indegree_item);
            for posterior_item in order_specific_dependencies
                .get(&next_zero_indegree_item)
                .unwrap()
            {
                item_indegrees
                    .entry(*posterior_item)
                    .and_modify(|indegree| {
                        if *indegree - 1 == 0 {
                            zero_indegree_items.push(*posterior_item);
                        }
                        *indegree -= 1
                    });
            }
        }
        if topologically_sorted_order.len() != order.len() {
            error!(?order_specific_dependencies, "order cannot be topologically sorted, a cycle must exist");
            return None
        }

        debug!(?topologically_sorted_order, "topological sort complete");
        Some(topologically_sorted_order)
    }
}
