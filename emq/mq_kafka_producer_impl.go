package emq

import (
	"strconv"
	"sync"

	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
)

type KafkaProducer struct {
	asyncProducer sarama.AsyncProducer
	client        *KafkaClient
	lock          *sync.Mutex
	callbackMap   map[int64]func(msg *sarama.ProducerMessage, e error)
	idx           int64
	returnCloseCh chan struct{}
}

func (p *KafkaProducer) InitProducerReturn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.returnCloseCh:
			return
		case success := <-p.asyncProducer.Successes():
			p.lock.Lock()
			idx := success.Metadata.(int64)
			p.callbackMap[idx](success, nil)
			delete(p.callbackMap, idx)
			p.lock.Unlock()
		case err := <-p.asyncProducer.Errors():
			p.lock.Lock()
			idx := err.Msg.Metadata.(int64)
			p.callbackMap[idx](err.Msg, err)
			delete(p.callbackMap, idx)
			p.lock.Unlock()
		}
	}
}

func (p *KafkaProducer) WriteMsg(ctx context.Context, key string, value []byte, properties map[string]string) (partition int32, msgId string, err error) {
	var headersIdx = 0
	headers := make([]sarama.RecordHeader, len(properties))
	for k, v := range properties {
		header := sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		}
		headers[headersIdx] = header
		headersIdx++
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	p.lock.Lock()
	msg := &sarama.ProducerMessage{
		Topic:    p.client.baseConfig.GetMqTopic(),
		Key:      sarama.StringEncoder(key),
		Value:    sarama.ByteEncoder(value),
		Headers:  headers,
		Metadata: p.idx,
	}

	var callbackErr error
	var tempMsgid int64
	var tempPartition int32
	p.callbackMap[p.idx] = func(msg *sarama.ProducerMessage, e error) {
		callbackErr = e
		tempMsgid = msg.Offset
		tempPartition = msg.Partition
		wg.Done()
	}
	p.idx++

	p.lock.Unlock()
	p.asyncProducer.Input() <- msg

	wg.Wait()
	return tempPartition, strconv.FormatInt(tempMsgid, 10), callbackErr
}

func WriteAsyncMsg(ctx context.Context, key string, value []byte, properties map[string]string) error {
	return nil
}

func (p *KafkaProducer) Close() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	close(p.returnCloseCh)
	return p.asyncProducer.Close()
}
