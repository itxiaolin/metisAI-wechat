package daemon

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/pkg/system"
	"os"
	"strconv"
	"strings"
	"time"
)

const stageVar = "__DAEMON_STAGE"

func getStage() (stage int, advanceStage func() error, resetEnv func() error) {
	var origValue string
	stage = 0

	daemonStage := os.Getenv(stageVar)
	stageTag := strings.SplitN(daemonStage, ":", 2)
	stageInfo := strings.SplitN(stageTag[0], "/", 3)

	if len(stageInfo) == 3 {
		stageStr, tm, check := stageInfo[0], stageInfo[1], stageInfo[2]

		hash := sha1.New()
		hash.Write([]byte(stageStr + "/" + tm + "/"))

		if check != hex.EncodeToString(hash.Sum([]byte{})) {
			origValue = daemonStage
		} else {
			stage, _ = strconv.Atoi(stageStr)

			if len(stageTag) == 2 {
				origValue = stageTag[1]
			}
		}
	} else {
		origValue = daemonStage
	}

	advanceStage = func() error {
		base := fmt.Sprintf("%d/%09d/", stage+1, time.Now().Nanosecond())
		hash := sha1.New()
		hash.Write([]byte(base))
		tag := base + hex.EncodeToString(hash.Sum([]byte{}))

		if err := os.Setenv(stageVar, tag+":"+origValue); err != nil {
			return fmt.Errorf("can't set %s: %s", stageVar, err)
		}
		return nil
	}
	resetEnv = func() error {
		return os.Setenv(stageVar, origValue)
	}

	return stage, advanceStage, resetEnv
}

func MakeDaemon() error {
	//global.Logger.Info("start make daemon .")
	stage, advanceStage, resetEnv := getStage()
	fatal := func(err error) error {
		if stage > 0 {
			os.Exit(1)
		}
		resetEnv()
		return err
	}

	if stage < 2 {
		procName, err := os.Executable()
		//fmt.Println(procName, zap.Int("stage", stage))
		if err != nil {
			return fatal(fmt.Errorf("can't determine full path to executable: %s", err))
		}
		if procName == "" {
			return fatal(fmt.Errorf("can't determine full path to executable"))
		}

		if err = advanceStage(); err != nil {
			return fatal(err)
		}
		dir, err := os.Getwd()
		if err != nil {
			return fatal(err)
		}
		files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
		osAttrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files}

		if stage == 0 {
			osAttrs.Sys = getSysProcAttr()
		}

		proc, err := os.StartProcess(procName, os.Args, &osAttrs)
		if err != nil {
			return fatal(fmt.Errorf("can't create process %s: %s", procName, err))
		}
		proc.Release()
		os.Exit(0)
	}
	system.Umask(0)
	resetEnv()

	return nil
}
