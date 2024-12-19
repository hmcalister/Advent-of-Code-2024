use std::collections::{HashMap, HashSet};

#[allow(unused_imports)]
use tracing::{debug, error, info, trace};
use tracing::{span, Level};

#[derive(Debug)]
pub struct DependencyGraph {
    // Dependencies encoded as key item depends on (must come after) all items in the value set
    dependencies: HashMap<i32, HashSet<i32>>,
}

/// Create a new dependency graph with no dependencies.
///
/// Dependencies are added later using `graph.add_dependency(a,b);`.
pub fn new_dependency_graph() -> DependencyGraph {
    DependencyGraph {
        dependencies: HashMap::new(),
    }
}

impl DependencyGraph {
    /// Given a prior and posterior item, enforce the dependency that
    /// prior item must exist before posterior item, equivalently that posterior item depends on prior item.
    ///
    /// Note that dependencies are stored as a set, hence repeated calls to `graph.add_dependency(a,b)` for the same pair `(a,b)` are idempotent.
    pub fn add_dependency(&mut self, prior_item: i32, posterior_item: i32) {
        self.dependencies
            .entry(prior_item)
            .or_default()
            .insert(posterior_item);
    }

    /// Perform a topological sort on the order given by a vector of items.
    ///
    /// Sorts the items such that all items come after their dependencies.
    /// If this is not possible (e.g. because of cyclical dependencies) then returns None.
    ///
    /// Note this method first filters out any irrelevant dependencies, such that cyclical / transitive dependencies 
    /// in the graph will only prevent sorting as long as the cycle consists of only items from the order.
    /// For example, if the graph has dependencies {a->b, b->c, c->d, d->a}, the order [a,b,c,d] can *not* be sorted, but the order [c,b,a] *can* be (the result is [a,b,c]).
    pub fn topological_sort(&self, order: &Vec<i32>) -> Option<Vec<i32>> {
        let _span = span!(Level::DEBUG, "order", ?order, "topological sort").entered();

        // Make a quick lookup for all items in the order
        let mut order_items = HashSet::new();
        // Track all indegrees for each item (the number of items that depend on this item)
        let mut item_indegrees: HashMap<i32, usize> = HashMap::new();
        for item in order {
            order_items.insert(*item);
            // Items start with indegree zero, updated as we walk the dependency graph
            item_indegrees.insert(*item, 0);
        }

        // Filter out only the dependencies that are relevant to this order, i.e. only items in this order
        let mut order_specific_dependencies: HashMap<i32, Vec<i32>> = HashMap::new();
        for item in order {
            // Track the relevant dependencies for the current item
            let mut order_specific_item_dependencies: Vec<i32> = Vec::new();

            // Get all dependencies for a given item. If no dependencies exist, use an empty set to show this
            let all_item_dependencies = match self.dependencies.get(item) {
                Some(dep) => dep,
                None => &HashSet::new(),
            };

            // Walk over all dependencies, filtering only the relevant dependencies, and adding these to the list
            // Also increment the indegree as we go
            let _ = all_item_dependencies
                .iter()
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

        // Track all items that have zero indegree --- these can be added to the topologically sorted order immediately as they have no further dependencies.
        // This vector is treated as a stack.
        let mut zero_indegree_items: Vec<i32> = Vec::new();
        
        // Find the existing items with zero indegree (if they exist).
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

        // Hold the final topologically sorted order.
        let mut topologically_sorted_order = Vec::new();

        // Pop each zero indegree item, 
        // add it to the topologically sorted vector, 
        // and decrement the dependant items (hopefully creating new zero indegree items in the process).
        //
        // If there are no more zero indegree items, we cannot continue to add these, 
        // but we may not have added all items to the topologically sorted vector.
        // This is handled in the next if statement.
        while let Some(next_zero_indegree_item) = zero_indegree_items.pop()  {
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

        // If we did not add all the items, we cannot topologically sort the vector.
        if topologically_sorted_order.len() != order.len() {
            error!(
                ?order_specific_dependencies,
                "order cannot be topologically sorted, a cycle must exist"
            );
            return None;
        }

        debug!(?topologically_sorted_order, "topological sort complete");
        Some(topologically_sorted_order)
    }
}
