package app

import (
	"errors"
	"fmt"
	"github.com/itxiaolin/openai-wechat/internal/cfg"
	"github.com/itxiaolin/openai-wechat/internal/core/logger"
	"github.com/itxiaolin/openai-wechat/internal/global"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func usageAndError(cmd *cobra.Command) error {
	if err := cmd.Usage(); err != nil {
		return err
	}
	return errors.New("invalid arguments")
}

func signalLoop() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		sig := <-sigCh
		logger.Info(nil, fmt.Sprintf("got signal %v", sig))
		switch sig {
		case syscall.SIGHUP:
			initConfig()
			cfg.LoadConfig(&global.Config, _cfgFile)
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			if robot != nil {
				robot.Stop()
			}
			_ = killProcess()
			return
		}
	}

}
