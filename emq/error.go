package emq

import "errors"

var (
	ErrKafkaConfigIllegal = errors.New("kafka config illegal")
	ErrKafkaLoggerNil     = errors.New("kafka logger nil")

	ErrConsumerClosed = errors.New("consumer closed")
	ErrIllegalMqType  = errors.New("illegal mq type")
)
