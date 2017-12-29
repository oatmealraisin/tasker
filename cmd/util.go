package cmd

import (
	"os"
	"path/filepath"
)

// TODO: Configure w/ Viper
func setup() error {
	XDG_DATA_HOME := os.Getenv("XDG_DATA_HOME")
	HOME := os.Getenv("HOME")

	if workDir == "" {
		if len(XDG_DATA_HOME) != 0 {
			workDir = filepath.Join(XDG_DATA_HOME, "tasker")
		} else if len(HOME) != 0 {
			workDir = filepath.Join(os.Getenv("HOME"), ".local", "tasker")
		} else {
			workDir = filepath.Join("/etc", "tasker")
		}

		if err := os.MkdirAll(workDir, os.FileMode(int(0777))); err != nil {
			return err
		}
	}

	return nil
}
