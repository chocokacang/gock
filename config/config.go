package config

import (
	"os"

	"github.com/chocokacang/gock/utils"
)

type Config struct {
	// Application name
	AppName string
	// Application environment LOCAL, TEST, PRODUCTION & etc
	AppEnv string
	// Enable application debug mode
	AppDebug bool

	// Database URL
	DBURL string
	// ORM perform single create, update, delete operations in transactions by default to ensure database data integrity.
	// You can disable it by setting `DBSkipDefaultTransaction` to true
	DBSkipDefaultTransaction bool

	// Write log into file. You can disable it by setting `LogFile` to empty string
	LogFile string
	// Show the log by level, available value is INFO, WARNING and ERROR.
	// If `AppDebug` is true, the log level will automatically set to INFO
	LogLevel string

	// HTTP port default is 8080
	HTTPPort string
	// Enable H2C support
	HTTPH2C bool
}

func Default() Config {
	config := &Config{
		AppEnv:   os.Getenv("APP_ENV"),
		AppName:  os.Getenv("APP_NAME"),
		AppDebug: utils.GetEnvBool("APP_DEBUG", true),
		DBURL:    os.Getenv("DB_URL"),
		LogLevel: utils.GetEnv("LOG_LEVEL", "WARNING"),
		LogFile:  os.Getenv("LOG_FILE"),
		HTTPPort: utils.GetEnv("HTTP_PORT", "8080"),
		HTTPH2C:  utils.GetEnvBool("HTTP_H2C", false),
	}
	if config.AppDebug {
		config.LogFile = "INFO"
	}
	return *config
}
