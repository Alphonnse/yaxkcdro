package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	ResourceURL   string `yaml:"resource_url"`
	PathDB        string `yaml:"path_db"`
	PathStopwords string `yaml:"path_stopwords"`
}

func NewAppConfig(path string) (AppConfig, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

func (config AppConfig) GetResourceURL() string {
	return config.ResourceURL
}

func (config AppConfig) GetPathDB() string {
	return config.PathDB
}

func (config AppConfig) GetPathStopwords() string {
	return config.PathStopwords
}
