package main

import (
        "net/http"
        "html/template"
        "log"
        "flag"
        "os"
        
        "backgo/internal/models"
        
        _ "github.com/go-sql-driver/mysql"
        
    )
    
type application struct{
    errLog *log.Logger
    infoLog *log.Logger
    dbconn *models.Conn
    templateCache map[string]*template.Template
}

func main(){
    add := flag.String("code", ":4000", "HTTP network address")

    dsn := flag.String("dsn", "backgo:ppp@/dropbox?parseTime=true", "MySQL data source name")

    flag.Parse()
    
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile|log.LUTC)
    
    db, err := openDb(*dsn)
    if err != nil{
        errLog.Fatal(err)
    }
    
    defer db.Close()
    
    templateCache, err := newTemplateCache()
    if err != nil{
        errLog.Fatal(err)
    }
    
    app := &application{
        errLog: errLog,
        infoLog: infoLog,
        dbconn : &models.Conn{DB: db},
        templateCache : templateCache,
    }

    srv := &http.Server{
        Addr: *add,
        ErrorLog: errLog,
        Handler: app.routes(),
    }
   
    infoLog.Printf("Running on %s",*add)
    errr := srv.ListenAndServe()
    errLog.Fatal(errr)
    
}
