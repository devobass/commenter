// main
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type Comment struct {
	Id	int	`json:"id"`
	Time	string `json:"time"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

var (
	db *sql.DB
	ip, port, path string
	maxSize		int)

func delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete{
		w.Header().Set("Allow", "DELETE")

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize * 1024))

	q := r.URL.Query()

	id := q.Get("id")

	if id == "" {
		http.Error(w, "Missing Id To Delete", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id)

	if err != nil {
		slog.Error("error while deleting message", "err", err)
		http.Error(w, "could not delete message from database", http.StatusInternalServerError)
		return
	}	

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("deletion successful"))
}

func comment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize * 1024))

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
	rows, err := db.Query("SELECT * FROM messages")

	if err != nil {
		slog.Error("failed to retrieve message from db", "err", err)
		http.Error(w, "could not retrieve message", http.StatusInternalServerError)
		return
	}	

	defer rows.Close()

	for rows.Next() {
		var cmt Comment
		err = rows.Scan(&cmt.Id, &cmt.Time, &cmt.Username, &cmt.Message)

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

func createDatabase() (error) {
	fd, err := os.Create(path)

	if err != nil {
		return err
	}
	defer fd.Close()

	db, err := sql.Open("sqlite", path)

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE messages (id INTEGER PRIMARY KEY AUTOINCREMENT, time TIMESTAMP, author TEXT, data TEXT);")

	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.IntVar(&maxSize, "m", 500, "Specify the maximum message size. (KB)")
	flag.StringVar(&ip, "h", "localhost", "Specify the IP for the server to listen on.")
	flag.StringVar(&port, "p", "8000", "Specify the Port for the server to listen on.")
	flag.StringVar(&path, "d", "messages.sqlite", "Specify the SQLite database file for the server to use.")
	flag.Parse()

	listen := net.JoinHostPort(ip, port)

	var err error
	
	_, err = os.Stat(path)

	if errors.Is(err, fs.ErrNotExist) {
		slog.Info("database does not exist, creating database.")
		err = createDatabase()

		if err != nil {
			slog.Error("error while creating database", "err", err)
			os.Exit(1)
		}
	} 

	if err != nil {
		slog.Error("error while opening database", "err", err)
		os.Exit(1)
	}

	db, err = sql.Open("sqlite", path)

	if err != nil {
		slog.Error("error while opening database", "err", err)
		os.Exit(1)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		slog.Error("error while accessing database", "err", err)
		os.Exit(1)
	}

	fileServer := http.FileServer(http.Dir("./site"))
	http.HandleFunc("/comment", comment)
	http.HandleFunc("/fetch", fetch)
	http.HandleFunc("/delete", delete)
	http.Handle("/", fileServer)

	err = http.ListenAndServe(listen, nil)

	if err != nil {
		slog.Error("error while starting web server", "err", err)
		os.Exit(1)
	}
}
