package main

import (
    "fmt"
    "net/http"
    "strconv"
    "crypto/sha512"
    "encoding/base64"
    "time"
    "strings"
)

var counter int64
var hashMap = map[int64]string{}
var durationMap = map[int64]int64{}

func incrementCounter() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}

func HashHandler(writer http.ResponseWriter, request *http.Request) {
    if shuttingDown {
        fmt.Fprintln(writer, "Server is shutting down...")
        return
    }

    logger.Printf("Received /hash %s request\n", request.Method)
    if request.Method == http.MethodPost {
        HandlePost(writer, request)
    } else if request.Method == http.MethodGet {
        HandleGet(writer, request)
    } else {
        fmt.Fprintf(writer, "Unsupported method: %s\n", request.Method)
    }
}

func HandlePost(writer http.ResponseWriter, request *http.Request) {
    request.ParseForm()
    var password = request.FormValue("password")
    if password != "" {
        incrementCounter()
        go hashPassword(password, counter)
        fmt.Fprintln(writer, strconv.FormatInt(counter, 10))
    } else {
        fmt.Fprintln(writer, "password to hash not provided")
    }
}

func HandleGet(writer http.ResponseWriter, request *http.Request) {
    id := id(request.URL.EscapedPath())
    if hashMap[id] != "" {
        fmt.Fprintln(writer, hashMap[id])
    }
}

func hashPassword(password string, counter int64) string {
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
