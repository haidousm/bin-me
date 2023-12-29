package main

import (
	"net/http"

	"binme.haido.us/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/bins/:id", dynamicMiddleware.ThenFunc(app.binView))
	router.Handler(http.MethodGet, "/users/signup", dynamicMiddleware.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/users/signup", dynamicMiddleware.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/users/login", dynamicMiddleware.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/users/login", dynamicMiddleware.ThenFunc(app.userLoginPost))

	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/bin/new", protectedMiddleware.ThenFunc(app.binCreate))
	router.Handler(http.MethodPost, "/bins", protectedMiddleware.ThenFunc(app.binCreatePost))
	router.Handler(http.MethodPost, "/users/logout", protectedMiddleware.ThenFunc(app.userLogoutPost))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(router)
}
