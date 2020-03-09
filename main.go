package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db       *sql.DB
	m        sync.Mutex
	conns    map[*websocket.Conn]struct{}
	upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

func addConnection(conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()
	if conns == nil {
		conns = make(map[*websocket.Conn]struct{})
	}
	conns[conn] = struct{}{}
}

func initTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages(
		id        TEXT PRIMARY KEY,
		timestamp timestamp,
		text      TEXT
	)`)
	return err
}

func publishMessage(message message) error {
	_, err := db.Exec(`INSERT INTO messages(id, timestamp, text) VALUES (?, ?, ?)`,
		message.ID,
		message.Timestamp,
		message.Text,
	)
	if err != nil {
		return err
	}

	m.Lock()
	defer m.Unlock()

	for conn := range conns {
		if err := conn.WriteJSON(message); err != nil {
			delete(conns, conn)
		}
	}
	return nil
}

func readAllMessages() ([]message, error) {
	rows, err := db.Query(`SELECT id, timestamp, text FROM messages ORDER BY timestamp DESC, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mm []message
	for rows.Next() {
		var m message
		if err := rows.Scan(&m.ID, &m.Timestamp, &m.Text); err != nil {
			return nil, err
		}
		mm = append(mm, m)
	}
	return mm, rows.Err()
}

type message struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	mm, err := readAllMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mm)
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	m := message{
		ID:        uuid.New(),
		Timestamp: time.Now(),
	}
	if err := json.NewDecoder(r.Body).Decode(&m.Text); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m.Text == "" {
		http.Error(w, "missing text message", http.StatusBadRequest)
		return
	}

	if err := publishMessage(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addConnection(conn)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "file:foo.db?loc=auto")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := initTable(); err != nil {
		panic(err)
	}

	http.HandleFunc("/get-messages", handleGetMessages)
	http.HandleFunc("/post-message", handlePostMessage)
	http.HandleFunc("/ws", handleWebsocket)
	http.HandleFunc("/script.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/script.js")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			pusher.Push("/script.js", nil)
		}
		http.ServeFile(w, r, "web/index.html")
	})

	if err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil); err != nil {
		panic(err)
	}
}
