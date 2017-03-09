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

func save(cp *Config) error {
	var buf bytes.Buffer
	var encoder = toml.NewEncoder(&buf)
	if err := encoder.Encode(*cp); err != nil {
		return err
	}
	fp, err := os.Create(getConfigPath())
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

func load() (*Config, error) {
	cp := new(Config)
	if !isExistConfigFile() {
		save(cp)
	}
	if _, err := toml.DecodeFile(getConfigPath(), cp); err != nil {
		return nil, err
	}
	return cp, nil
}
