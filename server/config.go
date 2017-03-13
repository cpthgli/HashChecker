package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	AutoCheck AutoCheckConfig
}

// AutoCheckConfig ...
type AutoCheckConfig struct {
	Enable bool
	Md5    bool
	Sha1   bool
	Sha256 bool
}

func getConfigPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := dir + "/config.toml"
	return path
}
func isExistConfigFile() bool {
	_, err := os.Stat(getConfigPath())
	return !(err != nil)
}

// Save ...
func (cp *Config) Save() error {
	var buf bytes.Buffer
	var encoder = toml.NewEncoder(&buf)
	if err := encoder.Encode(*cp); err != nil {
		return err
	}
	path := getConfigPath()
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	writer := bufio.NewWriter(fp)
	if _, err = writer.WriteString(buf.String()); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// Load ...
func (cp *Config) Load() error {
	if !isExistConfigFile() {
		if err := cp.Save(); err != nil {
			return err
		}
	}
	path := getConfigPath()
	if _, err := toml.DecodeFile(path, cp); err != nil {
		return err
	}
	return nil
}
