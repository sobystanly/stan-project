package config

// Global defines the global configuration values
var Global = struct {
	PostgresAddress  string
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string
}{
	PostgresAddress:  "postgres",
	PostgresUsername: "postgres",
	PostgresPassword: "postgres",
	PostgresDatabase: "risks",
}
