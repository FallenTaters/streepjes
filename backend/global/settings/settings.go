package settings

type Config struct {
	DisableSecure      bool   `mapstructure:"disable_secure"`
	Port               int    `mapstructure:"port"`
	LogLevel           string `mapstructure:"log_level"`
	DBConnectionString string `mapstructure:"db_connection_string"`
	TLSCertPath        string `mapstructure:"tls_cert_path"`
	TLSKeyPath         string `mapstructure:"tls_key_path"`
}

func DefaultConfig() Config {
	return Config{
		Port:               80,
		LogLevel:           "info",
		DBConnectionString: "postgresql://postgres@127.0.0.1:5432?sslmode=disable",
		TLSCertPath:        "streepjes.pem",
		TLSKeyPath:         "key.pem",
	}
}
