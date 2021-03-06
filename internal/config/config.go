package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	REDIS_ADDR               string
	POSTGRES_USER            string
	POSTGRES_PASS            string
	POSTGRES_DBNAME          string
	MAX_DB_OPEN_CONN         int64
	MAIN_LISTEN_ENDPOINT     string
	ALLOWED_ORIGINS          []string
	REDIS_PASS               string
	LOGFILE                  string
	GIT_SERVER_LOGFILE       string
	PATH_PREFIX              string
	GIT_USER_REPOS_DIR       string
	GIT_USER_PULLRQ_DIR      string
	GIT_SERVER_ENDPOINT      string
	NEWS_SERVER_ENDPOINT     string
	CSRF_SECRET_KEY          string
	COOKIE_EXPIRE_HOURS      int64
	DEFAULT_USER_AVATAR_NAME string
}

// New returns a new Config struct(!!!пароль не задан по дефолту и csrf secret key)
func New() *Config {
	return &Config{
		REDIS_ADDR:           getEnv("REDIS_ADDR", "localhost:6379"),
		POSTGRES_USER:        getEnv("POSTGRES_USER", "codehub_dev"),
		POSTGRES_PASS:        getEnv("POSTGRES_PASS", ""),
		POSTGRES_DBNAME:      getEnv("POSTGRES_DBNAME", "codehub_dev"),
		MAX_DB_OPEN_CONN:     getEnvAsInt("MAX_DB_OPEN_CONN", 10),
		MAIN_LISTEN_ENDPOINT: getEnv("MAIN_LISTEN_ENDPOINT", "127.0.0.1:8080"),
		REDIS_PASS:           getEnv("REDIS_PASS", ""),
		PATH_PREFIX:          getEnv("PATH_PREFIX", "/home/ubuntu/CodeHub/"),
		LOGFILE:              getEnv("LOGFILE", "logfile.log"),
		GIT_SERVER_LOGFILE:   getEnv("GIT_SERVER_LOGFILE", "gitserver_logfile.log"),
		ALLOWED_ORIGINS: getEnvAsSlice("ALLOWED_ORIGINS",
			[]string{"http://code-hub.space"}, ","),
		GIT_USER_REPOS_DIR:       getEnv("GIT_USER_REPOS_DIR", "codehub_repositories"),
		GIT_USER_PULLRQ_DIR:      getEnv("GIT_USER_PULLRQ_DIR", "codehub_pullrequests"),
		CSRF_SECRET_KEY:          getEnv("CSRF_SECRET_KEY", ""),
		COOKIE_EXPIRE_HOURS:      getEnvAsInt("COOKIE_EXPIRE_HOURS", 72),
		GIT_SERVER_ENDPOINT:      getEnv("GIT_SERVER_ENDPOINT", "127.0.0.1:5000"),
		NEWS_SERVER_ENDPOINT:     getEnv("NEWS_SERVER_ENDPOINT", "127.0.0.1:8083"),
		DEFAULT_USER_AVATAR_NAME: getEnv("DEFAULT_USER_AVATAR_NAME", "default.png"),
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
//func getEnvAsBool(name string, defaultVal bool) bool {
//	valStr := getEnv(name, "")
//	if val, err := strconv.ParseBool(valStr); err == nil {
//		return val
//	}
//
//	return defaultVal
//}

//вернет slice
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
