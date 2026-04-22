package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	snippets "github.com/pauljomy/snippetbox_backend/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *snippets.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "Http network address")
	dsn := flag.String("dsn", "web:pan123jomy@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{logger: logger, snippets: &snippets.SnippetModel{DB: db}}

	logger.Info("Starting server on %s", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
