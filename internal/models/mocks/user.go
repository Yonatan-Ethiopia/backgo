package mocks

import (
    "backgo/internal/models"
    )
    
type User struct{}

type ResponseUser struct{}

func (u *User) Insert(name, email, password string) error {
    switch email {
        case "dupe@example.com":
            return models.ErrDuplicateEmail
        default:
            return nil
    }
}

func (u *User) Authenticate(email, password string) (int, error){
    if email == "dupe@exmaple.com" && password == "validPassword"{
         return 1, nil
    }
    
    return 0, models.ErrInvalidCredentials
}

func (u *User) Exists(id int) (bool, error){
    switch id{
        case 1:
            return true, nil
        default: 
            return false, nil
    }
}

func (u *User) Get(id int) (*ResponseUser, error){
    validResponse := &ResponseUser{
        Name: "Dup",
        Email: "dup@gmail.com",
        Created: time.Now(),
    }
    if id == 1{
        return validResponse, nil
    }
    return nil, models.ErrNoRecord
}
