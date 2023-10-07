package mq

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/rs/zerolog"
)

func NewKafkaSubscriber(brokers []string, consumerGroup string, l zerolog.Logger) *kafka.Subscriber {
	config := kafka.DefaultSaramaSubscriberConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	sub, err := kafka.NewSubscriber(kafka.SubscriberConfig{
		Brokers:               brokers,
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: config,
		ConsumerGroup:         consumerGroup,
	}, newWatermillLogger(l))
	if err != nil {
		panic(fmt.Errorf("error creating kafka subscriber: %w", err))
	}
	return sub
}

func NewKafkaPublisher(brokers []string, l zerolog.Logger) *kafka.Publisher {
	pub, err := kafka.NewPublisher(kafka.PublisherConfig{
		Brokers:   brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}, newWatermillLogger(l))
	if err != nil {
		panic(fmt.Errorf("error creating kafka publisher: %w", err))
	}
	return pub
}

type watermillLogger struct {
	l zerolog.Logger
}

func newWatermillLogger(l zerolog.Logger) watermill.LoggerAdapter {
	return &watermillLogger{l: l}
}

func (w *watermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	w.l.Err(err).Fields(map[string]interface{}(fields)).Msg(msg)
}

func (w *watermillLogger) Info(msg string, fields watermill.LogFields) {
	w.l.Info().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (w *watermillLogger) Debug(msg string, fields watermill.LogFields) {
	w.l.Debug().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (w *watermillLogger) Trace(msg string, fields watermill.LogFields) {
	w.l.Trace().Fields(map[string]interface{}(fields)).Msg(msg)
}

func (w *watermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &watermillLogger{l: w.l.With().Fields(map[string]interface{}(fields)).Logger()}
}
