package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB //indicate the DB instance memory location into the sql.
var tmplt *template.Template
var notFoundTmpl *template.Template

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
	// Try connecting up to 5 times
	for i := 1; i <= 5; i++ {
		log.Printf("Connecting to database (Attempt %d/5)...", i)
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to the database!")
				break
			}
		}

		log.Printf("Database not ready yet, waiting 3 seconds... Error: %v", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to database after 10 attempts: %v", err)
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
	notFoundTmpl = template.Must(template.ParseFiles("templates/404.html"))

	http.HandleFunc("/", showTasksHandler)
	http.HandleFunc("/addTask", addTasksHandler)

	fmt.Println("🚀 Server running smoothly at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func showTasksHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		renderCustom404(writer, request)
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
			http.Error(writer, "Row processing error. Cause: ", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, taskTitle)
	}
	tmplt.Execute(writer, tasks) // send tasks as writer to http
}

func renderCustom404(writer http.ResponseWriter, request *http.Request) {
	// Step A: Explicitly set the response status code header to 404
	writer.WriteHeader(http.StatusNotFound)

	// Step B: Render the beautiful 404 HTML file into the writer
	notFoundTmpl.Execute(writer, nil)
}

func addTasksHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		taskTitle := request.FormValue("task")
		if taskTitle != "" {
			insertSQL := "INSERT INTO tasks (title) VALUE (?)"
			_, err := db.Exec(insertSQL, taskTitle)
			if err != nil {
				http.Error(writer, "Data not saved because of ", http.StatusInternalServerError)
				return
			}
		}
	}
	http.Redirect(writer, request, "/", http.StatusSeeOther) // redirect to the home page after successfully add data
}
