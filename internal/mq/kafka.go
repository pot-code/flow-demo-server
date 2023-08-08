package mq

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

func NewKafkaSubscriber(brokers []string, consumerGroup string) *kafka.Subscriber {
	config := kafka.DefaultSaramaSubscriberConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	sub, err := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:               brokers,
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: config,
		ConsumerGroup:         consumerGroup,
	}, watermill.NewStdLogger(true, true))
	if err != nil {
		panic(fmt.Errorf("error creating kafka subscriber: %w", err))
	}
	return sub
}

func NewKafkaPublisher(brokers []string) *kafka.Publisher {
	pub, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:   brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}, watermill.NewStdLogger(true, true))
	if err != nil {
		panic(fmt.Errorf("error creating kafka publisher: %w", err))
	}
	return pub
}
