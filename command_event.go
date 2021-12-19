package leikari

type DoneEvent struct{}

func Done() DoneEvent {
	return DoneEvent{}
}