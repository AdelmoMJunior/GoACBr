package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Encryption EncryptionConfig
	B2         B2Config
	CNPJWS     CNPJWSConfig
	ACBr       ACBrConfig
	RateLimit  RateLimitConfig
	Log        LogConfig
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type EncryptionConfig struct {
	Key string // 32-byte hex key for AES-256-GCM
}

type B2Config struct {
	KeyID        string
	AppKey       string
	BucketName   string
	Endpoint     string
	Region       string
	PublicCDNURL string
}

type CNPJWSConfig struct {
	APIURL   string
	APIToken string
}

type ACBrConfig struct {
	LibPath     string
	SchemasPath string
	LogPath     string
	LogLevel    string
}

type RateLimitConfig struct {
	RequestsPerMinute int
	Burst             int
}

type LogConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnv("SERVER_PORT", "8080"),
			ReadTimeout:     getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "goacbr"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "goacbr"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenTTL:  getDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getDuration("JWT_REFRESH_TOKEN_TTL", 168*time.Hour),
		},
		Encryption: EncryptionConfig{
			Key: getEnv("ENCRYPTION_KEY", ""),
		},
		B2: B2Config{
			KeyID:        getEnv("B2_KEY_ID", ""),
			AppKey:       getEnv("B2_APP_KEY", ""),
			BucketName:   getEnv("B2_BUCKET_NAME", "goacbr-storage"),
			Endpoint:     getEnv("B2_ENDPOINT", ""),
			Region:       getEnv("B2_REGION", "us-west-004"),
			PublicCDNURL: getEnv("B2_PUBLIC_CDN_URL", ""),
		},
		CNPJWS: CNPJWSConfig{
			APIURL:   getEnv("CNPJWS_API_URL", "https://publica.cnpj.ws/cnpj"),
			APIToken: getEnv("CNPJWS_API_TOKEN", ""),
		},
		ACBr: ACBrConfig{
			LibPath:     getEnv("ACBR_LIB_PATH", "/app/lib/libacbrnfe64.so"),
			SchemasPath: getEnv("ACBR_SCHEMAS_PATH", "/app/lib/Schemas/NFe"),
			LogPath:     getEnv("ACBR_LOG_PATH", "/app/logs"),
			LogLevel:    getEnv("ACBR_LOG_LEVEL", "4"),
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: getInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
			Burst:             getInt("RATE_LIMIT_BURST", 10),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required fields.
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if cfg.Encryption.Key == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required")
	}

	return cfg, nil
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// Helper functions for reading environment variables.

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
