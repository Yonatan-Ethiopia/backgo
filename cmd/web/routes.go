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
    router.Handler(http.MethodGet, "/create", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.formCreateGet))))
    router.Handler(http.MethodPost, "/create", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.formCreatePost))))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))

         
    router.Handler(http.MethodGet, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignUpGet)))
    router.Handler(http.MethodPost, "/user/signup", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userSignUpPost)))
    
    router.Handler(http.MethodGet, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogInGet)))
    router.Handler(http.MethodPost, "/user/login", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogInPost)))
    router.Handler(http.MethodPost, "/user/logout", app.sessionManager.LoadAndSave(http.HandlerFunc(app.userLogOutPost)))
    
    return router
}
