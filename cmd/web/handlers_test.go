package main

import (
        "net/http"
        "testing"
        "net/url"
        "strconv"
        
        "backgo/internal/assert"
    )
    
func TestPing(t *testing.T){
    
    app := newTestApplication(t)
    
    ts := newTestServer(t, app.routes())
    defer ts.Close()
    
    code, _, body := ts.get(t, "/ping")
    
    assert.Equal(t, code, http.StatusOK)
    assert.Equal(t, body, "OK")
        
}

func TestBoxView(t *testing.T){
    
    app := newTestApplication(t)
    
    ts := newTestServer(t, app.routes())
    defer ts.Close()
    
    tests := []struct{
            name        string
            urlPath     string
            wantCode    int
            wantBody    string
        
        }{
            {
                name: "Valid ID",
                urlPath: "/still/1",
                wantCode: http.StatusOK,
                wantBody: "An old silent pond...",
            },
            {
                name: "Non-existent id",
                urlPath: "/still/2",
                wantCode: http.StatusNotFound,
            },
            {
                name: "Negative ID",
                urlPath: "/still/-1",
                wantCode: http.StatusNotFound,
            },
            {
                name: "Decimal ID",
                urlPath: "/still/1.23",
                wantCode: http.StatusNotFound,
            },
            {
                name: "String ID",
                urlPath: "/still/foo",
                wantCode: http.StatusNotFound,
            },
            {
                name: "Empty ID",
                urlPath: "/still/",
                wantCode: http.StatusNotFound,
            },
            
        }
        
    for _,tt := range tests{
        
        t.Run(tt.name, func(t *testing.T){
            
            code, _, body := ts.get(t, tt.urlPath)
            
            assert.Equal(t, code, tt.wantCode)
            
            if tt.wantBody != ""{
                assert.StringContains(t, body, tt.wantBody)
            }
            
        })
        
    }
}

func TestUserSignup(t *testing.T){
    
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()
    
    _,_, body := ts.get(t, "/user/signup")
    validCSRFToken := extractCSRFToken(t, body)
    t.Log("CSRF:", validCSRFToken)
    const (
        validName = "Bob"
        validPassword = "validPassword"
        validEmail = "due@gmail.com"
        formTag = "<form action='/user/signup' method='POST' novalidate>"
    )
    
    tests := []struct {
        name         string
        userName     string
        userEmail    string
        userPassword string
        csrfToken    string
        wantCode     int
        wantFormTag  string
    }{
        {
            name:         "Valid submission",
            userName:     validName,
            userEmail:    validEmail,
            userPassword: validPassword,
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusSeeOther,
        },
        {
            name:         "Invalid CSRF Token",
            userName:     validName,
            userEmail:    validEmail,
            userPassword: validPassword,
            csrfToken:    "wrongToken",
            wantCode:     http.StatusBadRequest, // Adjusted assuming standard CSRF failure code
        },
        {
            name:         "Empty name",
            userName:     "",
            userEmail:    validEmail,
            userPassword: validPassword,
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
        {
            name:         "Empty email",
            userName:     validName,
            userEmail:    "",
            userPassword: validPassword,
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
        {
            name:         "Empty password",
            userName:     validName,
            userEmail:    validEmail,
            userPassword: "",
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
        {
            name:         "Invalid email",
            userName:     validName,
            userEmail:    "bob@example.",
            userPassword: validPassword,
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
        {
            name:         "Short password",
            userName:     validName,
            userEmail:    validEmail,
            userPassword: "pa$$",
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
        {
            name:         "Duplicate email",
            userName:     validName,
            userEmail:    "dupe@example.com",
            userPassword: validPassword,
            csrfToken:    validCSRFToken,
            wantCode:     http.StatusUnprocessableEntity,
            wantFormTag:  formTag,
        },
    }
    
    for _, tt := range tests{

        
        t.Run(tt.name, func(t *testing.T){
            form := url.Values{}
            form.Add("name", tt.userName)
            form.Add("email", tt.userEmail)
            form.Add("password", tt.userPassword)
            form.Add("csrf_token", tt.csrfToken)
            code, _, body := ts.postForm(t, "/user/signup", form)
            
            assert.Equal(t, code, tt.wantCode)
            if tt.wantFormTag != ""{
                assert.StringContains(t, body, tt.wantFormTag )
            }
        })
    }
}

func TestFormCreatevalid(t *testing.T) {
    app := newTestApplication(t)
    
    ts := newTestServer(t, app.routes())
    defer ts.Close()
    
    _, _, body := ts.get(t, "/create")
    validCSRFToken := extractCSRFToken(t, body)
    validToken := validCSRFToken
    
    const(
        validTitle="The poems of a test"
        validContent="Never trust, Always test"
        validExpiry = 365
    )
    
    tests := []struct{
        name string
        title string
        content string
        expires int
        csrfToken string
        expected int
    }{
       {
            name: "Valid Submission",
            title: validTitle,
            content: validContent,
            expires: validExpiry,
            csrfToken: validToken,
            expected: http.StatusSeeOther,
        }, 
       {
            name: "Empty title",
            title: "",
            content: validContent,
            expires: validExpiry,
            csrfToken: validToken,
            expected: http.StatusUnprocessableEntity,
        }, 
       {
            name: "Empty content",
            title: validTitle,
            content: "",
            expires: validExpiry,
            csrfToken: validToken,
            expected: http.StatusUnprocessableEntity,
        }, 
       {
            name: "Invalid Expiry",
            title: validTitle,
            content: validContent,
            expires: 0,
            csrfToken: validToken,
            expected: http.StatusUnprocessableEntity,
        }, 
       {
            name: "Invalid Token",
            title: validTitle,
            content: validContent,
            expires: validExpiry,
            csrfToken: "",
            expected: http.StatusUnprocessableEntity,
            
        }, 
    }
    
    for _, tt := range tests{
        t.Run(tt.name, func(t *testing.T){
            
            form := url.Values{}
            form.Add("title", tt.title)
            form.Add("content", tt.content)
            form.Add("expires_at", strconv.Itoa(tt.expires))
            form.Add("csrf_token", tt.csrfToken)
            
            code, _, _ := ts.postForm(t, "create", form)
            
            assert.Equal(t, code, tt.expected)
            
            
        })
    }
    
    
}

func TestFormCreate(t *testing.T){
    app := newTestApplication(t)
    
    ts := newTestServer(t, app.routes())
    defer ts.Close()
    t.Run("Invalid Account", func(t *testing.T){
        
        code, header, _ := ts.get(t, "/create")
        
        assert.Equal(t, code, http.StatusSeeOther)
        
        assert.Equal(t, header.Get("Location"), "/user/login")
        
    })
    
    t.Run("Valid", func(t *testing.T){
        _, _, body := ts.get(t, "/user/loging")
        
        csrfToken := extractCSRFToken(t, body)
        
        form := url.Values{}
        form.Add("email", "dupe@exmaple.com")
        form.Add("password", "validPassword")
        form.Add("csrf_token", csrfToken)
        
        ts.postForm(t, "/user/login", form)
        
        code, _, body := ts.get(t, "/create")
        
        assert.Equal(t, code, http.StatusOK)
        assert.StringContains(t, body, "<form action='/create' method='POST'>")
    })
}
