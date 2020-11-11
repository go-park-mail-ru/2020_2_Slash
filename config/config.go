package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var logLevelsCode = map[string]int{
	"DEBUG": 10,
	"INFO":  20,
	"WARN":  30,
	"ERROR": 40,
	"FATAL": 50,
}

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
	PostersDir string `json:"posters"`
	VideosDir  string `json:"videos"`
	LoggerFile string `json:"logger"`
	LogLevel   string `json:"log_level"`
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

func (c *Config) GetPostersPath() string {
	return fmt.Sprintf("./%s", c.PostersDir)
}

func (c *Config) GetVideosPath() string {
	return fmt.Sprintf("./%s", c.VideosDir)
}

func (c *Config) GetLoggerDir() string {
	return c.LoggerFile
}

func (c *Config) GetLogLevel() int {
	return logLevelsCode[c.LogLevel]
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
