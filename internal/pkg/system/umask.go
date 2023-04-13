//go:build !windows
// +build !windows

package system

import (
	"golang.org/x/sys/unix"
)

// Umask 设置当前进程的文件模式创建掩码为newmask并返回oldmask
func Umask(newmask int) (oldmask int, err error) {
	return unix.Umask(newmask), nil
	//return syscall.Umask(0), nil
}
