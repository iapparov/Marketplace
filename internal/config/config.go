package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
    Env	string	`yaml:"env" env-default:"local"`
    Http_port	int	`yaml:"http_port" env-default:"8080"`
	Db	string	`yaml:"db" env-default:"./storage/marketplace.db"`
	JWT_ACCESS_SECRET	string	`yaml:"JWT_ACCESS_SECRET" env-default:"YOUR_JWT_SECRET"`
	JWT_REFRESH_SECRET	string	`yaml:"JWT_REFRESH_SECRET" env-default:"YOUR_JWT_SECRET"`
	JWT_EXP_ACCESS_TOKEN	int	`yaml:"JWT_EXP_ACCESS_TOKEN" env-default:"15"`
	JWT_EXP_REFRESH_TOKEN	int	`yaml:"JWT_EXP_REFRESH_TOKEN" env-default:"24"`
}


func MustLoad() *Config{
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		log.Fatalf("CONFIG_PATH env is required")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("%s", "failed to read config file:" + err.Error())
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil{
		log.Fatalf("%s", "failed to unmarshal config: " + err.Error())
	}


	return &cfg
}