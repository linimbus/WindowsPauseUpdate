package main

import (
	"encoding/json"
	"time"
)

type AppConfig struct {
	FeatureEndTime   time.Time `json:"feature_end"`
	FeatureStartTime time.Time `json:"feature_start"`
	QualityEndTime   time.Time `json:"quality_end"`
	QualityStartTime time.Time `json:"quality_start"`
	PauseWeeks       float64   `json:"pause_weeks"`
}

var configCache = AppConfig{
	PauseWeeks: 4.0,
}

func ConfigGet() *AppConfig {
	return &configCache
}

func ConfigInit() error {
	return RegistryRead(&configCache)
}

func ConfigString() string {
	value, err := json.Marshal(configCache)
	if err == nil {
		return string(value)
	}
	return err.Error()
}
