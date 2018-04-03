package apns

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/RobotsAndPencils/buford/push"
	"github.com/pkg/errors"

	"github.com/as/micromdm/platform/config"
	"github.com/as/micromdm/platform/pubsub"
	"github.com/as/micromdm/platform/queue"
)

type Service interface {
	Push(ctx context.Context, udid string) (string, error)
}

type Store interface {
	PushInfo(udid string) (*PushInfo, error)
}

type PushService struct {
	store    Store
	start    chan struct{}
	provider PushCertificateProvider

	mu      sync.RWMutex
	pushsvc *push.Service
}

type PushCertificateProvider interface {
	PushCertificate() (*tls.Certificate, error)
}

type Option func(*PushService)

func WithPushService(svc *push.Service) Option {
	return func(p *PushService) {
		p.pushsvc = svc
	}
}

func New(db Store, provider PushCertificateProvider, sub pubsub.Subscriber, opts ...Option) (*PushService, error) {
	pushSvc := PushService{
		store:    db,
		provider: provider,
		start:    make(chan struct{}),
	}
	for _, opt := range opts {
		opt(&pushSvc)
	}
	// if there is no push service, the push certificate hasn't been provided.
	// start a goroutine that delays the run of this service.
	if err := updateClient(&pushSvc, sub); err != nil {
		return nil, errors.Wrap(err, "wait for push service config")
	}

	if err := pushSvc.startQueuedSubscriber(sub); err != nil {
		return &pushSvc, err
	}
	return &pushSvc, nil
}

func (svc *PushService) startQueuedSubscriber(sub pubsub.Subscriber) error {
	commandQueuedEvents, err := sub.Subscribe(context.TODO(), "push-info", queue.CommandQueuedTopic)
	if err != nil {
		return errors.Wrapf(err,
			"subscribing push to %s topic", queue.CommandQueuedTopic)
	}
	go func() {
		if svc.pushsvc == nil {
			log.Println("push: waiting for push certificate before enabling APNS service provider")
			<-svc.start
			log.Println("push: service started")
		}
		for {
			select {
			case event := <-commandQueuedEvents:
				cq, err := queue.UnmarshalQueuedCommand(event.Message)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = svc.Push(context.TODO(), cq.DeviceUDID)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}()

	return nil
}

func updateClient(svc *PushService, sub pubsub.Subscriber) error {
	configEvents, err := sub.Subscribe(context.TODO(), "push-server-configs", config.ConfigTopic)
	if err != nil {
		return errors.Wrap(err, "update push service client")
	}
	go func() {
		for {
			select {
			case <-configEvents:
				pushsvc, err := NewPushService(svc.provider)
				if err != nil {
					log.Println("push: could not get push certificate %s", err)
					continue
				}
				svc.mu.Lock()
				svc.pushsvc = pushsvc
				svc.mu.Unlock()
				go func() { svc.start <- struct{}{} }() // unblock queue
			}
		}
	}()
	return nil
}

func NewPushService(provider PushCertificateProvider) (*push.Service, error) {
	cert, err := provider.PushCertificate()
	if err != nil {
		return nil, errors.Wrap(err, "get push certificate from store")
	}

	client, err := push.NewClient(*cert)
	if err != nil {
		return nil, errors.Wrap(err, "create push service client")
	}

	svc := push.NewService(client, push.Production)
	return svc, nil
}
