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
		for keyIntoMap, value := range newJson {
			if value.Port == port {
				log.Fatal("ERROR: Port is already in use")
			}
			if keyIntoMap == key {
				log.Fatal("ERROR: Name of service is already in use")
			}
		}
		newJson[key] = Json{Port: port, IsRunning: false}
	case "start":
		if val, ok := newJson[key]; ok {
			if serv == "bc" {
				for key, _ := range newJson {
					if valueIntoBranching, ok := newJson[key]; ok {
						newJson[key] = Json{Port: valueIntoBranching.Port, IsRunning: true}
					}
				}
			}
			newJson[key] = Json{Port: val.Port, IsRunning: true}
		} else {
			log.Fatalf("ERROR: Server with key %s does not exist", key)
		}
	case "stop":
		if val, ok := newJson[key]; ok {
			newJson[key] = Json{Port: val.Port, IsRunning: false}
		} else {
			log.Fatalf("ERROR: Server with key %s does not exist", key)
		}
	case "chport":
		if val, ok := newJson[serv]; ok {
			newJson[serv] = Json{Port: port, IsRunning: val.IsRunning}
		} else {
			log.Fatalf("ERROR: Server with key %s does not exist", serv)
		}
	case "del":
		if _, exists := newJson[key]; exists {
			delete(newJson, key)
		} else {
			log.Fatalf("ERROR: Server with key %s does not exist", key)
		}
	default:
		log.Fatal("ERROR: Invalid action specified")
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
