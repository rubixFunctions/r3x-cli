package cmd

import (
	"fmt"
	"io"
	"os"
)

func exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		fmt.Println(err)
	}
	return false
}

func isEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}
