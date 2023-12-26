package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)

	mux.HandleFunc("/bin/view", app.binView)
	mux.HandleFunc("/bin/create", app.binCreate)
	return app.logRequest(secureHeaders(mux))
}
