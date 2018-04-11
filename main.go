package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"io/ioutil"
)

func main() {
	http.HandleFunc("/", serv)

	http.HandleFunc("/ws", servWs)
	http.ListenAndServe(":8080", nil)

}

func serv(w http.ResponseWriter, req *http.Request) {
	var html, _ = ioutil.ReadFile("index.html")
	w.Write(html)
}

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func servWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	clients[conn] = true
	if err != nil {
		log.Println(err)
		return
	}
	for {
		mt, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for client := range clients {
			go func(client *websocket.Conn) {
				client.WriteMessage(mt, message)
				//err = conn.WriteMessage(mt, message)
				if err != nil {
					log.Println("write:", err)
					delete(clients, client)
					client.Close()
				}
			}(client)

		}
	}
}
