package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

const TIME_LAYOUT = "15:04:05"

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
		var msg *Message
		err := websocket.JSON.Receive(c.conn, &msg)
		if err != nil {
			fmt.Println(err)
			c.room.Leave(c)
			return
		} else {
			msg.SetTime(time.Now().Format(TIME_LAYOUT))
			c.room.receive <- msg
		}
	}
}
func (c *client) write(msg *Message) {
	err := websocket.JSON.Send(c.conn, *msg)
	if err != nil {
		fmt.Println(err)
		fmt.Println(len(c.room.clients))
	}
}
