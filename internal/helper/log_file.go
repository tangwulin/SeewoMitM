package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"time"
)

func GetLogFile(logPath string) *os.File {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("get workdirectory failed: %v\n", err)
		panic(err)
	}

	var logDir string
	if logPath == "" {
		logDir = path.Join(wd, "logs")
	} else {
		logDir = logPath
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, fs.ModePerm)
		if err != nil {
			fmt.Printf("create log directory failed: %v\n", err)
			panic(err)
		}
	}

	now := time.Now()

	logFile, err := os.OpenFile(path.Join(logDir, now.Format("2006-01-02_150304")+".log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)

	if err != nil {
		fmt.Printf("open log file failed: %v\n", err)
		panic(err)
	}
	return logFile
}
