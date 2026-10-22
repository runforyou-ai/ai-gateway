// Package config 加载与校验启动配置（yaml；开发期可用环境变量覆盖）。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// EnvDevelopment 允许环境变量覆盖 yaml。
	EnvDevelopment = "development"
	// EnvProduction 忽略环境变量，只认 yaml。
	EnvProduction = "production"

	appDirName = "lite-ai-gateway"
	configName = "config.yaml"
	pidName    = "lite-ai-gateway.pid"
)

// Config 是进程启动配置。
type Config struct {
	Env      string         `yaml:"env"`
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	NATS     NATSConfig     `yaml:"nats"`
	Log      LogConfig      `yaml:"log"`
	Auth     AuthConfig     `yaml:"auth"`
}

// ServerConfig 是 HTTP 服务配置。
type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	PublicURL string `yaml:"public_url"`
}

// PostgresConfig 是数据库连接配置。
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

// NATSConfig 是 NATS 连接；URL 为空时嵌入 nats-server。
type NATSConfig struct {
	URL string `yaml:"url"`
}

// LogConfig 是日志目录；Dir 为空时按 env 使用默认路径。
type LogConfig struct {
	Dir string `yaml:"dir"`
}

// AuthConfig 是认证相关启动配置。
type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

// Default 返回 init 写入的默认配置（开发向）。
func Default() Config {
	return Config{
		Env: EnvDevelopment,
		Server: ServerConfig{
			Host:      "127.0.0.1",
			Port:      3000,
			PublicURL: "http://127.0.0.1:3000",
		},
		Postgres: PostgresConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			User:     "lite-ai-gateway",
			Password: "lite-ai-gateway",
			Database: "lite-ai-gateway",
			SSLMode:  "disable",
		},
		NATS: NATSConfig{},
		Log:  LogConfig{},
		Auth: AuthConfig{
			JWTSecret: "change-me-in-production",
		},
	}
}

// DefaultPath 返回 XDG 用户配置路径下的默认 config.yaml。
func DefaultPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configName), nil
}

// ConfigDir 返回配置目录（XDG）。
func ConfigDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, appDirName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".config", appDirName), nil
}

// DefaultPIDPath 返回与配置同目录的 pid 文件路径。
func DefaultPIDPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, pidName), nil
}

// Addr 返回 host:port。
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsDevelopment 表示当前为开发环境。
func (c Config) IsDevelopment() bool {
	return strings.EqualFold(c.Env, EnvDevelopment)
}

// LogDir 返回生效的日志目录。
func (c Config) LogDir() string {
	if strings.TrimSpace(c.Log.Dir) != "" {
		return c.Log.Dir
	}
	if c.IsDevelopment() {
		return filepath.Join("data", "logs")
	}
	return "/var/log/" + appDirName
}

// PostgresDSN 返回 bun/pgdriver 可用的 DSN。
func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
		c.Postgres.SSLMode,
	)
}

// Load 从 path 读取 yaml；开发环境可用环境变量覆盖缺省项。
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	if cfg.Env == "" {
		cfg.Env = EnvDevelopment
	}
	if cfg.IsDevelopment() {
		applyEnvOverrides(&cfg)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate 校验必要字段。
func (c Config) Validate() error {
	switch strings.ToLower(c.Env) {
	case EnvDevelopment, EnvProduction:
	default:
		return fmt.Errorf("invalid env %q: want %s or %s", c.Env, EnvDevelopment, EnvProduction)
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port must be positive")
	}
	if c.Postgres.Host == "" || c.Postgres.Database == "" || c.Postgres.User == "" {
		return fmt.Errorf("postgres host, database and user are required")
	}
	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		return fmt.Errorf("auth.jwt_secret is required")
	}
	return nil
}

// Write 将配置写成 yaml 文件，按需创建目录。
func Write(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := os.Getenv("SERVER_PUBLIC_URL"); v != "" {
		cfg.Server.PublicURL = v
	}
	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		cfg.Postgres.Host = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Postgres.Port = p
		}
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		cfg.Postgres.User = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		cfg.Postgres.Password = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		cfg.Postgres.Database = v
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		cfg.Postgres.SSLMode = v
	}
	if v := os.Getenv("NATS_URL"); v != "" {
		cfg.NATS.URL = v
	}
	if v := os.Getenv("LOG_DIR"); v != "" {
		cfg.Log.Dir = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
}
