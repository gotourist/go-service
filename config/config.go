package config

import (
	"github.com/joho/godotenv"
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	PostgresSSL      string

	RPCPort string

	KafkaHost string
	KafkaPort int
}

func Load() Config {
	_ = godotenv.Load()

	config := Config{}

	config.PostgresHost = cast.ToString(getOrReturnDefault("DATABASE_HOST", "localhost"))
	config.PostgresPort = cast.ToInt(getOrReturnDefault("DATABASE_PORT", 5432))
	config.PostgresDatabase = cast.ToString(getOrReturnDefault("DATABASE_NAME", "postgres"))
	config.PostgresUser = cast.ToString(getOrReturnDefault("DATABASE_USER", "postgres"))
	config.PostgresPassword = cast.ToString(getOrReturnDefault("DATABASE_PASSWORD", "postgres"))
	config.PostgresSSL = cast.ToString(getOrReturnDefault("DATABASE_SSL", "disable"))

	config.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":50051"))

	config.KafkaHost = cast.ToString(getOrReturnDefault("KAFKA_HOST", "localhost"))
	config.KafkaPort = cast.ToInt(getOrReturnDefault("KAFKA_PORT", 9092))

	return config
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
