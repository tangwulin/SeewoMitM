package timer

import (
	"SeewoMitM/internal/log"
	"fmt"
	"github.com/lextoumbourou/idle"
	"time"
)

func LaunchScreenSaverTimer(timeout time.Duration, callback func()) error {
	log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Info("launching ScreenSaverTimer")
	for {
		idleTime, err := idle.Get()
		if err != nil {
			log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Error(fmt.Sprintf("failed to get idle time: %v", err))
			return err
		}
		log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Trace("idleTime:", idleTime)
		if idleTime >= timeout {
			log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Trace("timeout reached!")
			callback()
			log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Trace("nextWait:", timeout)
			time.Sleep(timeout)
		} else {
			nextWait := (timeout - idleTime) / 2
			if nextWait > 5*time.Second {
				log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Trace("nextWait:", nextWait)
				time.Sleep(nextWait)
			} else {
				log.WithFields(log.Fields{"type": "ScreenSaverTimer", "timeout": timeout}).Trace("nextWait: 1s")
				time.Sleep(1 * time.Second)
			}
			continue
		}
	}
}
