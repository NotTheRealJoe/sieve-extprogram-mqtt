package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	internal "github.com/nottherealjoe/sieve-extprogram-mqtt/internal"
)

func findConfigFile() (string, error) {
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json", nil
	} else if _, err := os.Stat("/etc/sieve-extprogram-mqtt/config.json"); err == nil {
		return "/etc/sieve-extprogram-mqtt/config.json", nil
	} else {
		return "", fmt.Errorf("can't find the credentials file in current directory or at /etc/sieve-extprogram-mqtt/config.json")
	}
}

func main() {
	configFilePath, err := findConfigFile()
	if err != nil {
		panic(err)
	}
	configFile, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	confFileContent, err := io.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	config := internal.Config{}
	json.Unmarshal(confFileContent, &config)

	client, err := internal.MQTTSetup(config)
	if err != nil {
		panic(err)
	}

	email, err := internal.ReadEmail(os.Stdin)
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(email)

	dest, _ := os.Create("./sampleout.json")
	defer dest.Close()
	dest.Write(json)

	internal.MQTTPublish(*client, config.Topic, json)
}
