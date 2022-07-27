package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	var filename string
	pathPrefix := os.Getenv("WS_PATH_PREFIX")
	session := uuid.New()
	if err := os.Mkdir(fmt.Sprintf("%s%s", pathPrefix, session), 0750); err != nil {
		log.Println(err)
	} else {
		log.Println("Created dir")
	}
	for {
		msgType, p, err := conn.ReadMessage()
		fmt.Println("Message, type:", msgType)
		switch msgType {
		case 1:
			fmt.Println("String")
			filename = string(p)
		case 2:
			fmt.Println("Binary")
			if err := ioutil.WriteFile(fmt.Sprintf("%s%s/%s", pathPrefix, session,
				filename), p, 0644); err != nil {
				log.Println(err)
			}
			//if err := ioutil.WriteFile(fmt.Sprintf("/Volumes/ram/%s", filename), p, 0644); err != nil {
			//	log.Println(err)
			//}
		}
		if err != nil {
			log.Println(err)
			return
		}

		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	log.Println(err)
		//	return
		//}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("WS_PORT")), nil))
}
