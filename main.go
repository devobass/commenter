// main
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"time"
	"log/slog"
	"net"
	"net/http"

	_ "modernc.org/sqlite"
)

type Comment struct {
	Time	string `json:"time"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var (
	db *sql.DB
	ip, port, path string
)

func comment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		slog.Error("error while parsing form", "err", err)
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)	
		return
	}

	form := r.Form

	if form == nil {
		slog.Error("form is not application/x-www-form-urlencoded", "err", err)
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)	
		return
	}

	username := form.Get("username")
	message := form.Get("message")
	time := time.Now().Unix()

	if username == "" || message == "" {
		http.Error(w, "Missing Argument", http.StatusBadRequest)	
		return
	}

	_, err = db.Exec("INSERT INTO messages(time, author, data) VALUES (?, ?, ?)", time, username, message)

	if err != nil {
		slog.Error("error while inserting to database", "err", err)
		http.Error(w, "could not insert message into database", http.StatusInternalServerError)
		return
	}	

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("insert successful"))
}

func fetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data := []Comment{}
	rows, err := db.Query("SELECT time, author, data FROM messages")

	if err != nil {
		slog.Error("failed to retrieve message from db", "err", err)
		http.Error(w, "could not retrieve message", http.StatusInternalServerError)
		return
	}	

	defer rows.Close()

	for rows.Next() {
		var cmt Comment
		err = rows.Scan(&cmt.Time, &cmt.Username, &cmt.Message)

		if err != nil {
			slog.Error("error occurred while reading data from db", "err", err)
			http.Error(w, "could not retrieve message", http.StatusInternalServerError)
			break;
		}

		if cmt.Username == "" || cmt.Message == "" {
			continue;
		}
		
		data = append(data, cmt)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("failed to retrieve message from db", "err", err)
		http.Error(w, "could not retrieve message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		slog.Error("failed to encode message data in json", "err", err)
	}
}

func main() {
	flag.StringVar(&ip, "h", "localhost", "Specify the IP for the server to listen on.")
	flag.StringVar(&port, "p", "8000", "Specify the Port for the server to listen on.")
	flag.StringVar(&path, "d", "./messages.sqlite", "Specify the SQLite database file for the server to use.")
	flag.Parse()

	listen := net.JoinHostPort(ip, port)

	var err error
	
	db, err = sql.Open("sqlite", path)

	if err != nil {
		slog.Info("error while opening database", "err", err)
		return
	}

	err = db.Ping()

	if err != nil {
		slog.Info("error while accessing database", "err", err)
		return
	}

	defer db.Close()

	fileServer := http.FileServer(http.Dir("./site"))
	http.HandleFunc("/comment", comment)
	http.HandleFunc("/fetch", fetch)
	http.Handle("/", fileServer)

	err = http.ListenAndServe(listen, nil)

	if err != nil {
		slog.Error("error while starting web server", "err", err)
		return
	}
}
