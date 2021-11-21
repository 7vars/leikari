package leikari

type Message interface {
	Value() interface{}
	Reply(interface{})
}

type sendOnly struct {
	value interface{}
}

func Send(v interface{}) Message {
	return sendOnly{
		value: v,
	}
}

func (so sendOnly) Value() interface{} {
	return so.value
}

func (so sendOnly) Reply(interface{}) {}

type request struct {
	reply chan<- interface{}
	value interface{}
}

func Request(reply chan<- interface{}, v interface{}) Message {
	return &request{
		reply: reply,
		value: v,
	}
}

func (r request) Value() interface{} {
	return r.value
}

func (r request) Reply(v interface{}) {
	r.reply <- v
}

type Pusher interface {
	Push(Message) error
}