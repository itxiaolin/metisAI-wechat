package pidfile

import (
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"github.com/itxiaolin/metisAi-wechat/internal/pkg/system"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDFile 存储正在运行进程的进程ID的文件.
type PIDFile struct {
	Path string
}

func GetPidFile() string {
	filePath := global.Config.System.PidFile
	if filePath == "" {
		filePath = getDefaultPid(global.Config.System.AppName)
	}
	return filePath
}

// CheckPIDFileAlreadyExists 检查pid文件是否已经存在
func CheckPIDFileAlreadyExists(path string) (bool, int) {
	if pidByte, err := os.ReadFile(path); err == nil {
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil {
			if processExists(pid) {
				return true, pid
			}
		}
	}
	return false, -1
}

// New 使用指定的路径创建PID文件.
func New(path string) (*PIDFile, error) {
	if ok, _ := CheckPIDFileAlreadyExists(path); ok {
		return nil, fmt.Errorf("pid file found, ensure app is not running or delete %s", path)
	}
	// 注意:如果目录已经存在，MkdirAll返回nil
	if err := system.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		logger.Error(nil, "目录创建失败", zap.Error(err))
		return nil, err
	}
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		logger.Error(nil, "将pid写入文件失败", zap.Error(err))
		return nil, err
	}
	return &PIDFile{Path: path}, nil
}
