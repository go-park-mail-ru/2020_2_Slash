package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"database"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	AvatarsDir string `json:"avatars"`
}

func (c *Config) GetDbConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)
}

func (c *Config) GetServerConnString() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetAvatarsPath() string {
	return fmt.Sprintf("./%s", c.AvatarsDir)
}

func LoadConfig(name string) (*Config, error) {
	file, err := os.Open(name)

	defer file.Close()

	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
