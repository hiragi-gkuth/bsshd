package ids

import (
	"sync"
)

type IdsParam struct {
	*sync.RWMutex
	// threshold.Threshold
}
