package pidfile

import (
	"fmt"
	"golang.org/x/sys/windows"
)

const (
	processQueryLimitedInformation = 0x1000

	stillActive = 259
)

func processExists(pid int) bool {
	h, err := windows.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return false
	}
	var c uint32
	err = windows.GetExitCodeProcess(h, &c)
	windows.Close(h)
	if err != nil {
		return c == stillActive
	}
	return true
}

func getDefaultPid(appName string) string {
	return fmt.Sprintf("C:\\Program Files/%s/%s.lock", appName, appName)
}
