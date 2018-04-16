package regular

import (
	"configcenter/src/framework/core/timer/types"
	"time"
)

// Regular the conventianl timer
type Regular interface {

	// AddTimerFunc register the timing callback function
	AddTimerFunc(timeDuration time.Duration, timerFunc types.TimerFunc) types.TimerKey

	// RemoveTimerFunc remove the timing callback function by the timer key
	RemoveTimerFunc(timerKey types.TimerKey)
}
