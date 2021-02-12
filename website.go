package main

import (
	log "fmt"
	"html/template"
	"net/http"
	"time"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func log_request(r *http.Request) {
	t := time.Now()
	log.Printf("%s %s %s %s\n", t.Format("20060102150405"), r.RemoteAddr, r.Method, r.URL.Path)
}

// HandleTemplate...
func HandleTemplate(w http.ResponseWriter, r *http.Request) {
	log_request(r)
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

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", HandleTemplate)

	log.Println("Listening on :8087")
	http.ListenAndServe(":8087", nil)
}

