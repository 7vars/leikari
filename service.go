package leikari

type ServiceExecutor interface {
	ExecuteService(Receiver, string, ...Option) (ActorHandler, error)
}