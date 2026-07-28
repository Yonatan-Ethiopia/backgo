package main

import ( "strconv"; "net/http"; "encoding/json"; "fmt"; "html/template"; "github.com/julienschmidt/httprouter")


func (app *application) homePage(w http.ResponseWriter, r *http.Request){
    boxes,err := app.dbconn.Latest()
    if err != nil{
        app.serverError(w,err)
    }
    tempData := &templateData{
        Boxes: boxes,
    }
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.html",
        "./ui/html/pages/home.tmpl",
        }
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }
    
    err = ts.ExecuteTemplate(w, "base", tempData)
    if err != nil {
        app.serverError(w, err)
        return
    }
}

func (app *application )home(w http.ResponseWriter, r *http.Request){

    boxes, err := app.dbconn.Latest()
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    for _, box := range boxes{
        fmt.Fprintf(w, "%+v\n", box)
    }
    
}

func (app *application )apiview(w http.ResponseWriter, r *http.Request){
    print("here")
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1{
        http.NotFound(w, r)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    //fmt.Fprintf(w, "The id is %d", id)
    json.NewEncoder(w).Encode(map[string]int{"id":id})
}

func (app *application) apiCreatePost( w http.ResponseWriter, r *http.Request){
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7
    
    id, err := app.dbconn.Insert(title, content, expires)
    if err != nil{
        app.infoLog.Printf("Not succefull")
        app.serverError(w, err)
        return
    }
    
    http.Redirect(w,r, fmt.Sprintf("/still/%d", id), http.StatusSeeOther)
    w.Write([]byte("Data inserted"))
}


func (app *application) greet(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "HI there this is from b")
    app.infoLog.Printf("Hi there nigros")
}

func (app *application) cnt( w http.ResponseWriter, r *http.Request){
    params := httprouter.ParamsFromContext(r.Context())
    app.infoLog.Printf("This was called to cnt")
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil {
        app.serverError(w, err)                                  
        return
    }
    rec, err := app.dbconn.Get(id)
    if err != nil {
        app.serverError(w, err)
        return
    }
    fmt.Println("The value is ",rec.Id)
    fmt.Fprintf(w, "This is your id: %+v", rec.Title)
}

func (app *application) formCreate( w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is a place holder")
}

func (app *application) boxViewGet( w http.ResponseWriter, r *http.Request){
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
    
    rec, err := app.dbconn.Get(id)
    
    tempData := &templateData{
        Box: rec,
    }
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.html",
        "./ui/html/pages/view.html",
    }
    
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.errLog.Print(err)
        app.serverError(w,err)
    }
    err = ts.ExecuteTemplate(w, "base", tempData)
    
    if err != nil {
        app.errLog.Print(err)
        app.serverError(w, err)
    }
}
