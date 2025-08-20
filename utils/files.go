package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	vendor = "JeffsCasino"
)

func GetDataDir() (string, error) {
	var data string
	var err error

	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		data = filepath.Join(home, "Library/Application Support", vendor)
	case "windows":
		if v := os.Getenv("LOCALAPPDATA"); v != "" {
			data = filepath.Join(v, vendor)
		} else if v := os.Getenv("APPDATA"); v != "" {
			data = filepath.Join(v, vendor)
		} else {
			data, err = os.UserConfigDir()
			if err != nil {
				return "", err
			}
		}
	default:
		if x := os.Getenv("XDG_DATA_HOME"); x != "" {
			data = filepath.Join(x, vendor)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			data = filepath.Join(home, ".local", "share", vendor)
		}
	}

	return data, nil
}

func EnsureDirs(paths ...string) error {
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			return err
		}
	}

	return nil
}
