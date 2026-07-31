package models

import(
    "time"
    "database/sql"
)

type User struct{
    Id              int
    Name            string
    Email           string
    HashedPassword  []byte
    Created         time.Time
}

type UserModel struct{
    DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error{
    
    return nil
    
}
func (u *UserModel) Authenticate(email, password string) (int, error){
    
    return 0, nil
    
}
func (u *UserModel) Exists(email string) (bool, error){
    
    return true, nil
    
}
func (u *UserModel) Remove(id int) (bool, error){
    
    return true, nil
    
}

