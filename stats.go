package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type Stats struct {
    Total   int64   `json:"total"`
    Average int64   `json:"average"`
}

func StatsHandler(writer http.ResponseWriter, request *http.Request) {
    if shuttingDown {
        fmt.Fprintln(writer, "Server is shutting down...")
        return
    }

    logger.Printf("Received /stats %s request\n", request.Method)
    var stats = Stats{int64(len(durationMap)), average()}
    if bytes, err := json.Marshal(stats); err != nil {
        fmt.Fprintf(writer, "Error serializing..."+err.Error())
    } else {
        fmt.Fprintln(writer, string(bytes))
    }
}

func average() int64 {
    var total int64
    for _, v := range durationMap {
        total += v
    }

    if counter == 0 {
        return 0
    } else {
        if len(durationMap) == 0 {
            return 0
        } else {
            return total / int64(len(durationMap))
        }
    }
}
