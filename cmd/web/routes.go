package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router{
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    router := httprouter.New()
    router.HandlerFunc(http.MethodGet, "/", app.homePage)
    router.HandlerFunc(http.MethodGet, "/still/:id", app.boxViewGet)
    router.HandlerFunc(http.MethodGet, "/create", app.greet)
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))
    
    return router
}
