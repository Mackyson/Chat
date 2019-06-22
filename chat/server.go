package main

import (
	"log"
)

type room struct {
	pattern  string
	messages []*Message
	receive  chan *Message
	join     chan *client
	leave    chan *client
	clients  map[*client]bool
}

func newRoom(pattern string) *room {
	return &room{
		pattern:  pattern,
		messages: []*Message{},
		receive:  make(chan *Message),
		join:     make(chan *client),
		leave:    make(chan *client),
		clients:  map[*client]bool{},
	}
}

func (r *room) Join(c *client) {
	r.join <- c
}
func (r *room) Leave(c *client) {
	r.leave <- c
}
func (r *room) Receive(msg *Message) {
	r.receive <- msg
}
func (r *room) SendPastMessages(c *client) {
	for _, msg := range r.messages {
		c.write(msg)
	}
}
func (r *room) Broadcast(msg *Message) {
	for c := range r.clients {
		c.write(msg)
	}
}
func (r *room) listen() {
	for {
		select {
		case c := <-r.join:
			log.Printf("New User %s Joined %s", c.name, r.pattern)
			r.clients[c] = true
			r.SendPastMessages(c)
		case c := <-r.leave:
			log.Printf("%s left %s", c.name, r.pattern)
			delete(r.clients, c)
			c.conn.Close()
			if len(r.clients) == 0 {
				r.messages = nil
			}
		case msg := <-r.receive:
			log.Printf("%s says %s", msg.Name, msg.Payload)
			r.Broadcast(msg)
			r.messages = append(r.messages, msg)
		}
	}
}
