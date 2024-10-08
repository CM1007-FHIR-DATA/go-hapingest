package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Config struct {
	DataDir       string
	URLBase       string
	FHIRServerURL string
	Port          string
	PingServer    bool
}

const (
	defaultDataDir       = "./fhir-data"
	defaultURLBase       = "http://host.docker.internal"
	defaultFHIRServerURL = "http://host.docker.internal:8080/fhir"
	defaultPort          = "8001"
	defaultPingServer    = true
)

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			DataDir:       getEnv("DATA_DIR", defaultDataDir),
			URLBase:       getEnv("URL_BASE", defaultURLBase),
			FHIRServerURL: getEnv("FHIR_SERVER_URL", defaultFHIRServerURL),
			Port:          getEnv("PORT", defaultPort),
			PingServer:    getEnvBool("BLOCKING_PING_FHIR_SERVER", defaultPingServer),
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value) == "true"
	}
	return fallback
}

func (c *Config) String() string {
	pingServerStr := "true"
	if !c.PingServer {
		pingServerStr = "false"
	}
	return fmt.Sprintf("Config:\n\tDataDir: %s\n\tURLBase: %s\n\tFHIRServerURL: %s\n\tPort: %s\n\tBlockingPingFhirServer: %s\n\n",
		c.DataDir, c.URLBase, c.FHIRServerURL, c.Port, pingServerStr)
}
