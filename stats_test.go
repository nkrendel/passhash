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
    durationMap[0] = 1012; incrementCounter()
    durationMap[1] = 2034; incrementCounter()
    durationMap[2] = 536; incrementCounter()
    durationMap[3] = 9876; incrementCounter()

    rc := average()
    expected := (durationMap[0] + durationMap[1] + durationMap[2] + durationMap[3]) / 4
    if rc != expected {
        t.Errorf("Average was incorrect, got: %d, want: %d.", rc, expected)
    }
}
