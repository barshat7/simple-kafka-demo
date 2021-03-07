package sse

import (
	"fmt"
	"log"
	"net/http"
)

// Broker The Backbone of our SSE.
/**	It holds all open client connections
	Listens to incoming events on Notifier chan
	Broadcasts event to all connected clients
**/
type Broker struct {
	Notifier chan []byte

	newClients chan chan []byte

	closingClients chan chan []byte

	clients map[chan []byte]bool

}
// RegisterClient Handler Function that registers the client to subscribe for SSE
func (b *Broker) RegisterClient(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	messageChan := make(chan []byte)
	b.newClients <- messageChan
	defer func() {
		b.closingClients <- messageChan
	}()
	notify := r.Context().Done()
	go func() {
		<-notify
		b.closingClients <- messageChan
	}()
	for {
		fmt.Fprintf(w, "data: %s\n\n", <-messageChan)
		flusher.Flush()
	}
}

// listen Will listen to messages consumed from Kafka Learderboard topic
func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
			log.Printf("Client Added, Total Client %d", len(b.clients))
		case s := <-b.closingClients:
			delete(b.clients, s)
			log.Printf("Removed Client")
		case event := <-b.Notifier:
			for clientMessageChan := range b.clients {
				clientMessageChan <- event
			}
		}
	}
}

// SendMessage Will send the SSE
func (b *Broker) SendMessage(message string) {
	log.Println("Sending SSE")
	b.Notifier <- []byte(message)
}
// New Create new Broker
func New() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()
	return
}



