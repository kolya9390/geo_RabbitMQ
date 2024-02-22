package main

import (
	"fmt"
	"gitlab.com/ptflp/gopubsub/kafkamq"
	"log"
	"time"
)

func main() {
	mq, err := kafkamq.NewKafkaMQ("localhost:9092", "myGroup")
	if err != nil {
		log.Fatalf("Failed to create MQ: %s\n", err)
	}

	topic := "myTopic"

	msgs, err := mq.Subscribe(topic)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s\n", err)
	}

	n := 1
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			if msg.Err != nil {
				log.Fatalf("Failed to receive message: %s\n", msg.Err)
			}
			fmt.Printf("Received message: %s\n", string(msg.Data))
			// Подтверждаем получение сообщения
			err = mq.Ack(&msg)
			if err != nil {
				log.Fatalf("Failed to ack message: %s\n", err)
			}
		default:
			time.Sleep(100 * time.Millisecond)
			err = mq.Publish(topic, []byte(fmt.Sprintf("Message kafka %d", n)))
			if err != nil {
				log.Fatalf("Failed to publish message: %s\n", err)
			}
			n++
		}
	}
}
