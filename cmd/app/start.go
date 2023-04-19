package app

import (
	"errors"
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/application/wxrobot"
	"github.com/itxiaolin/metisAi-wechat/internal/core/application/worker"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"github.com/itxiaolin/metisAi-wechat/internal/pkg/daemon"
	"github.com/itxiaolin/metisAi-wechat/internal/pkg/pidfile"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var Start = &cobra.Command{
	Use:   "start",
	Short: "start server",
	RunE:  start,
}
var _mkDaemon bool
var robot *worker.Worker

func start(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return usageAndError(cmd)
	}
	if _mkDaemon {
		if err := daemon.MakeDaemon(); err != nil {
			logger.Error(nil, "make daemon 异常")
			return err
		}
	}
	initConfig()
	return run()
}

func run() error {
	pidFile := pidfile.GetPidFile()
	if ok, pid := pidfile.CheckPIDFileAlreadyExists(pidFile); ok {
		return errors.New(fmt.Sprintf("%s is running,pid : %d", global.Config.System.AppName, pid))
	} else {
		_ = os.Remove(pidFile)
	}
	_, err := pidfile.New(pidFile)
	if err != nil {
		fmt.Println("create pidFile error")
		logger.Error(nil, "create pidFile error", zap.Error(err))
		return errors.New("create pidFile error")
	}
	logger.Info(nil, "start success", zap.String("app name", global.Config.System.AppName))
	robot = worker.CreateWorker(wxrobot.NewWXBotEngine())
	go robot.Run()
	signalLoop()
	return nil
}
