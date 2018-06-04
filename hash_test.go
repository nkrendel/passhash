package main

import (
    "testing"
    "time"
)

const IncorrectHashMessage = "Hash was incorrect, got: %s, want: %s."
const IncorrectCounterMessage = "Counter has incorrect value, got: %d, want: %d."

func TestHash(t *testing.T) {
    hashPassword("angryMonkey", 13)

    time.Sleep(5 * time.Second) // wait for calculation

    expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
    if hashMap[13] != expected {
        t.Errorf(IncorrectHashMessage, hashMap[13], expected)
    }
}

func TestLongPassHash(t *testing.T) {
    hashPassword("ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", 13)

    time.Sleep(5 * time.Second) // wait for calculation

    expected := "SsGKniC25ry6qor88fVozdDyNspT4Xqpun8jnZLkPur+0HJWlqz2BUTzhTlGQDOixd0mru+nE8jt6HS90IXWyA=="
    if hashMap[13] != expected {
        t.Errorf(IncorrectHashMessage, hashMap[13], expected)
    }
}

func TestId(t *testing.T) {
    validateId(t, "/hash/34", 34)
    validateId(t, "/hash/", 0)
    validateId(t, "/hash", 0)
    validateId(t, "/hash/a", 0)
    validateId(t, "/hash/12345678", 12345678)
}

func validateId(t *testing.T, path string, expected int64) {
    id := id(path)

    if id != expected {
        t.Errorf("Extracted id was incorrect, got: %d, want: %d.", id, expected)
    }
}