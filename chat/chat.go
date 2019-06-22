package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

var openRoomList = make(map[string]bool)

type client struct {
	conn *websocket.Conn
	room *room
}

func newClient(ws *websocket.Conn, room *room) *client {
	return &client{
		conn: ws,
		room: room,
	}
}
func (c *client) listen() {
	for {
		var msg string
		err := websocket.Message.Receive(c.conn, &msg)
		if err != nil {
			fmt.Println(err)
			c.room.Leave(c)
			return
		} else {
			c.room.receive <- &msg
		}
	}
}
func (c *client) write(msg *string) {
	err := websocket.Message.Send(c.conn, *msg)
	if err != nil {
		fmt.Println(err)
		fmt.Println(len(c.room.clients))
	}
}

type room struct {
	pattern  string
	messages []*string
	receive  chan *string
	join     chan *client
	leave    chan *client
	clients  map[*client]bool
}

func newRoom(pattern string) *room {
	return &room{
		pattern:  pattern,
		messages: []*string{},
		receive:  make(chan *string),
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
func (r *room) Receive(msg *string) {
	r.receive <- msg
}
func (r *room) SendPastMessages(c *client) {
	for _, msg := range r.messages {
		c.write(msg)
	}
}
func (r *room) Broadcast(msg *string) {
	layout := "(15:04:05)"
	tmp := time.Now().Format(layout)
	*msg += "\t" + tmp
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
			fmt.Println("New comer!")
			r.clients[c] = true
			r.SendPastMessages(c)
		case c := <-r.leave:
			fmt.Println("Bye someone!")
			delete(r.clients, c)
			c.conn.Close()
			if len(r.clients) == 0 {
				r.messages = nil
			}
		case msg := <-r.receive:
			fmt.Println("New message!")
			r.Broadcast(msg)
			r.messages = append(r.messages, msg)
		}
	}
}

func entryListen(ws *websocket.Conn) {
	var pattern string
	websocket.Message.Receive(ws, &pattern)
	if pattern != "" {
		room := newRoom(pattern)
		go room.listen()
	}
	websocket.Message.Send(ws, "ok")
}

func main() {
	fmt.Println("Listening on 8080")
	http.Handle("/entry", websocket.Handler(entryListen))
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":8080", nil)
}
