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
func main() {
	room := newRoom()
	fmt.Println("Listening on 8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
