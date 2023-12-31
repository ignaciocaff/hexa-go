package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Default PORT to set if .env PORT var is empty or can't load
const DEFAULT_PORT_IF_EMPTY = "8080"

// Env config struct
type EnvApp struct {
	// Server Envs
	PORT string

	JWT_SECRET string
	// Database Envs
	DB_HOST        string
	DB_PORT        string
	DB_USERNAME    string
	DB_PASSWORD    string
	DB_SERVICE     string
	CIDI_PASS      string
	CIDI_KEY       string
	ID_APP         string
	BASE_CIDI_URI  string
	GIN_MODE       string
	TIMEZONE       string
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
	ENCRYPT_KEY    string
}

// Get the env configuration
func GetEnv(env_file string) EnvApp {
	err := godotenv.Load(env_file)
	if err != nil {
		log.Printf("error loading: %+v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT_IF_EMPTY
	}

	return EnvApp{
		PORT:           port,
		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_USERNAME:    os.Getenv("DB_USERNAME"),
		DB_PASSWORD:    os.Getenv("DB_PASSWORD"),
		DB_SERVICE:     os.Getenv("DB_SERVICE"),
		JWT_SECRET:     os.Getenv("JWTSECRET"),
		CIDI_PASS:      os.Getenv("CIDI_PASS"),
		CIDI_KEY:       os.Getenv("CIDI_KEY"),
		ID_APP:         os.Getenv("ID_APP"),
		BASE_CIDI_URI:  os.Getenv("BASE_CIDI_URI"),
		GIN_MODE:       os.Getenv("GIN_MODE"),
		TIMEZONE:       os.Getenv("TIMEZONE"),
		REDIS_HOST:     os.Getenv("REDIS_HOST"),
		REDIS_PORT:     os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
		ENCRYPT_KEY:    os.Getenv("ENCRYPT_KEY"),
	}
}
