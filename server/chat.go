package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"strconv"
)

func main() {
	const port = 8080
	fmt.Println("Listening on", port)
	http.Handle("/echo", websocket.Handler(echo))
	//http.HandleFunc("/", hello)
	fmt.Println(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
func echo(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}
