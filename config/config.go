package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type MSNServerConfiguration struct {
	DispatchServer     DispatchServer     `mapstructure:"dispatch_server"`
	NotificationServer NotificationServer `mapstructure:"notification_server"`
}

type DispatchServer struct {
	ServerAddr             string `mapstructure:"server_addr"`
	ServerPort             string `mapstructure:"server_port"`
	NotificationServerAddr string `mapstructure:"notification_server_addr"`
	NotificationServerPort string `mapstructure:"notification_server_port"`
}

type NotificationServer struct {
	ServerAddr string `mapstructure:"server_addr"`
	ServerPort string `mapstructure:"server_port"`
}

var Config MSNServerConfiguration

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/msnserver/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
