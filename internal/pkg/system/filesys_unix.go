//go:build !windows
// +build !windows

package system

import "os"

func MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// MkdirAllWithACL 是os的包装器。MkdirAll在unix系统。
func MkdirAllWithACL(path string, perm os.FileMode, sddl string) error {
	return os.MkdirAll(path, perm)
}

func GetServiceFilePath() string {
	return "/etc/systemd/system/go-chatGPT.service"
}
