//go:build !windows && !darwin
// +build !windows,!darwin

package pidfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func processExists(pid int) bool {
	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(pid))); err == nil {
		return true
	}
	return false
}

func getDefaultPid(appName string) string {
	return fmt.Sprintf("/run/%s.lock", appName)
}
