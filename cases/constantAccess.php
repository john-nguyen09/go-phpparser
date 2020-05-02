<?php

__DIR__;

$currentDir = __DIR__;
$randomInclude = include(__DIR__ . '/random.php');

if (TEST_CONST1) {
    echo __DIR__;
}
