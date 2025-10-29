package extensions

import (
	"log/slog"
	"os"
	"path/filepath"
)

func IsFile(path string) (bool, error) {
	pwd, _ := os.Getwd()
	slog.Debug("checking if it's a file", "path", path, "pwd", pwd)
	info, err := os.Stat(path)
	if err != nil {
		slog.Error("error checking if it's a file", "path", path, "error", err)
		return false, err
	}

	if info.IsDir() {
		slog.Debug("file object is a directory", "path", path)
		return false, nil
	}
	slog.Debug("file object is a file", "path", path)
	return true, nil
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}
	return false, nil
}

func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return -1, err
	}
	return info.Size(), nil
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return -1, err
	}
	return size, nil
}

func ListDir(path string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			//files = append(files, filepath.Join(dir, info.Name()))
			files = append(files, dir)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return files, err
}
