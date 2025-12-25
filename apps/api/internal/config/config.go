package config

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary       Primary              `koanf:"primary" validate:"required"`
	Server        ServerConfig         `koanf:"server" validate:"required"`
	Database      DatabaseConfig       `koanf:"database" validate:"required"`
	Auth          AuthConfig           `koanf:"auth" validate:"required"`
	Redis         RedisConfig          `koanf:"redis" validate:"required"`
	SMTP          SMTPConfig           `koanf:"smtp" validate:"required"`
	Observability *ObservabilityConfig `koanf:"observability"`
	Seeder        SeederConfig         `koanf:"seeder" validate:"required"`
}

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Primary struct {
	Env Env `koanf:"env" validate:"required,oneof=development staging production"`
}

type ServerConfig struct {
	Port               string        `koanf:"port" validate:"required"`
	ReadTimeout        time.Duration `koanf:"read_timeout" validate:"required"`
	WriteTimeout       time.Duration `koanf:"write_timeout" validate:"required"`
	IdleTimeout        time.Duration `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string      `koanf:"cors_allowed_origins" validate:"required"`
}

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
	SSLModeRequire SSLMode = "require"
	SSLModeVerify  SSLMode = "verify-full"
)

type DatabaseConfig struct {
	Host            string        `koanf:"host" validate:"required"`
	Port            int           `koanf:"port" validate:"required"`
	User            string        `koanf:"user" validate:"required"`
	Password        string        `koanf:"password"`
	Name            string        `koanf:"name" validate:"required"`
	SSLMode         SSLMode       `koanf:"ssl_mode" validate:"required,oneof=disable require verify-full"`
	MaxOpenConns    int           `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int           `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time" validate:"required"`
}

type SeederConfig struct {
	Enabled bool          `koanf:"enabled" validate:"required"`
	Timeout time.Duration `koanf:"timeout" validate:"required"`
}

type RedisConfig struct {
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	SMTP SMTPConfig `koanf:"smtp" validate:"required"`
}

type SMTPConfig struct {
	Host      string `koanf:"host" validate:"required"`
	Port      int    `koanf:"port" validate:"required"`
	Username  string `koanf:"username" validate:"required"`
	Password  string `koanf:"password" validate:"required"`
	FromEmail string `koanf:"from_email" validate:"required,email"`
	FromName  string `koanf:"from_name" validate:"required"`
}

type AuthConfig struct {
	SecretKey            string        `koanf:"secret_key" validate:"required"`
	AccessTokenTTL       time.Duration `koanf:"access_token_ttl"`
	RefreshTokenTTL      time.Duration `koanf:"refresh_token_ttl"`
	GoogleClientID       string        `koanf:"google_client_id"`
	EmailVerificationTTL time.Duration `koanf:"email_verification_ttl"`
	AccessCookieName     string        `koanf:"access_cookie_name"`
	RefreshCookieName    string        `koanf:"refresh_cookie_name"`
	CookieDomain         string        `koanf:"cookie_domain"`
	CookieSameSite       string        `koanf:"cookie_same_site"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("API_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "API_"))
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	if mainConfig.Auth.EmailVerificationTTL == 0 {
		mainConfig.Auth.EmailVerificationTTL = 24 * time.Hour
	}
	if mainConfig.Auth.RefreshTokenTTL == 0 {
		mainConfig.Auth.RefreshTokenTTL = 30 * 24 * time.Hour
	}
	if mainConfig.Auth.AccessCookieName == "" {
		mainConfig.Auth.AccessCookieName = "access_token"
	}
	if mainConfig.Auth.RefreshCookieName == "" {
		mainConfig.Auth.RefreshCookieName = "refresh_token"
	}
	if mainConfig.Auth.CookieSameSite == "" {
		mainConfig.Auth.CookieSameSite = "lax"
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	// Set default observability config if not provided
	if mainConfig.Observability == nil {
		mainConfig.Observability = DefaultObservabilityConfig()
	}

	// Override service name and environment from primary config
	mainConfig.Observability.ServiceName = "go-kickstart"
	mainConfig.Observability.Env = mainConfig.Primary.Env

	// Validate observability config
	if err := mainConfig.Observability.Validate(); err != nil {
		logger.Fatal().Err(err).Msg("invalid observability config")
	}

	return mainConfig, nil
}
