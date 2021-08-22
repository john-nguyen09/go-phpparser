<?php

fn(bool $a) => $a;
fn($x = 42) => $x;
static fn(&$x) => $x;
fn&($x) => $x;
fn($x, ...$rest) => $rest;
fn(): int => $x;

$fn1 = fn($x) => $x + $y;
