package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestWriteKafka(t *testing.T) {

	topic := "weimi"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	require.NoError(t, err)

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one")},
		kafka.Message{Value: []byte("two")},
		kafka.Message{Value: []byte("three")},
	)

	require.NoError(t, err)
}

func TestReadKafka(t *testing.T) {
	topic := "weimi"
	partition := 0

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{
			"localhost:9092",
			"localhost:9093",
			"localhost:9094",
			"localhost:9095",
		},
		Topic:       topic,
		Partition:   partition,
		StartOffset: kafka.SeekStart,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			println("Error!", err.Error())
			break
		}
		fmt.Println("Read message success:", string(m.Key), string(m.Value))
	}
}
