#![allow(dead_code)]

pub type Operation = fn(i64, i64) -> i64;

pub fn addition(a: i64, b:i64) -> i64 {
    a+b
}

pub fn multiplication(a: i64, b:i64) -> i64 {
    a*b
}

pub fn concatenation(a: i64, b:i64) -> i64 {
    a * 10i64.pow(b.ilog10() + 1) + b
    // (a.to_string() + &b.to_string()).parse().unwrap()
}