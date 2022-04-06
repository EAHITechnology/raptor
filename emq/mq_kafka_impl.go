package emq

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"golang.org/x/net/context"
)

const (
	KAFKA_MQ_TYP = "kafka"
)

type KafkaConfig struct {
	TopicType  string
	Topic      string
	BrokerAddr []string
	User       string
	Password   string
	L          MqLog
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{}
}

func (conf *KafkaConfig) SetMqType(mqTyp string) {
	conf.TopicType = mqTyp
}

func (conf *KafkaConfig) GetMqType() string {
	return conf.TopicType
}

func (conf *KafkaConfig) SetMqTopic(topic string) {
	conf.Topic = topic
}

func (conf *KafkaConfig) GetMqTopic() string {
	return conf.Topic
}

func (conf *KafkaConfig) SetBrokerAddr(addrs []string) {
	conf.BrokerAddr = addrs
}

func (conf *KafkaConfig) GetBrokerAddr() []string {
	return conf.BrokerAddr
}

func (conf *KafkaConfig) SetUser(user string) {
	conf.User = user
}

func (conf *KafkaConfig) GetUser() string {
	return conf.User
}

func (conf *KafkaConfig) SetPassword(password string) {
	conf.Password = password
}

func (conf *KafkaConfig) GetPassword() string {
	return conf.Password
}

func (conf *KafkaConfig) SetLog(logger MqLog) {
	conf.L = logger
}

func (conf *KafkaConfig) GetLog() MqLog {
	return conf.L
}

type KafkaClient struct {
	baseConfig *KafkaConfig
	sconf      *sarama.Config
	cconf      *cluster.Config
}

func GetSaramConfig(config *KafkaConfig) *sarama.Config {
	sconfig := sarama.NewConfig()
	sconfig.Version = sarama.V0_10_2_0

	sconfig.Producer.Partitioner = sarama.NewRandomPartitioner
	sconfig.Producer.RequiredAcks = sarama.WaitForAll
	sconfig.Producer.Return.Errors = true
	sconfig.Producer.Return.Successes = true
	sconfig.Producer.Flush.Bytes = 100 * 1024
	sconfig.Producer.Flush.Frequency = 1000 * time.Millisecond

	sconfig.Net.SASL.User = config.GetUser()
	sconfig.Net.SASL.Password = config.GetPassword()
	sconfig.Net.SASL.Handshake = true
	sconfig.Net.SASL.Enable = true
	return sconfig
}

func GetClusterConfig(config *KafkaConfig) *cluster.Config {
	cconf := cluster.NewConfig()
	cconf.Version = sarama.V0_10_2_0

	cconf.Consumer.Offsets.CommitInterval = 1 * time.Second
	cconf.Consumer.Offsets.Initial = sarama.OffsetNewest
	cconf.Consumer.Return.Errors = true
	cconf.ChannelBufferSize = 1024

	cconf.Group.Return.Notifications = true

	cconf.Net.SASL.User = config.GetUser()
	cconf.Net.SASL.Password = config.GetPassword()
	cconf.Net.SASL.Handshake = true
	cconf.Net.SASL.Enable = true

	return cconf
}

func NewKafkaClient(ctx context.Context, config *KafkaConfig) (KafkaClient, error) {
	if config.L == nil {
		return KafkaClient{}, ErrKafkaLoggerNil
	}
	sconfig := GetSaramConfig(config)
	cconfig := GetClusterConfig(config)
	return KafkaClient{
		baseConfig: config,
		sconf:      sconfig,
		cconf:      cconfig,
	}, nil
}

func (client *KafkaClient) NewConsumer(ctx context.Context, consumerName string) (Consumer, error) {
	consumer, err := cluster.NewConsumer(client.baseConfig.GetBrokerAddr(), consumerName, []string{client.baseConfig.GetMqTopic()}, client.cconf)
	if err != nil {
		return nil, err
	}

	c := &KafkaConsumer{
		consumer:      consumer,
		client:        client,
		closeCh:       make(chan struct{}),
		returnCloseCh: make(chan struct{}),
	}

	go c.InitConsumerReturn(ctx)

	return c, nil
}

func (client *KafkaClient) NewProducer(ctx context.Context) (Producer, error) {
	producer, err := sarama.NewAsyncProducer(client.baseConfig.GetBrokerAddr(), client.sconf)
	if err != nil {
		return nil, err
	}

	p := &KafkaProducer{
		asyncProducer: producer,
		client:        client,
		lock:          &sync.Mutex{},
		callbackMap:   make(map[int64]func(msg *sarama.ProducerMessage, e error), 1024),
		returnCloseCh: make(chan struct{}),
	}
	go p.InitProducerReturn(ctx)

	return p, nil
}
