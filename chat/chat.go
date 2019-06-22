package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var openRoomList = make(map[string]bool)

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
	log.Println("Listening on 8080")
	http.Handle("/entry", websocket.Handler(entryListen))
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":8080", nil)
}
