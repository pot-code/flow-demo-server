package event

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type kafkaEventBus struct {
	pub *kafka.Publisher
}

func NewKafkaEventBus(pub *kafka.Publisher) *kafkaEventBus {
	return &kafkaEventBus{pub: pub}
}

func (k *kafkaEventBus) Publish(e Event) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(e); err != nil {
		panic(fmt.Errorf("error encoding event: %w", err))
	}

	if err := k.pub.Publish(e.Topic(), message.NewMessage(watermill.NewUUID(), buf.Bytes())); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}
	return nil
}
