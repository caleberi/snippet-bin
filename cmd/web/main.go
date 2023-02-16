package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Addr      string
	StaticDir string
	dns       string
}

type application struct {
	config      *config
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func (app *application) Serve() {
	srv := &http.Server{
		Addr:     app.config.Addr,
		ErrorLog: app.errorLogger,
		Handler:  app.routes(),
	}

	app.infoLogger.Printf("Starting server on :%s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		app.errorLogger.Fatal(err)
	}
}

func main() {

	cfg := new(config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "http network address")
	flag.StringVar(&cfg.StaticDir, "staticDir", "./ui/static/", "static directory")
	flag.StringVar(&cfg.dns, "dns", "caleberi:test1234@/snippet_bin?parseTime=true", "database connection url")

	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDatabase(cfg.dns)

	if err != nil {
		errorLogger.Fatal(err)
	}

	defer db.Close()

	app := &application{
		config:      cfg,
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	app.Serve()
}

func openDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
