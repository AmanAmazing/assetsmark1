package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AmanAmazing/assetsMark1/handlers"
	"github.com/AmanAmazing/assetsMark1/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
) 


var workingDirectory string


func init(){
    err := godotenv.Load(".env")
    if err!= nil {
        log.Fatalf("loading env error!; %s", err)
    }
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",os.Getenv("host"),os.Getenv("port"),os.Getenv("user"),os.Getenv("password"),os.Getenv("dbname"))

    models.DB, err = sql.Open("postgres", psqlInfo)
    if err!= nil {
        panic(err)
    }

    err = models.DB.Ping()
    if err != nil {
        panic(err)
    }
    // checking if all the user tables exist 
    err = models.CreateTables()
    if err != nil {
        log.Fatal(err)
    }
    // working directory 
    workingDirectory, err = os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }
    
    // jwt stuff 
    models.TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("jwtSecret")),nil)

    log.Println("Connected")
}


func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/",handlers.HelloWorld)
    r.Get("/login",handlers.Login)
    r.Get("/signup",handlers.SignUp)
    r.Post("/signup",handlers.SignUpPost)
    
    r.Post("/login",handlers.LoginPost)

    // protected routes 
    r.Group(func(r chi.Router){
        // seek, verify and validate JWT tokens 
        r.Use(jwtauth.Verifier(models.TokenAuth))
        r.Use(jwtauth.Authenticator)

        r.Get("/admin", func(w http.ResponseWriter, r *http.Request){
            _, claims, _ := jwtauth.FromContext(r.Context())
            w.Write([]byte(fmt.Sprint("protected area. Hi ",claims["id"])))
        })
        
        r.Get("/organisations",handlers.Organisations)
        r.Post("/OrganisationAdd",handlers.OrganisationAdd)
    })
    log.Fatal(http.ListenAndServe(":3000", r))
}
