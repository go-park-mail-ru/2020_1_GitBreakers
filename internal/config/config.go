package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	REDIS_ADDR          string
	POSTGRES_USER       string
	POSTGRES_PASS       string
	POSTGRES_DBNAME     string
	MAX_DB_OPEN_CONN    int64
	MAIN_LISTEN_PORT    string
	ALLOWED_ORIGINS     []string
	REDIS_PASS          string
	LOGFILE             string
	HOST_TO_SAVE        string
	GIT_USER_REPOS_DIR  string
	GIT_SERVER_PORT     string
	CSRF_SECRET_KEY     string
	COOKIE_EXPIRE_HOURS int64
}

// New returns a new Config struct(!!!пароль не задан по дефолту и csrf secret key)
func New() *Config {
	return &Config{
		REDIS_ADDR:       getEnv("REDIS_ADDR", "localhost:6379"),
		POSTGRES_USER:    getEnv("POSTGRES_USER", "codehub_dev"),
		POSTGRES_PASS:    getEnv("POSTGRES_PASS", ""),
		POSTGRES_DBNAME:  getEnv("POSTGRES_DBNAME", "codehub_dev"),
		MAX_DB_OPEN_CONN: getEnvAsInt("MAX_DB_OPEN_CONN", 10),
		MAIN_LISTEN_PORT: getEnv("MAIN_LISTEN_PORT", ":8080"),
		REDIS_PASS:       getEnv("REDIS_PASS", ""),
		HOST_TO_SAVE:     getEnv("HOST_TO_SAVE", "http://89.208.198.186:8080"),
		LOGFILE:          getEnv("LOGFILE", "logfile.log"),
		ALLOWED_ORIGINS: getEnvAsSlice("ALLOWED_ORIGINS",
			[]string{"http://localhost:3000", "http://89.208.198.186:80", "http://89.208.198.186:3000"}, ","),
		GIT_USER_REPOS_DIR:  getEnv("GIT_USER_REPOS_DIR", "codehub_repositories"),
		CSRF_SECRET_KEY:     getEnv("CSRF_SECRET_KEY", ""),
		COOKIE_EXPIRE_HOURS: getEnvAsInt("COOKIE_EXPIRE_HOURS", 72),
		GIT_SERVER_PORT:     getEnv("GIT_SERVER_PORT", ":5000"),
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
func getEnvAsInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return int64(value)
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
