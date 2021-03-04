package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
)

type EventDto struct {
	Type string `json:type`
	Value string `json:"value"`
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/events", getEvents).Methods("GET")
	router.HandleFunc("/events", createEvent).Methods("POST")
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

func createEvent(w http.ResponseWriter, r *http.Request) {
	var eventDto EventDto
	_ = json.NewDecoder(r.Body).Decode(&eventDto)
	fmt.Println("Received Event ", eventDto.Value, " of Type", eventDto.Type)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventDto)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventDto := EventDto{Value: "Health Check Event"}
	json.NewEncoder(w).Encode(eventDto);
}