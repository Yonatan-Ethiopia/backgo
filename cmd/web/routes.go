package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router{
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    router := httprouter.New()
    router.Handler(http.MethodGet, "/", app.sessionManager.LoadAndSave(http.HandlerFunc(app.homePage)))
    router.Handler(http.MethodGet, "/still/:id", app.sessionManager.LoadAndSave(http.HandlerFunc(app.boxViewGet)))
    router.Handler(http.MethodGet, "/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.formCreateGet)))
    router.Handler(http.MethodPost, "/create", app.sessionManager.LoadAndSave(http.HandlerFunc(app.formCreatePost)))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

    
    return router
}
