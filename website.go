package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func logRequest(r *http.Request) {
	t := time.Now()
	fmt.Printf("%s %s %s %s\n", t.Format("20060102150405"), r.RemoteAddr, r.Method, r.URL.Path)
}

// HandleTemplate...
func HandleTemplate(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	tmpl := template.Must(template.ParseFiles("login.gtmp"))
	data := TodoPageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.Execute(w, data)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleCommunity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	fmt.Fprintf(w, "<p>id: %s</p>\n", id)

	randomID := uuid.New()
	fmt.Fprintf(w, "<p>uuid: %s</p\n", randomID)
}

func main() {
	r := mux.NewRouter()

	// server static stuff from the /static path
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", HandleTemplate)
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/community/{id}/home", handleCommunity)

	log.Println("Listening on :8087")

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8087",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	server.ListenAndServe()
}
