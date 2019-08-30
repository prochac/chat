package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	messages := make([]string, 0)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connections := make([]*websocket.Conn, 0)

	http.HandleFunc("/get-messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		json.NewEncoder(w).Encode(messages)
	})

	http.HandleFunc("/post-message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var message string
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		messages = append(messages, message)

		for _, conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println(err)
			}
		}

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		connections = append(connections, conn)
	})

	http.Handle("/", http.FileServer(http.Dir("web")))

	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
