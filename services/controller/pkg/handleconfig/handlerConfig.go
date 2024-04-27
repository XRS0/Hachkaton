package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Json struct {
	Port      string `json:"port"`
	IsRunning bool   `json:"isRunning"`
}

func ChangeJson(act, serv, port string) {
	pathToJson := "../config/config.json"

	cfgFile, err := os.ReadFile(pathToJson)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var ExistingJson map[string]Json
	if err := json.Unmarshal(cfgFile, &ExistingJson); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	newJson := make(map[string]Json)
	for key, value := range ExistingJson {
		if strings.HasPrefix(key, "m") {
			newJson[key] = value
		}
	}

	key := "m" + serv

	switch act {
	case "add":
		for _, value := range newJson {
			if value.Port == port {
				log.Fatal("ERROR: Port is already in use")
			}
		}
		newJson[key] = Json{Port: port, IsRunning: false}
	case "start":
		newJson[key] = Json{Port: port, IsRunning: true}
	case "stop":
		newJson[key].Port = ""
	case "chport":
		newJson[serv] = Json{Port: port, IsRunning: false}
	}

	jsondata, err := json.MarshalIndent(newJson, "", " ")
	if err != nil {
		log.Fatalf("ERROR: Troubles with parsing JSON: %v", err)
	}

	if err := os.WriteFile(pathToJson, jsondata, 0644); err != nil {
		log.Fatalf("ERROR: error writing config file: %v", err)
	}
	fmt.Println(newJson)
	fmt.Println("Json is fine!")
}
