package helper

import "os"

func GetFileInfo(path string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return fileInfo, nil
}
