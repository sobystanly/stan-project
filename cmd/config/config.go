package config

import "os"

// Global defines the global configuration values
var Global = struct {
	PostgresAddress  string
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string
}{
	PostgresAddress:  getEnv("POSTGRES_ADDRESS", "localhost"),
	PostgresUsername: "postgres",
	PostgresPassword: "postgres",
	PostgresDatabase: "risks",
}

func getEnv(key, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}
	return value
}
