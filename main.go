package main

import (
    "fmt"
    "sync"
    "net/http"
    "strconv"
    "crypto/sha512"
    "encoding/base64"
    "strings"
    "time"
    "encoding/json"
    "context"
    "log"
    "os"
)

type Configuration struct {
    Port int `json:"port"`
}

const ConfigFile = "config.json"
const LogFile = "passhash.log"
const DefaultPort = 8081

var counter int
var mutex = &sync.Mutex{}
var hashMap = map[int]string{}
var durationMap = map[int]int64{}
var shuttingDown = false
var wg sync.WaitGroup
var server *http.Server
var configuration Configuration
var logger *log.Logger

func incrementCounter() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}

func main() {
    loadConfiguration()
    writer, err := os.Create(LogFile); if err != nil { panic(err) }
    logger = log.New(writer, "", log.LstdFlags)
    defer writer.Close()

    server = &http.Server{Addr: ":"+strconv.Itoa(configuration.Port)}

    http.HandleFunc("/hash", hashHandler)
    http.HandleFunc("/hash/", hashHandler) // in order to support /{id} path parameter
    http.HandleFunc("/stats", statsHandler)
    http.HandleFunc("/shutdown", shutdownHandler)

    logger.Printf("Starting server on port %d...\n", configuration.Port)
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        // Error starting or closing listener:
        logger.Printf("HTTP server ListenAndServe: %v", err)
    }
}

func loadConfiguration() {
    configuration.Port = DefaultPort

    //filename is the path to the json config file
    file, err := os.Open(ConfigFile); if err != nil {
        logger.Printf("Unable to open config file: %s", ConfigFile)
        return
    }
    decoder := json.NewDecoder(file)
    if err = decoder.Decode(&configuration); err != nil {
        logger.Printf("Error loading configuration file: %v", err)
    }
}

type Stats struct {
    Total   int   `json:"total"`
    Average int64 `json:"average"`
}

func statsHandler(writer http.ResponseWriter, request *http.Request) {
    if shuttingDown {
        fmt.Fprintln(writer, "Server is shutting down...")
        return
    }

    var stats Stats
    stats.Total = counter
    stats.Average = average()
    if bytes, err := json.Marshal(stats); err != nil {
        fmt.Fprintf(writer, "Error serializing..."+err.Error())
    } else {
        fmt.Fprintln(writer, string(bytes))
    }
}

func hashHandler(writer http.ResponseWriter, request *http.Request) {
    if shuttingDown {
        fmt.Fprintln(writer, "Server is shutting down...")
        return
    }

    if request.Method == http.MethodPost {
        handlePost(writer, request)
    } else if request.Method == http.MethodGet {
        handleGet(writer, request)
    } else {
        fmt.Fprintf(writer, "Unsupported method: %s\n", request.Method)
    }
}

func shutdownHandler(writer http.ResponseWriter, request *http.Request) {
    writer.WriteHeader(200)
    fmt.Fprintln(writer, "Shutting down...")

    if !shuttingDown {
        go func() {
            mutex.Lock()
            shuttingDown = true
            mutex.Unlock()

            // wait for processing to finish and shut down server
            logger.Println("Shutting down server...")
            wg.Wait()
            if err := server.Shutdown(context.Background()); err != nil {
                // Error from closing listeners, or context timeout:
                logger.Printf("HTTP server Shutdown: %v", err)
            }
        }()
    }
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
    request.ParseForm()
    var password = request.FormValue("password")
    if password != "" {
        incrementCounter()
        go hashPassword(password, counter)
        fmt.Fprintln(writer, strconv.Itoa(counter))
    }
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
    id := id(request.URL.EscapedPath())
    fmt.Fprintln(writer, hashMap[int(id)])
}

func hashPassword(password string, counter int) string {
    wg.Add(1)     // add to waitgroup
    defer wg.Done()     // indicate done

    startTime := time.Now()     // save start time
    time.Sleep(5 * time.Second) // wait 5 seconds

    // take SHA512 sum of passed-in string
    sha := sha512.New()
    sha.Write([]byte(password))
    hash := sha.Sum(nil)

    // encode hash in base64
    encoded := base64.StdEncoding.EncodeToString(hash)
    hashMap[counter] = encoded

    duration := time.Since(startTime)
    durationMap[counter] = duration.Nanoseconds()

    // store hash
    return string(encoded)
}

func average() int64 {
    var total int64
    for _, v := range durationMap {
        total += v
    }

    if counter == 0 {
        return 0
    } else {
        return total / int64(counter)
    }
}

func id(path string) int64 {
    v := strings.SplitAfter(path, "/")
    if len(v) >= 3 {
        if i, err := strconv.ParseInt(v[2], 10, 64); err == nil {
            return i
        } else {
            return 0
        }
        return 0
    }
    return 0
}
