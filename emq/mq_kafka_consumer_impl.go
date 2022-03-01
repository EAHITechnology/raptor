package emq

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type KafkaConsumer struct {
	consumer      *cluster.Consumer
	client        *KafkaClient
	closeCh       chan struct{}
	returnCloseCh chan struct{}
}

type KafkaMssage struct {
	msg *sarama.ConsumerMessage
	c   *cluster.Consumer
}

func (m *KafkaMssage) Commit(metadata string) {
	m.c.MarkOffset(m.msg, metadata)
}

func (c *KafkaConsumer) InitConsumerReturn(ctx context.Context) {
	fun := "InitConsumerReturn -->"
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.returnCloseCh:
			return
		case err := <-c.consumer.Errors():
			c.client.baseConfig.L.Errorf("%s mark offset err:%v", fun, err)
		case ntf := <-c.consumer.Notifications():
			c.client.baseConfig.L.Infof("%s %s", fun, ntf.Type.String())
		}
	}
}

func (c *KafkaConsumer) ReadMsg(ctx context.Context, value interface{}) (context.Context, error) {
	select {
	case <-c.closeCh:
		return nil, nil
	case cm, ok := <-c.consumer.Messages():
		if !ok {
			return nil, ConsumerClosedErr
		}
		if err := json.Unmarshal(cm.Value, value); err != nil {
			return ctx, err
		}

		// ctx

		c.consumer.MarkOffset(cm, "")
		return ctx, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *KafkaConsumer) FetchMsg(ctx context.Context, value interface{}) (context.Context, Message, error) {
	select {
	case <-c.closeCh:
		return nil, nil, nil
	case cm, ok := <-c.consumer.Messages():
		if !ok {
			return nil, nil, ConsumerClosedErr
		}
		if err := json.Unmarshal(cm.Value, value); err != nil {
			return ctx, nil, err
		}

		message := &KafkaMssage{
			msg: cm,
			c:   c.consumer,
		}

		// ctx

		return ctx, message, nil
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
}

func (c *KafkaConsumer) ReadPayloadMsg(ctx context.Context) (context.Context, []byte, error) {
	select {
	case <-c.closeCh:
		return nil, nil, nil
	case cm, ok := <-c.consumer.Messages():
		if !ok {
			return nil, nil, ConsumerClosedErr
		}

		// ctx

		c.consumer.MarkOffset(cm, "")
		return ctx, cm.Value, nil
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
}

func (c *KafkaConsumer) FetchPayloadMsg(ctx context.Context) (context.Context, []byte, Message, error) {
	select {
	case <-c.closeCh:
		return nil, nil, nil, nil
	case cm, ok := <-c.consumer.Messages():
		if !ok {
			return nil, nil, nil, ConsumerClosedErr
		}

		message := &KafkaMssage{
			msg: cm,
			c:   c.consumer,
		}

		// ctx

		return ctx, cm.Value, message, nil
	case <-ctx.Done():
		return nil, nil, nil, ctx.Err()
	}
}

func (c *KafkaConsumer) Close() error {
	close(c.closeCh)
	close(c.returnCloseCh)
	if err := c.consumer.Close(); err != nil {
		return err
	}
	return nil
}
