package mq

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMQ interface {
	ReadMessage() ([]byte, []byte, error)
	Push(key, value []byte) error
}

var _ KafkaMQ = &Storage{}

type Storage struct {
	Writer  *kafka.Writer
	Reader  *kafka.Reader
	Context context.Context
}

func NewKafka(writer *kafka.Writer, reader *kafka.Reader, context context.Context) KafkaMQ {
	return &Storage{
		Writer:  writer,
		Reader:  reader,
		Context: context,
	}
}

func (s *Storage) ReadMessage() ([]byte, []byte, error) {
	m, err := s.Reader.ReadMessage(s.Context)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return m.Key, m.Value, err
}

func (s *Storage) Push(key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return s.Writer.WriteMessages(s.Context, message)
}
