package main

import (
        "log"
        "io"
        "net/http"
        "net/http/httptest"
        "net/http/cookiejar"
        "testing"
        "bytes"
        "github.com/alexedwards/scs/v2"
        "github.com/go-playground/form/v4"
        "backgo/internal/models/mocks"
        "time"
    )
    
type testServer struct{
    *httptest.Server
}

func newTestApplication(t *testing.T) *application{
    
    templateCache, err := newTemplateCache()
    if err != nil{
        t.Fatal(err)
    }
    
    formDecoder := form.NewDecoder()
    
    sessionManager := scs.New()
    sessionManager.Lifetime = 12 * time.Hour
    sessionManager.Cookie.Secure = true
    
    return &application{
        errLog: log.New(io.Discard, "", 0),
        infoLog: log.New(io.Discard, "", 0),
        dbconn: &mocks.RecRow{},
        userconn: &mocks.User{},
        formDecoder: formDecoder,
        templateCache: templateCache,
        sessionManager: sessionManager,
    }
}

func newTestServer(t *testing.T,h http.Handler) *testServer{
    ts := httptest.NewTLSServer(h)
    
    jar, err := cookiejar.New(nil)
    if err != nil{
        t.Fatal(err)
    }
    
    ts.Client().Jar = jar
    
    ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }
    return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string){
    rs, err := ts.Client().Get(ts.URL + urlPath)
    if err != nil{
        t.Fatal(err)
    }
    
    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)
    if err != nil{
        t.Fatal(err)
    }
    bytes.TrimSpace(body)
    
    return rs.StatusCode, rs.Header, string(body)
}
