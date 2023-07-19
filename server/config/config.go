package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/pelletier/go-toml"

	"github.com/rs/zerolog/log"
)

type Config struct {
	DB                 *DBConfig
	ServiceBindAddress string
	ServicePort        int
	ServiceKey         string
	LoversEar          string
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func DefaultConfig() *Config {

	/// TODO Add support for config from environment variable
	/*
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "8080"
		}
	*/

	// first we try to parse the config at the global configuration path
	if fileExists("/etc/festivals-identity-server.conf") {
		config := ParseConfig("/etc/festivals-identity-server.conf")
		if config != nil {
			return config
		}
	}

	// if there is no global configuration check the current folder for the template config file
	// this is mostly so the application will run in development environment
	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Msg("server initialize: could not read default config file with error:" + err.Error())
	}
	path = path + "/config_template.toml"
	return ParseConfig(path)
}

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal().Msg("server initialize: could not read config file at '" + cfgFile + "' with error: " + err.Error())
	}

	serviceBindAdress := content.Get("service.bind-address").(string)
	servicePort := content.Get("service.port").(int64)
	serviceKey := content.Get("service.key").(string)

	dbHost := content.Get("database.host").(string)
	dbPort := content.Get("database.port").(int64)
	dbUsername := content.Get("database.username").(string)
	dbPassword := content.Get("database.password").(string)
	databaseName := content.Get("database.database-name").(string)

	loversear := content.Get("heartbeat.endpoint").(string)

	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     dbHost,
			Port:     int(dbPort),
			Username: dbUsername,
			Password: dbPassword,
			Name:     databaseName,
			Charset:  "utf8",
		},
		ServiceBindAddress: serviceBindAdress,
		ServicePort:        int(servicePort),
		ServiceKey:         serviceKey,
		LoversEar:          loversear,
	}
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
// see: https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}
