<?php

function fetchSomethingLarge() {
    $arr = [1, 2, 3];

    yield from $arr;
}
