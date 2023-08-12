package event

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
)

type kafkaEventBus struct {
	pub *kafka.Publisher
}

func NewKafkaEventBus(pub *kafka.Publisher) EventBus {
	return &kafkaEventBus{pub: pub}
}

func (k *kafkaEventBus) Publish(e Event) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(e); err != nil {
		panic(fmt.Errorf("error encoding event: %w", err))
	}

	if err := k.pub.Publish(e.Topic(), message.NewMessage(watermill.NewUUID(), buf.Bytes())); err != nil {
		log.Warn().Err(err).Msg("failed to publish event")
	}
}
