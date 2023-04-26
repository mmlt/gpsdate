//go:build linux

package clock

import (
	"errors"
	"syscall"
	"time"
)

func Set(t time.Time) error {
	tv := syscall.NsecToTimeval(t.UnixNano())
	if err := syscall.Settimeofday(&tv); err != nil {
		return errors.New("settimeofday: " + err.Error())
	}
	return nil
}
