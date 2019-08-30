package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	messages := make([]string, 0)

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

	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
