package app

import (
	"errors"
	"fmt"
	"github.com/itxiaolin/openai-wechat/internal/pkg/pidfile"
	"github.com/itxiaolin/openai-wechat/internal/pkg/system"
	"github.com/spf13/cobra"
)

var Stop = &cobra.Command{
	Use:   "stop",
	Short: "stop server",
	RunE:  stop,
}

func stop(cmd *cobra.Command, args []string) error {
	initConfig()
	return killProcess()
}

func killProcess() error {
	pidFile := pidfile.GetPidFile()
	fmt.Printf("killProcess pidFile: %s\n", pidFile)
	if ok, pid := pidfile.CheckPIDFileAlreadyExists(pidFile); ok {
		fmt.Println(pid)
		err := system.KillProcess(pid)
		if err != nil {
			return fmt.Errorf("term PID %d error: %v", pid, err)
		}
		fmt.Printf("term signal has send to PID %d\n", pid)
		//_ = os.Remove(pidFile)
		return nil
	}
	return errors.New("no running process")
}
