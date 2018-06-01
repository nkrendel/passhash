package main

import (
    "fmt"
    "net/http"
    "context"
)

func ShutdownHandler(writer http.ResponseWriter, request *http.Request) {
    logger.Printf("Received /shutdown %s request\n", request.Method)

    writer.WriteHeader(200)
    fmt.Fprintln(writer, "Shutting down...")

    if !shuttingDown {
        go func() {
            mutex.Lock()
            shuttingDown = true
            mutex.Unlock()

            // wait for processing to finish and shut down server
            wg.Wait()
            logger.Println("Shutting down server...")
            if err := server.Shutdown(context.Background()); err != nil {
                // Error from closing listeners, or context timeout:
                logger.Printf("HTTP server Shutdown: %v", err)
            }
        }()
    }
}
