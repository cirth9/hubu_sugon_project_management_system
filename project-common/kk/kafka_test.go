package kk

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	w := GetWriter("localhost:9092")
	m := make(map[string]string)
	m["projectCode"] = "1200"
	bytes, _ := json.Marshal(m)
	w.Send(LogData{
		Topic: "msproject_log",
		Data:  bytes,
	})
	time.Sleep(2 * time.Second)
}

func TestConsumer(t *testing.T) {
	reader := GetReader([]string{"localhost:9092"}, "group1", "log")
	reader.readMsg()
}
