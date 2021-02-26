package kafkawriter

import (
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// func Connect(ctx context.Context, topic string, partition int) *kafka.Conn {
// 	kconn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
// 	if err != nil {
// 		log.Fatalf("kafka connecting to DialLeader%s\n", err)
// 	}
// 	return kconn
// }

func Configure(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func CreateTopic(kafkaURL, topic string) {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
