package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	config *Config
)

//Config holds environment variables.
type Config struct {
	Port        string
	TokenSecret string

	DbUser     string
	DbPassword string
	DbName     string
	DbPort     string
}

//New constructs config from environment variables.
func New() *Config {
	//Read ENV`s only once
	once.Do(func() {
		if modeEnv := os.Getenv("MODE"); modeEnv == "" || modeEnv == "development" || modeEnv == "dev" {
			dev()
		}
		config = &Config{
			Port:        getEnv("PORT"),
			TokenSecret: getEnv("TOKEN_SECRET"),
			DbUser:      getEnv("DB_USER"),
			DbPassword:  getEnv("DB_PASSWORD"),
			DbName:      getEnv("DB_NAME"),
			DbPort:      getEnv("DB_PORT"),
		}

		configJSON, err := json.MarshalIndent(config, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		//Print configuration as JSON into terminal
		log.Println("Application ENVs:", string(configJSON))

	})
	return config
}

//getEnv is an helper function to check existence of passed environment variable and exit if there is no such environment variable.
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
