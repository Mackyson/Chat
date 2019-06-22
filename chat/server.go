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
	connect := func(ws *websocket.Conn) {
		client := newClient(ws, r)
		r.Join(client)
		client.listen()
	}
	if !openRoomList[r.pattern] {
		openRoomList[r.pattern] = true
		fmt.Println(r.pattern, "opened")
		http.Handle("/"+r.pattern, websocket.Handler(connect))
	}
	for {
		select {
		case c := <-r.join:
			log.Printf("New User %s Joined", c.name)
			r.clients[c] = true
			r.SendPastMessages(c)
		case c := <-r.leave:
			log.Printf("%s left chat", c.name)
			delete(r.clients, c)
			c.conn.Close()
			if len(r.clients) == 0 {
				r.messages = nil
			}
		case msg := <-r.receive:
			log.Printf("%s : %s", msg.Name, msg.Payload)
			r.Broadcast(msg)
			r.messages = append(r.messages, msg)
		}
	}
}
