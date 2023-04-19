package worker

import (
	"fmt"
	"github.com/itxiaolin/metisAi-wechat/internal/core/application"
	"github.com/itxiaolin/metisAi-wechat/internal/core/logger"
	"go.uber.org/zap"
	"sync"
)

var _ application.Application = (*Worker)(nil)

type Worker struct {
	workerEngine Engine
	startOnce    sync.Once
}

func CreateWorker(engine Engine) *Worker {
	return &Worker{
		workerEngine: engine,
	}
}

func (w *Worker) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "global recover in Run function", zap.Error(err.(error)), zap.Stack("stacktrace"))
		}
		_ = logger.GetZapLogger().Sync()
	}()
	w.start()
}

func (w *Worker) start() {
	w.startOnce.Do(func() {
		go w.workerEngine.Process()
	})
}

func (w *Worker) Stop() {
	fmt.Println("worker is stoping")
	w.workerEngine.GracefullyShutdown()
}
