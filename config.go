package main

import (
    "encoding/json"
    "os"
)

type Configuration struct {
    Port int `json:"port"`
}

const ConfigFile = "config.json"
const DefaultPort = 8081

var configuration Configuration

func LoadConfiguration() {
    configuration.Port = DefaultPort    // establish default value

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
