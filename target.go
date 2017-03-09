package main

import (
	"os"
	"strings"
)

func getTargetFilesPath() []string {
	if paths := os.Getenv("NAUTILUS_SCRIPT_SELECTED_FILE_PATHS"); len(paths) != 0 {
		return strings.Split(paths[:len(paths)-1], " ")
	}
	if len(os.Args) > 1 {
		return os.Args[1:]
	}
	return nil
}

func loadTargetFile(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		b = append(b, buf[:n]...)
	}
	return b, nil
}
