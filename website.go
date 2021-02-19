package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/toml"
)

type TomlConfig struct {
	HTTP struct {
		Enabled   bool   `toml:"enabled"`
		IPAddress string `toml:"ip"`
		Port      int    `toml:"port"`
	} `toml:"http"`

	HTTPS struct {
		Enabled   bool   `toml:"enabled"`
		Port      int    `toml:"port"`
		IPAddress string `toml:"ip"`
		Cert      string `toml:"certificate"`
		Key       string `toml:"privatekey"`
	} `toml:"https"`

	DB struct {
		Type   string `toml:"type"`
		File   string `toml:"file"`
		Host   string `toml:"host"`
		User   string `toml:"username"`
		Pass   string `toml:"password"`
		Schema string `toml:"schema"`
	} `toml:"database"`
}

var config TomlConfig

func load_config() {
	var cfgfile string
	if len(os.Args) == 2 {
		cfgfile = os.Args[1]
	} else {
		cfgfile = "config.toml"
	}

	f, err := os.Open(cfgfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		panic(err)
	}
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	MainTitle string
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

	files := []string{
		"template/home.html",
		"template/header.html",
		"template/footer.html",
	}

	tmpl, err := template.ParseFiles(files...)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data := TodoPageData{
		MainTitle: "Test title",
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
	id, err := strconv.ParseInt(params["id"], 10, 64)
	fmt.Fprintf(w, "<p>id: %d</p>\n", id)

	db, err := sql.Open("sqlite3", "./website.db")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("SELECT name FROM stuff WHERE id = ? LIMIT 1")
	defer stmt.Close()

	row, err := stmt.Query(id)
	if err != nil {
		fmt.Printf("ERR: %s\n", err.Error())
		return
	}

	defer row.Close()

	for row.Next() {
		var name string
		row.Scan(&name)
		fmt.Println(name)
	}

	randomID := uuid.New()
	fmt.Fprintf(w, "<p>uuid: %s</p\n", randomID)
}

func main() {
	load_config()
	r := mux.NewRouter()

	// server static stuff from the /static path
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", HandleTemplate)
	r.HandleFunc("/login", handleLogin)
	r.HandleFunc("/community/{id}/home", handleCommunity)

	server := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", config.HTTP.IPAddress, config.HTTP.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %s (http)\n", server.Addr)
	server.ListenAndServe()
}
