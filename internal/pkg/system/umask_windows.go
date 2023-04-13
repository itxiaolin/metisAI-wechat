// @Author  linxianqin  2022/11/25 10:11
package system

// Umask 不支持windows平台
func Umask(newmask int) (oldmask int, err error) {
	// should not be called on cli code path
	return 0, ErrNotSupportedPlatform
}
