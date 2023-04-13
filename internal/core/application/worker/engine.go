package worker

type Engine interface {
	Process()
	GracefullyShutdown()
}
