package main

import (
    "testing"
)

func TestAverage(t *testing.T) {
    // first clear the duration map
    for k := range durationMap {
        delete(durationMap, k)
    }

    // add some durations
    incrementCounter()
    if counter != 1 {
        t.Errorf(IncorrectCounterMessage, counter, 1)
    }

    durationMap[1] = 1012
    incrementCounter()
    durationMap[2] = 2034
    incrementCounter()
    durationMap[3] = 536
    incrementCounter()
    durationMap[4] = 9876

    if counter != 4 {
        t.Errorf(IncorrectCounterMessage, counter, 4)
    }

    rc := average()
    expected := (durationMap[1] + durationMap[2] + durationMap[3] + durationMap[4]) / 4
    if rc != expected {
        t.Errorf("Average was incorrect, got: %d, want: %d.", rc, expected)
    }
}
