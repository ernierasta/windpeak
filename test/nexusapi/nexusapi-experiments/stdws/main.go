package main

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://sso.nexusmods.com/"
	url := "wss://sso.nexusmods.com"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte(fmt.Sprintf("{ \"id\": \"%s\", \"appid\": \"Automaton\" }", "45a77daa-c819-47cc-9a16-cfdc6f11ac62"))); err != nil {
		log.Fatal("write:", err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal("read:", err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
