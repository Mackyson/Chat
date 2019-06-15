package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type client struct {
	conn  *websocket.Conn
	msgCh chan *string
	room  *room
}

func newClient(ws *websocket.Conn, room *room) *client {
	return &client{
		conn:  ws,
		msgCh: make(chan *string, 256),
		room:  room,
	}
}
type room struct {
	messages []*string
	receive  chan *string
	join     chan *client
	leave    chan *client
	clients  map[*client]bool
}

func newRoom() *room {
	return &room{
		messages: []*string{},
		receive:  make(chan *string),
		join:     make(chan *client),
		leave:    make(chan *client),
		clients:  make(map[*client]bool),
	}
}

func (r *room) Join(c *client) {
	r.join <- c
}
func (r *room) Leave(c *client) {
	r.leave <- c
}
func (r *room) Receive(msg *string) {
	r.receive <- msg
}
func (r *room) SendPastMessages(c *client) {
	for _, msg := range r.messages {
		c.write(msg)
	}
}
func (r *room) Broadcast(msg *string) {
	for c := range r.clients {
		c.write(msg)
	}
}
func (r *room) listen() {
	for {
		select {
		case c := <-r.join:
			fmt.Println("New comer!")
			r.clients[c] = true
			r.SendPastMessages(c)
		case c := <-r.leave:
			fmt.Println("Bye someone!")
			delete(r.clients, c)
			close(c.msgCh)
		case msg := <-r.receive:
			fmt.Println("New Message!")
			r.Broadcast(msg)
			r.messages = append(r.messages, msg)
		}
	}
}

func main() {
	room := newRoom()
	fmt.Println("Listening on 8080")
	go room.listen()
	fmt.Println(http.ListenAndServe(":8080", nil))
}