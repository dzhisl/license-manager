package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv" // Import godotenv package
)

// Config holds the application configuration.
type Config struct {
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	AuthData    AuthData   `yaml:"auth_data"`
}

// AuthData holds authentication credentials.
type AuthData struct {
	ApiKey string `env:"API_KEY" env-required:"true"`
}

// HTTPServer holds HTTP server configuration.
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad loads the configuration from the specified path.
func MustLoad() *Config {
	configPath := "config/local.yaml"
	envPath := "config/local.env"

	// Load environment variables from .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found, proceeding with defaults.")
	}

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
