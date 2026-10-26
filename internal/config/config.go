// Package config 加载与校验启动配置。
// init 可用当前目录 .env 覆盖默认值后写入 yaml；serve / reload 只读 yaml。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	// EnvDevelopment 开发环境（影响日志等默认路径）。
	EnvDevelopment = "development"
	// EnvProduction 生产环境（影响日志等默认路径）。
	EnvProduction = "production"

	appDirName = "ai-gateway"
	configName = "config.yaml"
	pidName    = "ai-gateway.pid"
	dotEnvName = ".env"
)

// Config 是进程启动配置。
type Config struct {
	Env      string         `yaml:"env"`
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Log      LogConfig      `yaml:"log"`
	Email    EmailConfig    `yaml:"email"`
}

// EmailConfig 是发信部署配置；仅支持 Resend。未配齐时发信不可用。
type EmailConfig struct {
	From     string       `yaml:"from"`
	FromName string       `yaml:"from_name"`
	Resend   ResendConfig `yaml:"resend"`
}

// ResendConfig 是 Resend 渠道凭证。
type ResendConfig struct {
	APIKey string `yaml:"api_key"`
}

// Configured 表示 from 与 Resend API Key 均已设置。
func (c EmailConfig) Configured() bool {
	return strings.TrimSpace(c.From) != "" && strings.TrimSpace(c.Resend.APIKey) != ""
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

// LogConfig 是日志目录；Dir 为空时按 env 使用默认路径。
type LogConfig struct {
	Dir string `yaml:"dir"`
}

// Default 返回内置默认配置。
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
			User:     "ai-gateway",
			Password: "ai-gateway",
			Database: "ai-gateway",
			SSLMode:  "disable",
		},
		Log: LogConfig{},
	}
}

// ForInit 返回 init 要写入的配置：默认值，若当前目录存在 .env 则用其中变量覆盖。
// 第二个返回值为实际使用的 .env 路径；未使用时为空。
func ForInit() (Config, string, error) {
	cfg := Default()
	if _, err := os.Stat(dotEnvName); err != nil {
		if os.IsNotExist(err) {
			return cfg, "", nil
		}
		return Config{}, "", err
	}
	values, err := godotenv.Read(dotEnvName)
	if err != nil {
		return Config{}, "", fmt.Errorf("read %s: %w", dotEnvName, err)
	}
	applyEnvMap(&cfg, values)
	return cfg, dotEnvName, nil
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

// Load 从 path 读取 yaml；不读取环境变量或 .env。
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
	return nil
}

// Write 将配置写成 yaml 文件，按需创建目录。
func Write(path string, cfg Config) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func applyEnvMap(cfg *Config, env map[string]string) {
	if v := env["APP_ENV"]; v != "" {
		cfg.Env = v
	}
	if v := env["SERVER_HOST"]; v != "" {
		cfg.Server.Host = v
	}
	if v := env["SERVER_PORT"]; v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := env["SERVER_PUBLIC_URL"]; v != "" {
		cfg.Server.PublicURL = v
	}
	if v := env["POSTGRES_HOST"]; v != "" {
		cfg.Postgres.Host = v
	}
	if v := env["POSTGRES_PORT"]; v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Postgres.Port = p
		}
	}
	if v := env["POSTGRES_USER"]; v != "" {
		cfg.Postgres.User = v
	}
	if v := env["POSTGRES_PASSWORD"]; v != "" {
		cfg.Postgres.Password = v
	}
	if v := env["POSTGRES_DB"]; v != "" {
		cfg.Postgres.Database = v
	}
	if v := env["POSTGRES_SSLMODE"]; v != "" {
		cfg.Postgres.SSLMode = v
	}
	if v := env["LOG_DIR"]; v != "" {
		cfg.Log.Dir = v
	}
	if v := env["EMAIL_FROM"]; v != "" {
		cfg.Email.From = v
	}
	if v := env["EMAIL_FROM_NAME"]; v != "" {
		cfg.Email.FromName = v
	}
	if v := env["EMAIL_RESEND_API_KEY"]; v != "" {
		cfg.Email.Resend.APIKey = v
	}
}
