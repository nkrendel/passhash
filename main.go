package main

import (
    "sync"
    "net/http"
    "strconv"
    "log"
    "os"
)

const LogFile = "passhash.log"

var mutex = &sync.Mutex{}
var shuttingDown = false
var wg sync.WaitGroup
var server *http.Server
var logger *log.Logger

func main() {
    LoadConfiguration() // load configuration, currently only port

    // create/overwrite log file
    writer, err := os.Create(LogFile)
    if err != nil {
        panic(err)
    }
    logger = log.New(writer, "", log.LstdFlags)
    defer writer.Close()

    server = &http.Server{Addr: ":" + strconv.Itoa(configuration.Port)}

    http.HandleFunc("/hash", HashHandler)
    http.HandleFunc("/hash/", HashHandler) // in order to support /{id} path parameter
    http.HandleFunc("/stats", StatsHandler)
    http.HandleFunc("/shutdown", ShutdownHandler)

    logger.Printf("Starting server on port %d...\n", configuration.Port)
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        // Error starting or closing listener:
        logger.Printf("HTTP server ListenAndServe: %v", err)
    }
}
