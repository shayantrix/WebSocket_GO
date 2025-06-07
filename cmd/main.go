package main

import (
	"fmt"
	"net/http"
	"github.com/shayantrix/WebSocket_GO/pkg/controller"
)

func main(){
	http.HandleFunc("/ws", controller.WsHandler)
	fmt.Println("Websocket Server started on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("Error starting server: ", err)
	}
}


