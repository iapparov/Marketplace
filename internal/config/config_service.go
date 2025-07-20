package config

import (
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"strings"
)


func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.Ad.AllowedImgTypesMap = make(map[string]bool)
	for _, ext := range cfg.Ad.ImgType {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		cfg.Ad.AllowedImgTypesMap[strings.ToLower(ext)] = true
	}

	return &cfg, nil
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		log.Fatal("CONFIG_PATH env is required")
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}