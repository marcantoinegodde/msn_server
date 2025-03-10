package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MSNServerConfiguration struct {
	Database           Database           `mapstructure:"database"`
	Redis              Redis              `mapstructure:"redis"`
	DispatchServer     DispatchServer     `mapstructure:"dispatch_server"`
	NotificationServer NotificationServer `mapstructure:"notification_server"`
	SwitchboardServer  SwitchboardServer  `mapstructure:"switchboard_server"`
	WebServer          WebServer          `mapstructure:"web_server"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Redis struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	DB            int    `mapstructure:"db"`
	PubSubChannel string `mapstructure:"pubsub_channel"`
}

type DispatchServer struct {
	ServerAddr             string `mapstructure:"server_addr"`
	ServerPort             int    `mapstructure:"server_port"`
	NotificationServerAddr string `mapstructure:"notification_server_addr"`
	NotificationServerPort int    `mapstructure:"notification_server_port"`
}

type NotificationServer struct {
	ServerAddr            string `mapstructure:"server_addr"`
	ServerPort            int    `mapstructure:"server_port"`
	SwitchboardServerAddr string `mapstructure:"switchboard_server_addr"`
	SwitchboardServerPort int    `mapstructure:"switchboard_server_port"`
}

type SwitchboardServer struct {
	ServerAddr            string `mapstructure:"server_addr"`
	ServerPort            int    `mapstructure:"server_port"`
	SwitchboardServerAddr string `mapstructure:"switchboard_server_addr"`
	SwitchboardServerPort int    `mapstructure:"switchboard_server_port"`
}

type WebServer struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

func LoadConfig() (*MSNServerConfiguration, error) {
	var config MSNServerConfiguration

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/msnserver/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("fatal error unmarshal config: %w", err)
	}

	return &config, nil
}
