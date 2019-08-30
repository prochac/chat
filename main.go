package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	messages := make([]string, 0)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket")); err != nil {
			log.Println(err)
			return
		}
	})

	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
