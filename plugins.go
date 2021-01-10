package main

import (
	"os"
	"path/filepath"
	"plugin"
	"runtime"
)

func loadPlugins() error {
	switch runtime.GOOS {
	case "darwin", "freebsd", "linux":
	default:
		return nil
	}
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
