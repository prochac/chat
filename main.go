package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func main() {
	var messages []string

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

		r.ParseForm()
		if len(r.Form["message"]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no message"))
			return
		}

		messages = append(messages, r.Form["message"][0])

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.Handle("/", func() http.HandlerFunc {
		t, err := template.ParseFiles("web/index.html")
		if err != nil {
			panic(err)
		}
		return func(w http.ResponseWriter, r *http.Request) {
			data := struct {
				Messages []string
			}{
				Messages: messages,
			}

			t.Execute(w, data)
		}
	}())

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static")),
		),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
