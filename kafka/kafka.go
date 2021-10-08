package kafka

import (
	"HezzelTask/config"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
)

const (
	Topic = "addUserLogs"
)

type Kafka struct {
	Kf *kafka.Conn
	Writer *kafka.Writer
}

func Connect(cfg *config.Config) (*kafka.Conn,*kafka.Writer, error) {

	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%s",cfg.Kafka.Host,cfg.Kafka.Port))
	if err != nil {
		return nil,nil, errors.Errorf("failed to dial leader:", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%s",cfg.Kafka.Host,cfg.Kafka.Port)),
		Topic:   cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	err = createTopic(conn,Topic)
	if err != nil {
		return nil,nil,errors.Errorf("failed to create topic")
	}
	return conn,writer, nil
}

func createTopic(conn *kafka.Conn, name string) error {
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) WriteMessage(key []byte,value []byte) error {
    mes := kafka.Message{
    	Key: key,
    	Value: value,
	}

	err := k.Writer.WriteMessages(context.Background(), mes)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return nil
}