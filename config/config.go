package config

import "os"

type Config struct {
	Pg Pg
	Clickhouse Clickhouse
	Kafka Kafka
	Redis Redis
}

type Pg struct {
	PgUrl string
	PgMigrationPath string
}

type Clickhouse struct {
	Database string
	Port string
	Host string
}

type Kafka struct {
	Host string
	Port string
	Topic string
}

type Redis struct {
	Host string
	Port string
}

func Get() *Config{
	migPath, _ := os.LookupEnv("PG_MIGRATIONS_PATH")
	return &Config{
		Pg: Pg{
			PgUrl:  "postgres://postgres:postgres@postgres/postgres?sslmode=disable",
			PgMigrationPath: migPath,
		},
		Clickhouse: Clickhouse{
			Host: "clickhouse",
			Port: "9000",
			Database: "logs",
		},
		Kafka: Kafka{
			Host: "kafka",
			Port: "9092",
			Topic: "addUserLogs",
		},
		Redis: Redis{
			Port: "6379",
			Host: "redis",
		},
	}
}