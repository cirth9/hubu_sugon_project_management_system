package kk

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaReader struct {
	R *kafka.Reader
}

func (r *KafkaReader) readMsg() {
	for {
		m, err := r.R.ReadMessage(context.Background())
		if err != nil {
			log.Printf("kafka readMsg err %s \n", err.Error())
			continue
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

	}
}

func (r *KafkaReader) Close() {
	r.R.Close()
}

func GetReader(brokers []string, groupId, topic string) *KafkaReader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId, //同一个组下的consumer 协同工作 共同消费topic队列中的内容
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	k := &KafkaReader{R: r}
	//go k.readMsg()
	return k
}
