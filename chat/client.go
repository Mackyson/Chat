package main

import (
	"fmt"
	"golang.org/x/net/websocket"
)

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
