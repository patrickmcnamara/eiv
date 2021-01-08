package main

import (
	"os"
	"path/filepath"
	"plugin"
)

func loadPlugins() error {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	pluginPath := filepath.Join(configPath, "eiv", "plugin")
	return filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".eivp" {
			_, err := plugin.Open(path)
			return err
		}
		return nil
	})
}
