package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	Addr      string
	StaticDir string
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

	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		config:      cfg,
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	app.Serve()
}
