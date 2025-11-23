// Package config handles project level config variables
package config

import "slices"

// ConfigKey represents keys tied to environment variables
type ConfigKey string

const (
	DBPassword       ConfigKey = "DB_PASSWORD"
	DBUser           ConfigKey = "DB_USER"
	DBName           ConfigKey = "DB_NAME"
	DBHost           ConfigKey = "DB_HOST"
	DBConnString     ConfigKey = "DATABASE_URL"
	JWTSecret        ConfigKey = "JWT_SECRET"
	Port             ConfigKey = "PORT"
	AllowedOrigins   ConfigKey = "ALLOWED_ORIGINS"
	TestDBConnString ConfigKey = "TEST_DB_CONN_STRING"
	CI               ConfigKey = "CI"
	TestDBName       ConfigKey = "TEST_DB_NAME"
	TestDBPassword   ConfigKey = "TEST_DB_PASSWORD"
	TestDBUser       ConfigKey = "TEST_DB_USER"
	TestDBHost       ConfigKey = "TEST_DB_HOST"
	ENV              ConfigKey = "ENV"
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
