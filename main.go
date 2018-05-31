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
)

var counter int
var mutex = &sync.Mutex{}
var hashMap = map[int]string{}
var durationMap = map[int]int64{}
var shuttingDown = false
var wg sync.WaitGroup
var server *http.Server

func incrementCounter() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}

func main() {

    server = &http.Server{Addr: ":8081"}

    http.HandleFunc("/hash", hashHandler)
    http.HandleFunc("/hash/", hashHandler) // in order to support /{id} path parameter
    http.HandleFunc("/stats", statsHandler)
    http.HandleFunc("/shutdown", shutdownHandler)

    server.ListenAndServe()
}

type Stats struct {
    Total   int   `json:"total"`
    Average int64 `json:"average"`
}

func statsHandler(writer http.ResponseWriter, request *http.Request) {
    if shuttingDown {
        fmt.Fprintln(writer, "Shutting down...")
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
        fmt.Fprintln(writer, "Shutting down...")
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

            wg.Wait()                               // wait for all processing to finish
            server.Shutdown(context.Background())   // shutdown
        }()
    }
}

func handlePost(writer http.ResponseWriter, request *http.Request) {
    request.ParseForm()
    var password = request.Form.Get("password")
    if password != "" {
        incrementCounter()
        go hashPassword(password, counter)
        fmt.Fprintln(writer, strconv.Itoa(counter))
    }
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
    id := id(request.URL.EscapedPath())
    fmt.Fprintln(writer, hashMap[int(id)])
    fmt.Fprint(writer, "Duration: ")
    fmt.Fprintln(writer, durationMap[int(id)])
}

func hashPassword(password string, counter int) string {
    wg.Add(1)   // add to waitgroup
    defer wg.Done()   // indicate done

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
