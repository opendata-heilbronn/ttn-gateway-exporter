package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type TargetConfig struct {
	DefaultBaseUrl string   `yaml:"default_base_url" json:"default_base_url"`
	Targets        []Target `yaml:"targets" json:"targets"`
}

type Target struct {
	GatewayID string `yaml:"gateway_id" json:"gateway_id"`
	APIKey    string `yaml:"api_key" json:"api_key"`
	BaseUrl   string `yaml:"base_url" json:"base_url"`
}

func ReadTargets(location string) (TargetConfig, error) {
	file, err := os.Open(location)
	if err != nil {
		return TargetConfig{}, err
	}

	targetConfig := TargetConfig{
		DefaultBaseUrl: "https://eu1.cloud.thethings.network",
	}
	err = yaml.NewDecoder(file).Decode(&targetConfig)
	if err != nil {
		return TargetConfig{}, nil
	}

	for i := range targetConfig.Targets {
		if targetConfig.Targets[i].BaseUrl == "" {
			targetConfig.Targets[i].BaseUrl = targetConfig.DefaultBaseUrl
		}
	}

	return targetConfig, err
}
