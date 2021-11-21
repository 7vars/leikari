package leikari

type ServiceExecutor interface {
	ExecuteService(Receiver, ...Option) (ActorHandler, error)
}