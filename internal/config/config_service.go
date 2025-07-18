package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"strings"
)


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

	cfg.Ad.AllowedImgTypesMap = make(map[string]bool)
	for _, ext := range cfg.Ad.ImgType {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		cfg.Ad.AllowedImgTypesMap[strings.ToLower(ext)] = true
	}


	return &cfg
}