package config

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"

	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
)

// Shard is a struct to allow shard location when files are uploaded
type Shard struct {
	Depth    int
	Width    int
	RestOnly bool
}

// AllowedSize is a struct used in the allowed_sizes option
type AllowedSize struct {
	Height int
	Width  int
}

// Options is a struct to add options to the application
type Options struct {
	EnableUpload     bool          `mapstructure:"enable_upload"`
	EnableDelete     bool          `mapstructure:"enable_delete"`
	EnableStats      bool          `mapstructure:"enable_stats"`
	AllowedSizes     []AllowedSize `mapstructure:"allowed_sizes"`
	DefaultUserAgent string        `mapstructure:"default_user_agent"`
	MimetypeDetector string        `mapstructure:"mimetype_detector"`
}

// Sentry is a struct to configure sentry using a dsn
type Sentry struct {
	DSN  string
	Tags map[string]string
}

// Config is a struct to load configuration flags
type Config struct {
	Debug          bool
	Engine         *engine.Config
	Sentry         *Sentry
	SecretKey      string `mapstructure:"secret_key"`
	Shard          *Shard
	Port           int
	Options        *Options
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	Storage        *storage.Config
	KVStore        *kvstore.Config
	Logger         logger.Config
}

// DefaultConfig returns a default config instance
func DefaultConfig() *Config {
	return &Config{
		Engine: &engine.Config{
			DefaultFormat: DefaultFormat,
			Quality:       DefaultQuality,
			Format:        "",
		},
		Options: &Options{
			EnableDelete:     false,
			EnableUpload:     false,
			DefaultUserAgent: fmt.Sprint(DefaultUserAgent, "/", constants.Version),
			MimetypeDetector: DefaultMimetypeDetector,
		},
		Port: DefaultPort,
		KVStore: &kvstore.Config{
			Type: "dummy",
		},
		Shard: &Shard{
			Width:    DefaultShardWidth,
			Depth:    DefaultShardDepth,
			RestOnly: DefaultShardRestOnly,
		},
	}
}

func load(content string, isPath bool) (*Config, error) {
	config := &Config{}

	defaultConfig := DefaultConfig()

	viper.SetDefault("options", defaultConfig.Options)
	viper.SetDefault("shard", defaultConfig.Shard)
	viper.SetDefault("port", defaultConfig.Port)
	viper.SetDefault("kvstore", defaultConfig.KVStore)
	viper.SetEnvPrefix("picfit")

	var err error

	if isPath == true {
		viper.SetConfigFile(content)
		err = viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		viper.SetConfigType("json")

		err = viper.ReadConfig(bytes.NewBuffer([]byte(content)))

		if err != nil {
			return nil, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if config.Engine == nil {
		config.Engine = defaultConfig.Engine
	}

	return config, nil
}

// Load creates a Config struct from a config file path
func Load(path string) (*Config, error) {
	return load(path, true)
}

// LoadFromContent creates a Config struct from a config content
func LoadFromContent(content string) (*Config, error) {
	return load(content, false)
}
