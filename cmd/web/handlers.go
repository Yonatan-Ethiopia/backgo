package main

import ( "strconv"; "net/http"; "encoding/json"; "fmt"; _ "html/template"; "github.com/julienschmidt/httprouter"; "strings"; "unicode/utf8")


type boxCreateForm struct{
    Title string
    Content string
    Expires_at int
    FieldErrors map[string]string
}


func (app *application) homePage(w http.ResponseWriter, r *http.Request){
    boxes,err := app.dbconn.Latest()
    if err != nil{
        app.serverError(w,err)
    }
    tempData := app.newTemplateData(r)
    tempData.Boxes = boxes
    app.render(w, 200, "home.tmpl", tempData)
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

func (app *application) formCreateGet( w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = &boxCreateForm{
        Expires_at: 365,
    }
    app.render(w, 200, "create.tmpl", data)
}

func (app *application) formCreatePost( w http.ResponseWriter, r *http.Request){
    err := r.ParseForm()
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")
    
    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    form := &boxCreateForm{
        Title: title,
        Content: content,
        Expires_at: expires,
        FieldErrors: map[string]string{},
    }
    
    if strings.TrimSpace(title) == "" {
        form.FieldErrors["title"] ="This field cannot be blank"
    }
    if utf8.RuneCountInString(title) > 100 {
        form.FieldErrors["title"] = "This field cannot be more than 100 characters"
    }
    
    if strings.TrimSpace(content) == "" {
        form.FieldErrors["content"] ="This field cannot be blank"
    }
    
    if expires != 1 && expires != 7 && expires != 365 {
        form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
    }
    
    if len(form.FieldErrors) > 0 {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
        
        return
    }
    
    id, err := app.dbconn.Insert(form.Title, form.Content, form.Expires_at)
    if err != nil{
        app.serverError(w, err)
        return
    }
    http.Redirect(w, r, fmt.Sprintf("/still/%d", id), http.StatusSeeOther)
}

func (app *application) boxViewGet( w http.ResponseWriter, r *http.Request){
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil{
        app.serverError(w, err)
    }
    
    rec, err := app.dbconn.Get(id)
    if err != nil{
        app.serverError(w, err)
    }
    
    tempData := app.newTemplateData(r)
    tempData.Box = rec
    app.render(w, 200, "view.tmpl", tempData)
}
