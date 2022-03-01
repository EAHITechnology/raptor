package emq

import (
	"context"
)

type MqConfig interface {
	SetMqType(mqTyp string)
	SetMqTopic(topic string)
	SetBrokerAddr(addrs []string)
	SetUser(user string)
	SetPassword(password string)

	GetMqType() string
	GetMqTopic() string
	GetBrokerAddr() []string
	GetUser() string
	GetPassword() string
}

type Consumer interface {
	// Read mq msg. You need to pass an object for serialization to the function.
	ReadMsg(ctx context.Context, value interface{}) (context.Context, error)

	// Fetch mq msg.You need to pass an object for serialization to the function.
	// It returns a handler for Ack msg.
	FetchMsg(ctx context.Context, value interface{}) (context.Context, Message, error)

	// read payload mq msg. Function will return original message.
	ReadPayloadMsg(ctx context.Context) (context.Context, []byte, error)

	FetchPayloadMsg(ctx context.Context) (context.Context, []byte, Message, error)

	// This method should be called when the Consumer's life cycle is ended.
	Close() error
}

type Producer interface {
	// write mq msg. You need to pass a "key" that can break up the data,a value,and an extra message to the function.
	WriteMsg(ctx context.Context, key string, value []byte, properties map[string]string) (partition int32, msgId string, err error)

	// This method should be called when the life cycle of the Producer is ended,
	// and you need to ensure that there are no write requests before calling.
	Close() error
}

type Client interface {
	// create a producer
	NewProducer(ctx context.Context) (Producer, error)

	// create a consumer
	NewConsumer(ctx context.Context, consumerName string) (Consumer, error)
}

type Message interface {
	Commit(metadata string)
}

type MqLog interface {
	Debugf(f string, args ...interface{})
	Infof(f string, args ...interface{})
	Warnf(f string, args ...interface{})
	Errorf(f string, args ...interface{})
}

func NewClient(ctx context.Context, m MqConfig) (Client, error) {
	switch m.GetMqType() {
	case KAFKA_MQ_TYP:
		kconfig, ok := m.(*KafkaConfig)
		if !ok {
			return nil, KafkaConfigIllegalErr
		}
		client, err := NewKafkaClient(ctx, kconfig)
		if err != nil {
			return nil, err
		}
		return &client, nil
	default:
		return nil, IllegalMqTypeErr
	}
}
