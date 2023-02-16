package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	staticServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", staticServer))

	return mux
}
