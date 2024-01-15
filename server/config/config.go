package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/pelletier/go-toml"

	"github.com/rs/zerolog/log"
)

type Config struct {
	ServiceBindAddress         string
	ServiceBindHost            string
	ServicePort                int
	ServiceKey                 string
	TLSRootCert                string
	TLSCert                    string
	TLSKey                     string
	LoversEar                  string
	Interval                   int
	JwtExpiration              int
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	DB                         *DBConfig
}

type DBConfig struct {
	Dialect    string
	Host       string
	Port       int
	Username   string
	Password   string
	ClientCA   string
	ClientCert string
	ClientKey  string
	Name       string
	Charset    string
}

func DefaultConfig() *Config {

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
		log.Fatal().Err(err).Msg("server initialize: could not read config file at '" + cfgFile + "'")
	}

	serviceBindAdress := content.Get("service.bind-address").(string)
	serviceBindHost := content.Get("service.bind-host").(string)
	servicePort := content.Get("service.port").(int64)
	serviceKey := content.Get("service.key").(string)

	tlsrootcert := content.Get("tls.festivaslapp-root-ca").(string)
	tlscert := content.Get("tls.cert").(string)
	tlskey := content.Get("tls.key").(string)

	loversear := content.Get("heartbeat.endpoint").(string)
	interval := content.Get("heartbeat.interval").(int64)

	jwtExpiration := content.Get("jwt.expiration").(int64)
	accessTokenPrivateKeyPath := content.Get("jwt.accessprivatekeypath").(string)
	accessTokenPublicKeyPath := content.Get("jwt.accesspublickeypath").(string)
	refreshTokenPrivateKeyPath := content.Get("jwt.refreshprivatekeypath").(string)
	refreshTokenPublicKeyPath := content.Get("jwt.refreshpublickeypath").(string)

	dbHost := content.Get("database.host").(string)
	dbPort := content.Get("database.port").(int64)
	dbUsername := content.Get("database.username").(string)
	dbPassword := content.Get("database.password").(string)
	dbClientCA := content.Get("database.festivalsapp-root-ca").(string)
	dbClientCert := content.Get("database.cert").(string)
	dbClientKey := content.Get("database.key").(string)

	return &Config{
		ServiceBindAddress:         serviceBindAdress,
		ServiceBindHost:            serviceBindHost,
		ServicePort:                int(servicePort),
		ServiceKey:                 serviceKey,
		TLSRootCert:                tlsrootcert,
		TLSCert:                    tlscert,
		TLSKey:                     tlskey,
		LoversEar:                  loversear,
		Interval:                   int(interval),
		JwtExpiration:              int(jwtExpiration),
		AccessTokenPublicKeyPath:   accessTokenPrivateKeyPath,
		AccessTokenPrivateKeyPath:  accessTokenPublicKeyPath,
		RefreshTokenPrivateKeyPath: refreshTokenPrivateKeyPath,
		RefreshTokenPublicKeyPath:  refreshTokenPublicKeyPath,
		DB: &DBConfig{
			Dialect:    "mysql",
			Host:       dbHost,
			Port:       int(dbPort),
			Username:   dbUsername,
			Password:   dbPassword,
			ClientCA:   dbClientCA,
			ClientCert: dbClientCert,
			ClientKey:  dbClientKey,
			Name:       "festivals_identity_database",
			Charset:    "utf8",
		},
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
