package main

import (
    "encoding/json"
    "os"
)

type Configuration struct {
    Port int `json:"port"`
    HashSize int64 `json:"hashSize"`
}

const ConfigFile = "config.json"
const DefaultPort = 8081           // default port to listen on
const DefaultHashSize = 1000000    // default is 1 million entries

var configuration Configuration

func LoadConfiguration() {
    // establish default values
    configuration.Port = DefaultPort
    configuration.HashSize = DefaultHashSize

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
