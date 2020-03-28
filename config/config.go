package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	REDIS_ADDR       string
	POSTGRES_USER    string
	POSTGRES_PASS    string
	POSTGRES_DBNAME  string
	MAX_DB_OPEN_CONN int
	MAIN_LISTEN_PORT string
	ALLOWED_ORIGINS  []string
}

// New returns a new Config struct(!!!пароль не задан по дефолту)
func New() *Config {
	return &Config{
		REDIS_ADDR:       getEnv("REDIS_ADDR", "localhost:6379"),
		POSTGRES_USER:    getEnv("POSTGRES_USER", "codehub_dev"),
		POSTGRES_PASS:    getEnv("POSTGRES_PASS", ""),
		POSTGRES_DBNAME:  getEnv("POSTGRES_DBNAME", "codehub_dev"),
		MAX_DB_OPEN_CONN: getEnvAsInt("MAX_DB_OPEN_CONN", 10),
		MAIN_LISTEN_PORT: getEnv("MAIN_LISTEN_PORT", ":8080"),
		ALLOWED_ORIGINS: getEnvAsSlice("ALLOWED_ORIGINS",
			[]string{"http://89.208.198.186:8080", "http://89.208.198.186:80"}, ","),
	}
}

//функции обертка для получения данных из переменной окружения
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

//вернет переменную окружения int
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// вернет переменную окружения bool
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

//вернет slice
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}