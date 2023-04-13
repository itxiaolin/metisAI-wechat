//go:build darwin
// +build darwin

package pidfile

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func processExists(pid int) bool {
	// OSX没有proc文件系统.
	// 使用kill -0 pid来判断进程是否存在。
	err := unix.Kill(pid, 0)
	return err == nil
}

func getDefaultPid(appName string) string {
	return fmt.Sprintf("/var/run/%s.lock", appName)
}
