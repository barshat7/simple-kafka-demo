package main

import (
	"context"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

const (
	brokers = "localhost:9092"
	sourceTopic = "vote_topic"
	destinationTopic = "up_vote_topic"
)

func initKafka() {
	c, _ := kafka.Dial("tcp", brokers)
	kt := kafka.TopicConfig{
		Topic: destinationTopic,
		NumPartitions: 1,
		ReplicationFactor: 1,
	}
	err := c.CreateTopics(kt)
	if err != nil {
		panic("Could Not Create Topic " + err.Error())
	}
	log.Println("Topic Created Successfully")
}

type producer func(ctx context.Context, topic, message string)

func producerImpl(ctx context.Context, topic, message string) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokers},
		Topic: topic,
	})
	err := w.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		panic("Could not write message " + err.Error())
	}
	log.Println("Message Sent to Upvote QUeue")
}

func main() {
	log.Println("Starting Upvote Processor")
	initKafka()
	ctx := context.Background()
	consume(ctx, producerImpl)
}

func consume(ctx context.Context, p producer) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string {brokers},
		Topic: sourceTopic,
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("Could Not Read Message " + err.Error())
		}
		message := string(msg.Value)
		isUpVote, _message := isUpvote(message)
		if isUpVote {
			log.Printf("Sending Upvote To -> %s", _message)
			newCtx := context.Background()
			p(newCtx, destinationTopic, _message)
		}
	}
}

func isUpvote(message string) (bool, string ){
	messageParts := strings.Split(message, "::")
	if len(messageParts) > 2 {
		return false, ""		
	}
	if messageParts[1] == "up" {
		return true, messageParts[0]
	}
	return false, ""
}

