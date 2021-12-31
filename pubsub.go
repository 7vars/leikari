package leikari

type PubSub interface {
	Subscribe(Ref, func(interface{}) bool)
	Unsubscribe(Ref)
	Publish(interface{})
}

type Subscribe struct {
	Ref Ref
	Filter func(interface{}) bool
}

type Unsubscribe struct {
	Ref Ref
}

type Publish struct {
	Content interface{}
}

