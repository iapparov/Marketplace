package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlContent := `
env: test
http_port: 8081
db: ":memory:"
jwt_access_secret: "access"
jwt_refresh_secret: "refresh"
jwt_exp_access_token: 10
jwt_exp_refresh_token: 20
username:
  min_length: 3
  max_length: 20
  allowed_characters: "A-Za-z0-9_-"
  case_sensitive: true
password:
  min_length: 8
  max_length: 64
  require_upper: true
  require_lower: true
  require_digit: true
ad:
  min_length_title: 3
  max_length_title: 100
  min_length_description: 10
  max_length_description: 1000
  img_type: ["jpg", "png"]
  price_min: 0.01
`

	_, err = tmpFile.Write([]byte(yamlContent))
	if err != nil {
		t.Fatalf("failed to write config to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.Env != "test" {
		t.Errorf("expected env to be 'test', got '%s'", cfg.Env)
	}

	if cfg.Http_port != 8081 {
		t.Errorf("expected http_port to be 8081, got %d", cfg.Http_port)
	}

	if !cfg.Ad.AllowedImgTypesMap[".jpg"] {
		t.Errorf("expected AllowedImgTypesMap to include .jpg")
	}

	if !cfg.Ad.AllowedImgTypesMap[".png"] {
		t.Errorf("expected AllowedImgTypesMap to include .png")
	}

	if cfg.Ad.AllowedImgTypesMap[".gif"] {
		t.Errorf("did not expect .gif in AllowedImgTypesMap")
	}
}