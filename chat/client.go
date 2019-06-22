package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
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
			if err != io.EOF {
				log.Println(err) //EOFは無視していい
			}
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
		log.Println(err)
	}
}
