package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func InitializeConfig(t *Tunnel) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Error getting config directory: ", err.Error())
	}

	configDir += "/elyntra"
	jsonConfigFilePath := configDir + "/elyntra.config.json"

	if _, err = os.Stat(jsonConfigFilePath); os.IsNotExist(err) {
		log.Printf("No config found, created new config at %s", configDir)

		if _, err = os.Stat(configDir); os.IsNotExist(err) {
			err := os.Mkdir(configDir, 0755)
			if err != nil {
				log.Fatal("Error creating config directory: ", err.Error())
			}
		}
		res, err := http.Get("https://raw.githubusercontent.com/v1noid/elyntra/refs/heads/main/config/elyntra.config.json")

		if err != nil {
			log.Fatal("Error getting default config: ", err.Error())
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatal("Error getting default config: received non-200 response code")
		}
		file, err := os.Create(jsonConfigFilePath)
		if err != nil {
			log.Fatal("Error creating default config file: ", err.Error())
		}
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		if err != nil {
			log.Fatal("Error creating default config file: ", err.Error())
		}
	}
	config, err := os.ReadFile(jsonConfigFilePath)
	if err != nil {
		log.Fatalf("Error reading config file")
	}

	t.Config = &Config{
		Tunnel: make(map[string]ConfigTunnel, 4),
	}
	tmpConfig := struct {
		Tunnel []ConfigTunnel `json:"tunnel"`
	}{}
	err = json.Unmarshal(config, &tmpConfig)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range tmpConfig.Tunnel {
		t.Config.Tunnel[c.Host] = c
	}

}
