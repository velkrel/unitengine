//go:build !windows

package boundedmpsc

import (
	"github.com/velkrel/unitengine/engine/internal/perftest"
	"time"
)

func processCPUTime() (time.Duration, error) { return perftest.CPUTime() }
