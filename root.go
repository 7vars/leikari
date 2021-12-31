package leikari

import "sync"

type rootActor struct {
	sync.RWMutex
	subscriptions []Subscribe
}

func root() Receiver {
	return &rootActor{
		subscriptions: make([]Subscribe, 0),
	}
}

func (r *rootActor) addSubsciption(s Subscribe) {
	r.Lock()
	defer r.Unlock()

	r.subscriptions = append(r.subscriptions, s)
}

func (r *rootActor) removeSubscription(ref Ref) {
	r.Lock()
	defer r.Unlock()

	for i, s := range r.subscriptions {
		if s.Ref == ref {
			r.subscriptions = append(r.subscriptions[:i], r.subscriptions[i+1:]...)
		}
	}
}

func (r *rootActor) Receive(ctx ActorContext, msg Message) {
	if msg.Value() == nil {
		return
	}

	switch val := msg.Value().(type) {
	case Subscribe:
		r.addSubsciption(val)
		msg.Reply(Done())
	case Unsubscribe:
		r.removeSubscription(val.Ref)
		msg.Reply(Done)
	case Publish:
		r.RLock()
		defer r.RUnlock()

		for _, s := range r.subscriptions {
			if s.Filter(val.Content) {
				s.Ref.Send(val.Content)
			}
		}
		msg.Reply(Done())
	}
}