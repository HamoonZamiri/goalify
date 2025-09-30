package config

import "slices"

// Constant keys representing environment variables
type ConfigKey string

const (
	DB_PASSWORD         ConfigKey = "DB_PASSWORD"
	DB_USER             ConfigKey = "DB_USER"
	DB_NAME             ConfigKey = "DB_NAME"
	DB_HOST             ConfigKey = "DB_HOST"
	DB_CONN_STRING      ConfigKey = "DATABASE_URL"
	JWT_SECRET          ConfigKey = "JWT_SECRET"
	PORT                ConfigKey = "PORT"
	TEST_DB_CONN_STRING ConfigKey = "TEST_DB_CONN_STRING"
	CI                  ConfigKey = "CI"
	TEST_DB_NAME        ConfigKey = "TEST_DB_NAME"
	TEST_DB_PASSWORD    ConfigKey = "TEST_DB_PASSWORD"
	TEST_DB_USER        ConfigKey = "TEST_DB_USER"
	TEST_DB_HOST        ConfigKey = "TEST_DB_HOST"
	ENV                 ConfigKey = "ENV"
)

type Environment string

const (
	Staging    Environment = "staging"
	Test       Environment = "test"
	Production Environment = "prod"
	Local      Environment = "local"
	LocalTest  Environment = "localtest"
	CIEnv      Environment = "ci"
)

var AppEnvironments = []Environment{Staging, Test, Production, Local, LocalTest, CIEnv}

func IsValidAppEnvironment(env Environment) bool {
	return slices.Contains(AppEnvironments, env)
}
