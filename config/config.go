package config

import (
	"github.com/agirot/superviGo/ressource"
	"os"
	"encoding/json"
)

var Config ressource.ConfigurationFile

func HydrateConfiguration() {
	file, err := os.Open("config.json")
	if err != nil {
		panic("config.json not found")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		panic("Failed parsing json config: "+ err.Error())
	}
}
