package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Pagination PaginationConfig
	Peers      PeersConfig
}

type PeersConfig struct {
	MdmsHost             string `mapstructure:"mdmsHost"`
	MdmsSearchPath       string `mapstructure:"mdmsSearchPath"`
	PropertyHost         string `mapstructure:"propertyHost"`
	PropertySearchPath   string `mapstructure:"propertySearchPath"`
	IdgenHost            string `mapstructure:"idgenHost"`
	IdgenGeneratePath    string `mapstructure:"idgenGeneratePath"`
	UserHost             string `mapstructure:"userHost"`
	UserSearchPath       string `mapstructure:"userSearchPath"`
	UserCreateNoValidate string `mapstructure:"userCreateNoValidatePath"`
}

type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	ContextPath string `mapstructure:"contextPath"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

type PaginationConfig struct {
	DefaultLimit  int `mapstructure:"defaultLimit"`
	MaxLimit      int `mapstructure:"maxLimit"`
	DefaultOffset int `mapstructure:"defaultOffset"`
}

// Load reads configs/application.yaml, then lets environment variables
// (SERVER_PORT, DATABASE_HOST, ...) override any value — no hardcoded
// credentials in code, per the audit's config-management finding.
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", 8091)
	v.SetDefault("server.contextPath", "/sw-services")
	v.SetDefault("pagination.defaultLimit", 50)
	v.SetDefault("pagination.maxLimit", 500)
	v.SetDefault("pagination.defaultOffset", 0)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: failed to read %s: %w", configPath, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal: %w", err)
	}
	return &cfg, nil
}
