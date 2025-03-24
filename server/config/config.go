package config

import (
	"os"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/rs/zerolog/log"
)

type Config struct {
	ServiceBindHost           string
	ServicePort               int
	ServiceKey                string
	TLSRootCert               string
	TLSCert                   string
	TLSKey                    string
	LoversEar                 string
	Interval                  int
	JwtExpiration             int
	AccessTokenPrivateKeyPath string
	AccessTokenPublicKeyPath  string
	InfoLog                   string
	TraceLog                  string
	DB                        *DBConfig
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

func ParseConfig(cfgFile string) *Config {

	content, err := toml.LoadFile(cfgFile)
	if err != nil {
		log.Fatal().Err(err).Msg("server initialize: could not read config file at '" + cfgFile + "'")
	}

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

	infoLogPath := content.Get("log.info").(string)
	traceLogPath := content.Get("log.trace").(string)

	dbPassword := content.Get("database.password").(string)

	tlsrootcert = expandTilde(tlsrootcert)
	tlscert = expandTilde(tlscert)
	tlskey = expandTilde(tlskey)
	accessTokenPublicKeyPath = expandTilde(accessTokenPublicKeyPath)
	accessTokenPrivateKeyPath = expandTilde(accessTokenPrivateKeyPath)
	infoLogPath = expandTilde(infoLogPath)
	traceLogPath = expandTilde(traceLogPath)

	return &Config{
		ServiceBindHost:           serviceBindHost,
		ServicePort:               int(servicePort),
		ServiceKey:                serviceKey,
		TLSRootCert:               tlsrootcert,
		TLSCert:                   tlscert,
		TLSKey:                    tlskey,
		LoversEar:                 loversear,
		Interval:                  int(interval),
		JwtExpiration:             int(jwtExpiration),
		AccessTokenPublicKeyPath:  accessTokenPublicKeyPath,
		AccessTokenPrivateKeyPath: accessTokenPrivateKeyPath,
		InfoLog:                   infoLogPath,
		TraceLog:                  traceLogPath,
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     "localhost",
			Port:     int(3306),
			Username: "festivals.identity.writer",
			Password: dbPassword,
			Name:     "festivals_identity_database",
			Charset:  "utf8",
		},
	}
}

// ExpandTilde expands a leading tilde (~) to the user's home directory
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}
