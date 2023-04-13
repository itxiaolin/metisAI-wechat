package app

import (
	"fmt"
	"github.com/itxiaolin/openai-wechat/internal/pkg/pidfile"
	"github.com/spf13/cobra"
	"os"
)

var Status = &cobra.Command{
	Use:   "status",
	Short: "status server",
	RunE:  status,
}

func status(cmd *cobra.Command, args []string) error {
	initConfig()
	if ok, pid := pidfile.CheckPIDFileAlreadyExists(pidfile.GetPidFile()); ok {
		fmt.Printf("process %d is running\n", pid)
	} else {
		fmt.Println("no running process")
		os.Exit(1)
	}
	return nil
}
