package api

import (
	"configcenter/src/framework/core/log"
)

// SetLoger replace the logger
func SetLoger(target log.Loger) {
	log.SetLoger(target)
}
