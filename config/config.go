package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
        Name     string `mapstructure:"name"`
        Env      string `mapstructure:"env"`
        Port     int    `mapstructure:"port"`
        BaseURL  string `mapstructure:"base_url"`
        Timezone string `mapstructure:"timezone"`
    } `mapstructure:"app"`

    Logging struct {
        Level     string `mapstructure:"level"`
        Format    string `mapstructure:"format"`
        Output    string `mapstructure:"output"`
        FilePath  string `mapstructure:"file_path"`
    } `mapstructure:"logging"`

    Server struct {
        ReadTimeout     time.Duration `mapstructure:"read_timeout"`
        WriteTimeout    time.Duration `mapstructure:"write_timeout"`
        IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
        MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
    } `mapstructure:"server"`

	ExternalAPI struct {
        NewsAPI struct {
            BaseURL string `mapstructure:"base_url"`
            APIKey  string `mapstructure:"api_key"`
        } `mapstructure:"newsapi"`
	}

	CORS struct {
		AllowedOrigins   []string `mapstructure:"allowed_origins"`
        AllowedMethods   []string `mapstructure:"allowed_methods"`
        AllowedHeaders   []string `mapstructure:"allowed_headers"`
        AllowCredentials bool     `mapstructure:"allow_credentials"`
        MaxAge           int      `mapstructure:"max_age"`
    } `mapstructure:"cors"`
}

var AppConfig *Config

func LoadConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error loading config file: %v", err)
    }

    AppConfig = &Config{}
    if err := viper.Unmarshal(AppConfig); err != nil {
        log.Fatalf("Unable to decode into struct: %v", err)
    }
}