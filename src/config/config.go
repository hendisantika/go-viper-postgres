package config

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Password PasswordConfig
	Cors     CorsConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	InternalPort string
	ExternalPort string
	RunMode      string
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}
