package controller

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

var clients = make(map[*websocket.Conn]bool) //Connected clients
var broadcast = make(chan []byte)	//Broadcast channel
var mutex = &sync.Mutex{}		//Protect client map

// websocket.Upgrade is used for upgrading HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		return true //origin == "<https://yourdomain.com>" 
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	//Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Listen for incoming messages
	for {
		//Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in reading the message: ", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- message
	}
}

func HandleConnection(conn *websocket.Conn){
	//websocket.Conn := a websocket connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil{
			fmt.Println("Error reading message: ", err)
			break
		}
		fmt.Printf("Received: %s\n", message)

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil{
			fmt.Println("Error writing message: ", err)
			break
		}
	}
}

func HandleMessages(){
	for {
		//Grab the next message from the broadcast channel
		message := <- broadcast
		
		//Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

