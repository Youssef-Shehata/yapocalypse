package auth

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHashing(t *testing.T) {
    pass := "1234ABCabc_)>"
    hashed, err := HashPassword(pass)
    if err != nil {
        t.Fatalf(`HashPassword(%v) = %v, Error : %v`,pass, hashed, err )
    }
    error := CheckHashedPassword(pass , hashed)
    if err != nil {
        t.Fatalf(`HashPassword(%v) = %v, Error : %v`,pass, hashed, error )
    }
    log.Println("Hasing test PASSED")
}
func TestHashing2(t *testing.T) {
    pass := "(*&^%$#12@!+?:AP{}\"<"

    hashed, err := HashPassword(pass)
    if err != nil {
        t.Fatalf(`HashPassword(%v) = %v, Error : %v`,pass, hashed, err )
    }
    error := CheckHashedPassword(pass , hashed)
    if err != nil {
        t.Fatalf(`HashPassword(%v) = %v, Error : %v`,pass, hashed, error )
    }
    log.Println("Hasing test PASSED")
}

func TestJWT(t *testing.T) {

    id ,err:=uuid.NewRandom()
    if err !=nil {
        log.Fatalf("why the fuck would a random genrator fail and return an error!!!!")
    }

	godotenv.Load()
  token, err := MakeJWT(id ,os.Getenv("SECRET") , time.Second*10 )
    if err != nil {
        t.Fatalf(`MakeJWT(%v) = %v, Error : %v`,id, token, err )
    }
    
    resId ,error := ValidateJWT(token, os.Getenv("SECRET"))
    if err != nil {
        t.Fatalf(`MakeJWT(%v) = %v, Error : %v`,token, resId, error )
    }
    if resId != id {

        t.Fatalf(`MakeJWT() : id (%v)  !=  result (%v)`,token, resId )
    }
    log.Println("JWT test PASSED")
}
