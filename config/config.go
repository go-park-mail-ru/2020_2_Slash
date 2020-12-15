package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var logLevelsCode = map[string]int{
	"DEBUG": 10,
	"INFO":  20,
	"WARN":  30,
	"ERROR": 40,
	"FATAL": 50,
}

type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Config struct {
	Database              Database `json:"database"`
	TestDatabase          Database `json:"test_database"`
	Server                Server   `json:"server"`
	UserblockMicroservice Server   `json:"userblock_microservice"`
	AuthMicroservice      Server   `json:"auth_microservice"`
	AvatarsDir            string   `json:"avatars"`
	PostersDir            string   `json:"posters"`
	VideosDir             string   `json:"videos"`
	LoggerFile            string   `json:"logger"`
	LogLevel              string   `json:"log_level"`
}

func getDbConnString(database Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		database.Host, database.Port, database.User, database.Password, database.Name)
}

func (c *Config) GetProdDbConnString() string {
	return getDbConnString(c.Database)
}

func (c *Config) GetTestDbConnString() string {
	return getDbConnString(c.TestDatabase)
}

func (c *Config) GetServerConnString() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetUserblockMSConnString() string {
	return fmt.Sprintf("%s:%d", c.UserblockMicroservice.Host,
		c.UserblockMicroservice.Port)
}

func (c *Config) GetAuthMSConnString() string {
	return fmt.Sprintf("%s:%d", c.AuthMicroservice.Host,
		c.AuthMicroservice.Port)
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
	file, err := os.Open(filepath.Clean(name))

	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return config, nil
}
