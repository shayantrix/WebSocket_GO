package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// websocket.Upgrade is used for upgrading HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//Upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Println("Error upgrading: ", err)
		return
	}
	defer conn.Close()
	// Listen for incoming messages
	for {
		//Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in reading the message: ", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)
		// Echo the message back to client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message: ", err)
			break
		}
	}
}