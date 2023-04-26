//go:build windows

package clock

import (
	"errors"
	"time"
)

func Set(t time.Time) error {
	// TODO
	return errors.New("clock Set() not implemented on Windows")
}
