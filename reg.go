package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"

	"golang.org/x/sys/windows/registry"
)

const TIME_LAYOUT = "2025-03-19T02:08:26Z"

func RegistryWriteTime(key registry.Key, name string, input time.Time) error {
	value := input.Format("2006-01-02T15:04:05Z")

	logs.Info("write registry name:%s value:%s", name, value)

	if err := key.SetStringValue(name, value); err != nil {
		return fmt.Errorf("write registry name: %s failed, %s", name, err.Error())
	}
	return nil
}

func RegistryReadTime(key registry.Key, name string, output *time.Time) error {
	value, _, err := key.GetStringValue(name)
	if err != nil {
		return fmt.Errorf("read registry name:%s failed, %s", name, err.Error())
	}
	logs.Info("read registry name:%s Value:%s", name, value)

	tm, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fmt.Errorf("parse timestamp failed, %s", err.Error())
	}
	*output = tm
	return nil
}

func RegistryRead(config *AppConfig) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.READ)
	if err != nil {
		return fmt.Errorf("open registry failed, %s", err.Error())
	}
	defer key.Close()

	err = RegistryReadTime(key, "PauseFeatureUpdatesStartTime", &config.FeatureStartTime)
	if err != nil {
		return err
	}

	err = RegistryReadTime(key, "PauseFeatureUpdatesEndTime", &config.FeatureEndTime)
	if err != nil {
		return err
	}

	err = RegistryReadTime(key, "PauseQualityUpdatesStartTime", &config.QualityStartTime)
	if err != nil {
		return err
	}

	err = RegistryReadTime(key, "PauseQualityUpdatesEndTime", &config.QualityEndTime)
	if err != nil {
		return err
	}

	err = RegistryReadTime(key, "PauseUpdatesStartTime", &config.QualityStartTime)
	if err != nil {
		return err
	}

	err = RegistryReadTime(key, "PauseUpdatesExpiryTime", &config.QualityEndTime)
	if err != nil {
		return err
	}

	logs.Info("registry config read all done")

	return nil
}

func RegistryWrite(config AppConfig) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\WindowsUpdate\UX\Settings`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("open registry failed, %s", err.Error())
	}
	defer key.Close()

	err = RegistryWriteTime(key, "PauseFeatureUpdatesStartTime", config.FeatureStartTime)
	if err != nil {
		return err
	}

	err = RegistryWriteTime(key, "PauseFeatureUpdatesEndTime", config.FeatureEndTime)
	if err != nil {
		return err
	}

	err = RegistryWriteTime(key, "PauseQualityUpdatesStartTime", config.QualityStartTime)
	if err != nil {
		return err
	}

	err = RegistryWriteTime(key, "PauseQualityUpdatesEndTime", config.QualityEndTime)
	if err != nil {
		return err
	}

	err = RegistryWriteTime(key, "PauseUpdatesStartTime", config.QualityStartTime)
	if err != nil {
		return err
	}

	err = RegistryWriteTime(key, "PauseUpdatesExpiryTime", config.QualityEndTime)
	if err != nil {
		return err
	}

	logs.Info("registry config write all done")

	return nil
}
