package settings

var (
	DisableSecure      = false
	Port               = 80
	Debug              = false
	DBConnectionString = "postgresql://postgres@127.0.0.1:5432?sslmode=disable"
	TLSCertPath        = "streepjes.pem"
	TLSKeyPath         = "key.pem"
)
