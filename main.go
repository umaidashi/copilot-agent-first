package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Task represents a todo task
type Task struct {
	Title   string    `json:"title"`
	DueDate time.Time `json:"due_date"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		panic(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"due_date" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "hello world"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO tasks (title, due_date) VALUES (?, ?)", task.Title, task.DueDate.Format(time.RFC3339))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT title, due_date FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate string
		if err := rows.Scan(&task.Title, &dueDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	initDB()

	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/tasks", createTaskHandler)
	http.HandleFunc("/tasks/list", listTasksHandler)
	http.ListenAndServe(":8080", nil)
}
