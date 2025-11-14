// Package config provides configuration loading utilities
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from file and environment variables
// It searches for config file in: current dir, home dir, /etc/appname/
func LoadConfig(appName string, configStruct interface{}) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Add search paths
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	
	// Enable environment variables
	viper.AutomaticEnv()
	
	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, continue with env vars only
			return nil
		}
		return fmt.Errorf("error reading config file: %w", err)
	}
	
	// Unmarshal config into struct
	if err := viper.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	return nil
}

// LoadConfigFromFile loads configuration from a specific file
func LoadConfigFromFile(configPath string, configStruct interface{}) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}
	
	viper.SetConfigFile(configPath)
	
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}
	
	if err := viper.Unmarshal(configStruct); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	return nil
}

// GetString gets a string value from config
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt gets an integer value from config
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool gets a boolean value from config
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	viper.Set(key, value)
}
