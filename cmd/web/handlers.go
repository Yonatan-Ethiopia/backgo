package main

import ( "strconv"; "net/http"; "encoding/json"; "fmt")

func (app *application )home(w http.ResponseWriter, r *http.Request){
    if r.URL.Path != "/" {
        app.clientError(w, http.StatusBadRequest)
        return
    }
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

func (app *application) apicreate( w http.ResponseWriter, r *http.Request){
    if r.Method != "POST" {
        w.Header().Set("Allow", http.MethodPost)
        app.clientError(w, http.StatusMethodNotAllowed)
    }
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7
    
    id, err := app.dbconn.Insert(title, content, expires)
    if err != nil{
        app.infoLog.Printf("Not succefull")
        app.serverError(w, err)
        return
    }
    
    http.Redirect(w,r, fmt.Sprintf("/still?id=%d", id), http.StatusSeeOther)
    w.Write([]byte("Data inserted"))
}


func (app *application) greet(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "HI there this is from b")
    app.infoLog.Printf("Hi there nigros")
}

func (app *application) cnt( w http.ResponseWriter, r *http.Request){
    app.infoLog.Printf("This was called to cnt")
    value, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        app.serverError(w, err)
        return
    }
    rec, err := app.dbconn.Get(value)
    if err != nil {
        app.serverError(w, err)
        return
    }
    fmt.Println("The value is ",rec.Id)
    fmt.Fprintf(w, "This is your id: %+v", rec.Title)
}
