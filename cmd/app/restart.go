package app

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var Restart = &cobra.Command{
	Use:   "restart",
	Short: "restart server",
	RunE:  restart,
}

func restart(c *cobra.Command, args []string) error {
	err := stop(c, args)
	if err != nil {
		fmt.Println("stop error")
		return errors.New("restart error")
	}
	fmt.Println("stop success.")
	time.Sleep(2 * time.Second)
	err = start(c, args)
	if err != nil {
		fmt.Println("start error.")
		return errors.New("restart error")
	}
	return nil
}
