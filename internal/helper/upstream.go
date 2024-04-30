package helper

import (
	"SeewoMitM/internal/log"
	"errors"
	"time"
)

func GetUpstream(retry int) (port int, err error) {
	for i := 0; i < retry; i++ {
		_, err := FindPidByName("SeewoCore.exe")
		if err != nil {
			log.WithFields(log.Fields{"type": "CheckSeewoCoreIsNoRunning"}).Error(err.Error())
			path, err := GetSeewoCorePath()
			if err != nil {
				log.WithFields(log.Fields{"type": "GetSeewoCorePath"}).Error(err.Error())
				return -1, errors.New("could not find SeewoCore.exe, please make sure SeewoCore is installed")
			}
			log.WithFields(log.Fields{"type": "GetSeewoCorePath"}).Info("GetSeewoCorePath success, path: " + path)

			log.WithFields(log.Fields{"type": "StartProcess"}).Info("SeewoCore.exe is not running, start it")
			err = StartProcess(path)
			if err != nil {
				log.WithFields(log.Fields{"type": "StartProcess"}).Error(err.Error())
				return -1, errors.New("could not start SeewoCore.exe, please check the path: " + path)
			}
			log.WithFields(log.Fields{"type": "StartProcess"}).Info("waiting for SeewoCore.exe to start")
			time.Sleep(5 * time.Second) //等待SeewoCore.exe启动
			continue
		}

		//如果SeewoCore.exe进程存在，获取端口
		port, err := GetUpstreamPort()
		if err != nil {
			log.WithFields(log.Fields{"type": "GetUpstreamPort"}).Error(err.Error())
			return -1, err
		}
		return port, nil
	}
	log.WithFields(log.Fields{"type": "GetUpstreamPort"}).Error("could not get upstream port")
	return -1, errors.New("could not get upstream port")
}
