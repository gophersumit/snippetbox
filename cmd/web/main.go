package web

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // New import
	"gophersumit.com/snippetbox/internal/models"
)

var (
	host     = "snippetbox-db.mysql.database.azure.com"
	database = "snippetbox"
	user     = "snpadmin"
	password = "Random@12345"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func Start() {

	addr := flag.String("addr", ":8080", "HTTP Network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, host, database)
	fmt.Println(connectionString)
	db, err := openDB(connectionString)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	infoLog.Printf("Starting server on port in Azure Container Service %s", *addr)

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
