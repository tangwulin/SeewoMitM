package watcher

import (
	"SeewoMitM/internal/helper"
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/screensaver"
	"SeewoMitM/internal/server"
	"github.com/fsnotify/fsnotify"
)

func LaunchConfigWatcher(configPath string) {
	configWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithFields(log.Fields{"type": "Watcher"}).Errorf("create config watcher failed:%v", err.Error())
		return
	}
	err = configWatcher.Add(configPath)
	if err != nil {
		log.WithFields(log.Fields{"type": "Watcher"}).Errorf("add config watcher failed:%v", err.Error())
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-configWatcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.WithFields(log.Fields{"type": "Watcher"}).Infof("modified file: %s", event.Name)
					configs, err := helper.ReadAndParseConfig(configPath)
					if err != nil {
						log.WithFields(log.Fields{"type": "Watcher"}).Errorf("parse config failed:%v", err.Error())
					}

					helper.SetConfig(*configs)

					server.ReloadResourceServer()
					screensaver.LoadScreenSaverConfig(configs)
				}
			case err, ok := <-configWatcher.Errors:
				if !ok {
					return
				}
				log.WithFields(log.Fields{"type": "Watcher"}).Errorf("error:%v", err.Error())
			}
		}
	}()
}
