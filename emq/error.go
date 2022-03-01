package emq

import "errors"

var (
	KafkaConfigIllegalErr = errors.New("kafka config illegal")
	KafkaLoggerNilErr     = errors.New("kafka logger nil")

	ConsumerClosedErr = errors.New("consumer closed")
	IllegalMqTypeErr  = errors.New("illegal mq type")
)
