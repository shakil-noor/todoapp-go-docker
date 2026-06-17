package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var db *sql.DB //indicate the DB instance memory location into the sql.
var tmplt *template.Template

func main() {
	// Gather configuration from environment variables (12-Factor App rule)
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build the MySQL Data Source Name (DSN) connection string
	// Format: username:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection pool: ", err)
	}
	defer db.Close()

	// Verify connection is physically alive before starting the engine
	err = db.Ping()
	if err != nil {
		log.Fatal("Database is not reachable checked by ping ", err)
	}

	// Create database
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Table not created because: ", err)
	}

	tmplt = template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/", showTasksHandler)

	fmt.Println("🚀 Server running smoothly at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func showTasksHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	// take data from database
	rows, err := db.Query("SELECT title FROM tasks ORDER BY id DESC")
	if err != nil {
		http.Error(writer, "Database execution error during read", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []string
	for rows.Next() {
		var taskTitle string
		err := rows.Scan(&taskTitle)
		if err != nil {
			http.Error(writer, "row processing error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, taskTitle)
	}
	tmplt.Execute(writer, tasks) // send tasks as writer to http
}
