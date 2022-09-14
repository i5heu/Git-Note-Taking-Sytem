package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
	GitSSH string `json:"gitSSH"`
}

func GetConfig() Config {
	jsonFile, err := os.Open(GetWorkDir() + "/config.json")
	if err != nil {
		log.Error().Err(err).Msg("open config.json failed")
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	defer jsonFile.Close()

	return config
}

func GetWorkDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("Error getting user home dir")
		panic(err)
	}
	return dirname + "/.Tyche"
}
