package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kalpitf1/job_scheduler/backend/models"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan *models.Job, 100000)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg models.Job
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func HandleMessages() {
	for {
		job := <-Broadcast

		for client := range clients {
			err := client.WriteJSON(job)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
