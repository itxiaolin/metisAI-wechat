//go:build linux || freebsd || darwin
// +build linux freebsd darwin

package system // import "github.com/docker/docker/pkg/system"

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// IsProcessAlive 如果带有给定pid的进程正在运行，则返回true.
func IsProcessAlive(pid int) bool {
	err := unix.Kill(pid, syscall.Signal(0))
	if err == nil || err == unix.EPERM {
		return true
	}

	return false
}

// KillProcess 强制停止进程.
func KillProcess(pid int) error {
	return unix.Kill(pid, syscall.SIGTERM)
}

// IsProcessZombie 如果进程的状态为"Z"则返回true
// http://man7.org/linux/man-pages/man5/proc.5.html
func IsProcessZombie(pid int) (bool, error) {
	statPath := fmt.Sprintf("/proc/%d/stat", pid)
	dataBytes, err := os.ReadFile(statPath)
	if err != nil {
		return false, err
	}
	data := string(dataBytes)
	sdata := strings.SplitN(data, " ", 4)
	if len(sdata) >= 3 && sdata[2] == "Z" {
		return true, nil
	}

	return false, nil
}
