package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/post-message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		fmt.Print(r.Form)

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.Handle("/",
		http.FileServer(http.Dir("web")),
	)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
