<?php

namespace TestNamespace2;

use TestNamespace1\TestClass;
use function TestNamespace1\testFunction;

TestClass::CONST1;

$instance1 = new TestClass;
$instance1->method1000();

testFunction();
