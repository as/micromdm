package inmem

import (
	"context"

	"github.com/as/micromdm/platform/pubsub"
)

func (p *Inmem) Subscribe(_ context.Context, name, topic string) (<-chan pubsub.Event, error) {
	events := make(chan pubsub.Event)
	sub := subscription{
		name:      name,
		topic:     topic,
		eventChan: events,
	}
	p.mtx.Lock()
	p.subscriptions[topic] = append(p.subscriptions[topic], sub)
	p.mtx.Unlock()

	return events, nil
}

func (p *Inmem) dispatch() {
	for {
		select {
		case ev := <-p.publish:
			p.mtx.Lock() // TODO(as): fix
			for _, sub := range p.subscriptions[ev.Topic] {
				go func(s subscription) { s.eventChan <- ev }(sub)
			}
			p.mtx.Unlock()
		}
	}
}
