package app

import (
	"bytes"
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/cfg"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"github.com/itxiaolin/metisAi-wechat/internal/global"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var _cfgFile, _exeName string

var RootCmd = &cobra.Command{
	Use: "metisAi-wechat",
	RunE: func(c *cobra.Command, args []string) error {
		return c.Usage()
	},
}

func init() {
	cobra.EnableCommandSorting = false
	_exeName = filepath.Base(os.Args[0])
	RootCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		if err := c.Usage(); err != nil {
			return err
		}
		fmt.Fprintln(c.OutOrStderr())
		return err
	})
	f := RootCmd.PersistentFlags()
	systemExecPath, _ := os.Executable()
	split := strings.Split(systemExecPath, "/")
	split[len(split)-1] = "config.yaml"
	var b bytes.Buffer
	for _, str := range split {
		if str == "" {
			continue
		}
		b.WriteString("/")
		b.WriteString(str)
	}
	f.StringVar(&_cfgFile, "config", "config/config.yaml", "config file")
	if _cfgFile == "config.yaml" {
		_cfgFile = b.String()
	}
	Start.PersistentFlags().BoolVarP(&_mkDaemon, "daemon", "d", false, "run as daemon")
	RootCmd.AddCommand(Start, Stop, Restart, Status)
}

func Main() {
	os.Setenv("LANG", "en_US.UTF-8")
	args := os.Args[:]
	if len(args) == 1 {
		args = append(args, "start")
	}
	RootCmd.SetArgs(args[1:])
	if err := RootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}

func initConfig() {
	cfg.LoadConfig(&global.Config, _cfgFile)
	logger.BuildZapLogger(global.Config.Logger)
}
