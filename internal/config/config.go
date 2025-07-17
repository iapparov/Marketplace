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