package publisher

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{nc: nc}
}

func (p *Publisher) Publish(subj string, t any) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return p.nc.Publish(subj, data)
}
