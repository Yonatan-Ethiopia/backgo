package main

import ( "strconv"; "net/http"; "fmt"; _ "html/template"; "github.com/julienschmidt/httprouter"; _ "strings"; _ "unicode/utf8"; "backgo/internal/validator"; "errors"; "backgo/internal/models")


type boxCreateForm struct{
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires_at          int    `form:"expires_at"`
    validator.Validator `form:"-"`
}

type userSignUpForm struct{
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

type userLoginForm struct{
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator `form:"-"`
}

type PassChangeForm struct {
    CurrentPassword     string `form:"currentPassword"`
    NewPassword         string `form:"newPassword"`
    ConfirmPassword     string `form:"newPasswordConfirmation"`
    validator.Validator `form:"-"`
}


func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
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

func (app *application) formCreateGet( w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = &boxCreateForm{
        Expires_at: 365,
    }
    app.render(w, 200, "create.tmpl", data)
}

func (app *application) formCreatePost( w http.ResponseWriter, r *http.Request){

    var form boxCreateForm
    
    err := app.decodePostForm(r, &form)
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer then 100 characters")
    form.CheckField(validator.PermittedInt(form.Expires_at, 1,7,365), "expires", "This field cannot be a value other then 1, 7 or 365")
    
    if !form.Valid(){
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
    
    app.sessionManager.Put(r.Context(),"flash","Box succefully created!")
    http.Redirect(w, r, fmt.Sprintf("/still/%d", id), http.StatusSeeOther)
}

func (app *application) boxViewGet( w http.ResponseWriter, r *http.Request){
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil{
        app.notFound(w)
    }
    
    rec, err := app.dbconn.Get(id)
    if err != nil{
        app.notFound(w)
    }
    
    tempData := app.newTemplateData(r)
    tempData.Box = rec
    app.render(w, http.StatusOK, "view.tmpl", tempData)
}


func (app *application) userSignUpGet(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = userSignUpForm{}
    app.render(w, 200, "signup.tmpl", data)
}
func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request){
    var user userSignUpForm
    
    err := app.decodePostForm(r, &user)
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    user.CheckField(validator.NotBlank(user.Name), "name", "This field cannot be blank")
    user.CheckField(validator.NotBlank(user.Email), "email", "This field cannot be blank")
    user.CheckField(validator.NotBlank(user.Password), "password", "This field cannot be blank")
    
    user.CheckField(validator.MinChars(user.Password, 8), "password", "Password must be atleast 8 characters long")
    user.CheckField(validator.ValidEmailAddress(user.Email, validator.EmailRX), "email", "Invalid email address")
    
    if !user.Valid(){
         data := app.newTemplateData(r)
         data.Form = user
         app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
         return
    }
    
    err = app.userconn.Insert(user.Name, user.Email, user.Password)
    if err != nil{
        if errors.Is(err, models.ErrDuplicateEmail){
            user.AddFieldError("email", "Email already exists!")
            data := app.newTemplateData(r)
            data.Form = user
            app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
            return
        } else {
            app.serverError(w, err)
        }
        
        return
    }
    app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please login.")
    http.Redirect(w, r, "/user/login", http.StatusSeeOther)
    
}
func (app *application) userLogInGet(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = userLoginForm{}
    
    app.render(w, http.StatusOK, "login.tmpl", data)
    
}
func (app *application) userLogInPost(w http.ResponseWriter, r *http.Request){
    
    var user userLoginForm
    var id int
    redirectURL := "/create"
    err := app.decodePostForm(r, &user)
    if err != nil{
        app.clientError(w, http.StatusBadRequest)
        return
    }
    
    user.CheckField(validator.NotBlank(user.Email), "email", "This section cannot be empty")
    user.CheckField(validator.ValidEmailAddress(user.Email, validator.EmailRX), "email", "Invalid email address")
    user.CheckField(validator.NotBlank(user.Password), "password", "This section cannot be empty")
    user.CheckField(validator.MinChars(user.Password, 8), "password", "Password needs to be atleast 8 characters long.")
    
    if !user.Valid(){
        data := app.newTemplateData(r)
        data.Form = user
        app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
        return
    }
    
    id, err = app.userconn.Authenticate(user.Email, user.Password)
    if err != nil{
        if errors.Is(err, models.ErrInvalidCredentials){
            user.AddNonFieldError("Incorrect email or password")
            data := app.newTemplateData(r)
            data.Form = user
            app.render(w, http.StatusBadRequest, "login.tmpl", data)
            return
        }
            app.serverError(w, err)
            return
    }
    err = app.sessionManager.RenewToken(r.Context())
    if err != nil{
        app.serverError(w, err)
        return
    }                                  
    
    app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
    url := app.sessionManager.PopString(r.Context(), "redirectUrl")
    if url != "" {
        redirectURL = url
    }
    http.Redirect(w, r, redirectURL, http.StatusSeeOther)
    
}
func (app *application) userLogOutPost(w http.ResponseWriter, r *http.Request){
    err := app.sessionManager.RenewToken(r.Context())
    if err != nil{
        app.serverError(w, err)
        return
    }
    
    app.sessionManager.Remove(r.Context(), "authenticatedUserId")
    
    app.sessionManager.Put(r.Context(), "flash", "You have successfully logged out! ")
    
    http.Redirect(w,r, "/", http.StatusSeeOther)
}

func(app *application) aboutPage(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    app.render(w, http.StatusOK, "about.tmpl", data)
}


func (app *application) viewAccount(w http.ResponseWriter, r *http.Request){
    id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
    flash := app.sessionManager.PopString(r.Context(), "flash")
    userDetail, err := app.userconn.Get(id)
    if err != nil{
        if errors.Is(err, models.ErrNoRecord){
             http.Redirect(w, r, "/user/login", http.StatusSeeOther)
        }
        app.serverError(w, err)
    }
    
    data := app.newTemplateData(r)
    data.Form = userDetail
    data.Flash = flash
    app.render(w, http.StatusOK, "account.tmpl", data)
    
}

func (app *application) changePassGet(w http.ResponseWriter, r *http.Request){
    data := app.newTemplateData(r)
    data.Form = &PassChangeForm{}
    
    app.render(w, http.StatusOK, "password.tmpl", data)
}

func (app *application) changePassPost(w http.ResponseWriter, r *http.Request){
    var passForm PassChangeForm
    var match bool
    id := app.sessionManager.GetInt(r.Context(), "authenticatedUserId")
    err := app.decodePostForm(r, &passForm)
    if err != nil{
        app.serverError(w, err)
    }
    passForm.CheckField(validator.NotBlank(passForm.CurrentPassword), "currentPassword", "This field cannot be blank" )
    passForm.CheckField(validator.NotBlank(passForm.NewPassword), "newPassword", "This field cannot be blank" )
    passForm.CheckField(validator.MinChars(passForm.NewPassword, 8), "newPassword", "This field cannot be blank" )
    
    passForm.CheckField(validator.NotBlank(passForm.ConfirmPassword), "newPasswordConfirmation", "This field cannot be blank" )
    
    passForm.CheckField((passForm.NewPassword == passForm.ConfirmPassword), "newPasswordConfirmation", "The passwords do not match")
    
    match, err = app.userconn.Match(id, passForm.CurrentPassword)
    if err != nil{
        fmt.Printf("Debug here!!")
        app.serverError(w, err)
    }
    passForm.CheckField(match, "currentPassword", "Wrong Password")
    
    if !passForm.Valid() {
        data := app.newTemplateData(r)
        data.Form = passForm
        app.render(w, http.StatusUnprocessableEntity,"password.tmpl", data)
        return
    }
    
    err = app.userconn.UpdatePass(id, passForm.NewPassword)
    if err != nil{
        if errors.Is(err, models.ErrInvalidCredentials){
            app.clientError(w, http.StatusBadRequest)
        }
        app.serverError(w, err)
    }
    
    app.sessionManager.Put(r.Context(), "flash", "Your password is succefully changed")
    http.Redirect(w,r, "/account/view", http.StatusSeeOther)
}
