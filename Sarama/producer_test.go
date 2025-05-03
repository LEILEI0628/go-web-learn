package Sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addrs = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	// kafka-console-consumer -topic=TestTopic -brokers=localhost:9094
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	// 使用Partitioner来保证同一个业务消息一定发送到同一个分区上（保证业务内信息有序）
	cfg.Producer.Partitioner = sarama.NewHashPartitioner // 默认是哈希
	producer, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "TestTopic",
		Key:   sarama.StringEncoder("oid-123"),
		Value: sarama.StringEncoder("测试消息AAA"),
		Headers: []sarama.RecordHeader{
			{Key: []byte("trace_id"),
				Value: []byte("123"),
			},
		},
		Metadata: "Metadata仅作用于发送过程",
	})
	assert.NoError(t, err)
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(addrs, cfg)
	assert.NoError(t, err)
	producer, err := sarama.NewAsyncProducerFromClient(client)
	assert.NoError(t, err)
	msgCh := producer.Input()
	msgCh <- &sarama.ProducerMessage{
		Topic: "TestTopic",
		Key:   sarama.StringEncoder("oid-123456"),
		Value: sarama.StringEncoder("测试消息BBB"),
		Headers: []sarama.RecordHeader{
			{Key: []byte("trace_id"),
				Value: []byte("123456"),
			},
		},
		Metadata: "Metadata仅作用于发送过程",
	}
	errCh := producer.Errors()
	succCh := producer.Successes()
	select {
	case err := <-errCh:
		t.Log("发送出错", err.Err)
	case <-succCh:
		t.Log("发送成功")
	}
}
