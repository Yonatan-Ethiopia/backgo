package models

import(
    "fmt"
    "time"
    "database/sql"
    "errors"
    "strings"
    "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

type User struct{
    Id              int
    Name            string
    Email           string
    HashedPassword  []byte
    Created         time.Time
}

type ResponseUser struct{
    Name    string
    Email   string
    Created time.Time
}

type UserInterface interface{
    Insert(name, email, password string) (error)
    Authenticate(email, password string) (int, error)
    Exists(id int) (bool, error)
    Get(id int) (*ResponseUser, error)
    Match(id int, password string) (bool, error)
    UpdatePass(id int, new_pass string) (error)
}

type UserModel struct{
    DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil{
        return err
    }
    stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES(?, ?, ?, UTC_TIMESTAMP())`
    _, err = u.DB.Exec(stmt, name, email, hashedPassword)
    if err != nil{
        
        var mySQLError *mysql.MySQLError
        if errors.As(err, &mySQLError){
            if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email"){
                return ErrDuplicateEmail
            }
        }
        return err
    }
    
    return nil
    
}
func (u *UserModel) Authenticate(email, password string) (int, error){
    
    var id int
    var hashedPassword []byte
    
    stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
    
    err := u.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }
    
    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err != nil{
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return 0, ErrInvalidCredentials
        } else {
            return 0, err
        }
    }
    
    return id, nil
    
}
func (u *UserModel) Exists(id int) (bool, error){
    
    var exists bool
    
    stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
    
    err := u.DB.QueryRow(stmt, id).Scan(&exists)
    
    return exists, err
    
}

func (u *UserModel) Get(id int) (*ResponseUser, error){
    responseUser := &ResponseUser{}
    
    stmt := "SELECT name, email, created FROM users WHERE id = ?"
    
    err := u.DB.QueryRow(stmt, id).Scan(&responseUser.Name, &responseUser.Email, &responseUser.Created )
    
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return nil, ErrNoRecord
        }
        return nil, err
    }
    
    return responseUser, nil
}

func (u *UserModel) Match(id int, password string) (bool, error){

    var hashedPassword []byte
    
    stmt := "SELECT hashed_password FROM users WHERE id = ?"
    
    err := u.DB.QueryRow(stmt, id).Scan(&hashedPassword)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            fmt.Printf("Debug no id for %d", id)
            return false, ErrInvalidCredentials
        } else {
            return false, err
        }
    }
    
    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err != nil{
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return false, nil
        } else {
            return false, err
        }
    }
    
    return true, nil
    
    
}

func (u *UserModel) UpdatePass(id int, new_pass string) error{
    
    var userData User
    
    stmt := "SELECT * FROM users WHERE id = ?"
    err := u.DB.QueryRow(stmt, id).Scan(&userData.Id, &userData.Name, &userData.Email, &userData.HashedPassword, &userData.Created)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return ErrInvalidCredentials
        }
        return err
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_pass), 12)
    if err != nil{
        return err
    }
    
    stmt = "UPDATE users SET hashed_password = ? WHERE id = ?"
    
    _,err = u.DB.Exec(stmt, string(hashedPassword), id)
    return err
    
    
}

