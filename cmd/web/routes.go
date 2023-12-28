package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/bins/:id", dynamicMiddleware.ThenFunc(app.binView))
	router.Handler(http.MethodGet, "/bin/new", dynamicMiddleware.ThenFunc(app.binCreate))
	router.Handler(http.MethodPost, "/bins", dynamicMiddleware.ThenFunc(app.binCreatePost))

	router.Handler(http.MethodGet, "/users/signup", dynamicMiddleware.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/users/signup", dynamicMiddleware.ThenFunc(app.userSignupPost))

	router.Handler(http.MethodGet, "/users/login", dynamicMiddleware.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/users/login", dynamicMiddleware.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/users/logout", dynamicMiddleware.ThenFunc(app.userLogoutPost))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standardMiddleware.Then(router)
}
