package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Targets ...
type Targets struct {
	paths []string
	data  [][]byte
}

func getNautilusTargetsPath() []string {
	if paths := os.Getenv("NAUTILUS_SCRIPT_SELECTED_FILE_PATHS"); len(paths) > 0 {
		return strings.Split(paths[:len(paths)-1], " ")
	}
	return nil
}

func getFlagTargetsPath() []string {
	paths := strings.Split(flag.Lookup("path").Value.String(), " ")
	for i, path := range paths {
		var err error
		paths[i], err = filepath.Abs(path)
		if err != nil {
			return nil
		}
	}
	return paths
}

func getTargetPath() []string {
	if flag.Lookup("path").Value.String() == flag.Lookup("path").DefValue {
		if paths := getNautilusTargetsPath(); len(paths) != 0 {
			return paths
		}
	} else {
		if paths := getFlagTargetsPath(); len(paths) != 0 {
			return paths
		}
	}
	return nil
}

// Load ...
func (t *Targets) Load() error {
	paths := getTargetPath()
	if paths == nil {
		return fmt.Errorf("Cannot load target file %v", paths)
	}
	t.paths = paths
	for _, path := range paths {
		t.data = append(t.data, FileLoad(path))
	}
	return nil
}

// FileLoad ...
func FileLoad(path string) []byte {
	// REMINDER: return err
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return nil
	}
	fp, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer fp.Close()
	var b []byte
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return nil
		}
		b = append(b, buf[:n]...)
	}
	return b
}
