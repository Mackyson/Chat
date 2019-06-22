package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var openRoomList = make(map[string]*room)
var nameCh = make(chan string, 1) //1以上にしないと二回目の送信が永遠にブロックされる

type ClientInfo struct {
	Pattern string `json:pattern`
	Name    string `json:name`
}

func entryListen(ws *websocket.Conn) {
	var clientInfo ClientInfo
	websocket.JSON.Receive(ws, &clientInfo)
	pattern := clientInfo.Pattern
	name := clientInfo.Name
	go func() {
		nameCh <- name //チャネルを経由してハンドル時に名前を決定
	}()
	var r *room
	if _, ok := openRoomList[pattern]; !ok {
		log.Printf("%s is opened", pattern)
		r = newRoom(string(pattern))
		openRoomList[pattern] = r
		connect := func(ws *websocket.Conn) {
			userName := <-nameCh
			client := newClient(ws, r, userName)
			r.Join(client)
			client.listen()
		}
		http.Handle("/"+pattern, websocket.Handler(connect))
		go r.listen()
	}
	websocket.Message.Send(ws, "ok")
}

func main() {
	log.Println("Listening on 8080")
	http.Handle("/entry", websocket.Handler(entryListen))
	http.Handle("/", http.FileServer(http.Dir("../static")))
	http.ListenAndServe(":8080", nil)
}
