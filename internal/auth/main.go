package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string ) (string , error){
    if len([]byte(pass)) > 20 {
        return "",errors.New("password cant exceed 20 bytes")
    }
    if len([]byte(pass)) < 10{
        return "",errors.New("password cant be less than 8 bytes")
    }
    byteHash ,err:=bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    if err != nil{
        return "",err
    }
    return string(byteHash), nil
}



func CheckHashedPassword(password string ,hashedPassword string ) error{
    err:=bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
    if err != nil{
        return err
    }
    return nil
}

