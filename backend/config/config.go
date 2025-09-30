package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	// Database configuration
	DBPassword   string
	DBUser       string
	DBName       string
	DBHost       string
	DBConnString string

	// Test database configuration
	TestDBPassword   string
	TestDBUser       string
	TestDBName       string
	TestDBHost       string
	TestDBConnString string

	// Application configuration
	JWTSecret string
	Port      string
	Env       Environment
	IsCI      bool
}

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// findProjectRoot walks up the directory tree to find the .env file
// This handles the case where tests run from nested directories
func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return ""
}

// GetConfig returns the singleton config instance
func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// ResetForTesting resets the config singleton for testing
// This should only be called from tests
func ResetForTesting() {
	mu.Lock()
	defer mu.Unlock()
	instance = nil
	once = sync.Once{}
}

// loadConfig loads configuration with automatic .env file detection
func loadConfig() *Config {
	// Skip .env loading in CI environment
	if os.Getenv("CI") == "true" {
		return loadFromEnvironment()
	}

	// Find and load .env file
	envPath := findProjectRoot()
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			// In production, .env file might not exist - that's okay
			// Only panic if we're clearly in development (local/testing env)
			if env := os.Getenv("ENV"); env == string(Local) || env == string(LocalTest) {
				panic(fmt.Errorf("failed to load .env file at %s: %w", envPath, err))
			}
		}
	}

	return loadFromEnvironment()
}

// loadFromEnvironment loads all config values from environment variables
func loadFromEnvironment() *Config {
	config := &Config{}
	missing := []string{}

	// Required variables
	requiredVars := map[ConfigKey]*string{
		DB_PASSWORD: &config.DBPassword,
		DB_USER:     &config.DBUser,
		DB_NAME:     &config.DBName,
		JWT_SECRET:  &config.JWTSecret,
		PORT:        &config.Port,
		ENV:         (*string)(&config.Env),
	}

	// Load required variables
	for key, field := range requiredVars {
		if value := os.Getenv(string(key)); value != "" {
			*field = value
		} else {
			missing = append(missing, string(key))
		}
	}

	// Validate environment
	if config.Env != "" && !IsValidAppEnvironment(config.Env) {
		panic(fmt.Errorf("invalid environment: %s", config.Env))
	}

	// Load optional variables
	config.DBHost = os.Getenv(string(DB_HOST))
	config.DBConnString = os.Getenv(string(DB_CONN_STRING))
	config.TestDBPassword = os.Getenv(string(TEST_DB_PASSWORD))
	config.TestDBUser = os.Getenv(string(TEST_DB_USER))
	config.TestDBName = os.Getenv(string(TEST_DB_NAME))
	config.TestDBHost = os.Getenv(string(TEST_DB_HOST))
	config.TestDBConnString = os.Getenv(string(TEST_DB_CONN_STRING))
	config.IsCI = os.Getenv(string(CI)) == "true"

	// Check for missing required variables
	if len(missing) > 0 {
		panic(fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", ")))
	}

	return config
}

func (c *Config) GetDBConnectionString() string {
	if c.DBConnString != "" {
		return c.DBConnString
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		c.DBUser, c.DBPassword, c.DBName, c.DBHost)
}

func (c *Config) GetTestDBConnectionString() string {
	if c.TestDBConnString != "" {
		return c.TestDBConnString
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		c.TestDBUser, c.TestDBPassword, c.TestDBName, c.TestDBHost)
}

func (c *Config) IsDevelopment() bool {
	return c.Env == Local || c.Env == LocalTest
}

func (c *Config) IsProduction() bool {
	return c.Env == Production
}

func (c *Config) IsTestEnvironment() bool {
	return c.Env == Local || c.Env == LocalTest
}

// GetEnv safely gets an environment variable
func GetEnv(key ConfigKey) (string, error) {
	val := os.Getenv(string(key))
	if val == "" {
		return "", fmt.Errorf("environment variable %s not found", key)
	}
	return val, nil
}

// MustGetEnv gets an environment variable or panics
func MustGetEnv(key ConfigKey) string {
	val := os.Getenv(string(key))
	if val == "" {
		panic(fmt.Sprintf("environment variable %s not found", key))
	}
	return val
}

// SetEnv sets an environment variable (primarily for testing)
func SetEnv(key ConfigKey, value string) {
	os.Setenv(string(key), value)
}
