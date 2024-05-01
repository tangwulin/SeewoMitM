package timer

import (
	"SeewoMitM/internal/log"
	"github.com/lextoumbourou/idle"
	"time"
)

func LaunchScreenSaverTimer(timeout time.Duration, callback func()) error {
	for {
		idleTime, err := idle.Get()
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("idleTime:", idleTime)
		if idleTime >= timeout {
			log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("timeout reached!")
			callback()
			log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("nextWait:", timeout)
			time.Sleep(timeout)
		} else {
			nextWait := (timeout - idleTime) / 2
			if nextWait > 5*time.Second {
				log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("nextWait:", nextWait)
				time.Sleep(nextWait)
			} else {
				log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("nextWait: 1s")
				time.Sleep(1 * time.Second)
			}
			continue
		}
	}
}
