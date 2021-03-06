package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"context"
	"github.com/segmentio/kafka-go"
)

type VoteDTO struct {
	Type string `json:"type"`
	UniqueID string `json:"uniqueID"`
}

func main() {
	go initKafka()
	router := mux.NewRouter()
	router.HandleFunc("/vote", createEvent).Methods("POST")
	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowCredentials: true,
    })
	handler := c.Handler(router)
	fmt.Println("Server starting a port :3010")
	err := http.ListenAndServe(
		":3010",
		handler,
	)
	if (err != nil) {
		panic("Could Not Start Server " + err.Error())
	}
}
/** API Module **/
func createEvent(w http.ResponseWriter, r *http.Request) {
	var dto VoteDTO
	_ = json.NewDecoder(r.Body).Decode(&dto)
	fmt.Println("Received Event ", dto.UniqueID)
	ctx := context.Background()
	go sendVote(ctx, dto)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}


/** Kafka Module**/

const (
	topic = "vote_topic"
	brokerAddress = "localhost:9092"
)

func initKafka() {
	c, _ := kafka.Dial("tcp", brokerAddress)
	kt := kafka.TopicConfig{
		Topic: topic,
		NumPartitions: 1,
		ReplicationFactor: 1,
	}
	err := c.CreateTopics(kt)
	if err != nil {
		panic("Could Not Create Topic " + err.Error())
	}
	fmt.Println("Topic Created Successfully")
}

func sendVote(ctx context.Context, dto VoteDTO) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic: topic,
	})
	message := dto.UniqueID + "::" +dto.Type
	err := w.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		panic("Message could not be written to kafka " + err.Error())
	}
	fmt.Println("Voted -> ", message)
}