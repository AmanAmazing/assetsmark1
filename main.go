package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
) 

type SignUser struct{
    FirstName string `json:"firstName"`
    Surname string `json:"surname"`
    Email string `json:"email"`
    Password string `json:"password"`
}

var Db *sql.DB
var workingDirectory string

func init(){
    err := godotenv.Load(".env")
    if err!= nil {
        log.Fatalf("loading env error!; %s", err)
    }
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",os.Getenv("host"),os.Getenv("port"),os.Getenv("user"),os.Getenv("password"),os.Getenv("dbname"))

    Db, err = sql.Open("postgres", psqlInfo)
    if err!= nil {
        panic(err)
    }

    err = Db.Ping()
    if err != nil {
        panic(err)
    }
    // working directory 
    workingDirectory, err = os.Getwd()
    if err != nil {
        log.Fatalln(err)
    }


    log.Println("Connected")
}


func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Get("/",helloWorld)
    r.Get("/login",login)
    r.Get("/signup",signUp)
    r.Post("/signup",signUpPost)

    log.Fatal(http.ListenAndServe(":3000", r))
}

func helloWorld(w http.ResponseWriter,r *http.Request){
    w.Write([]byte("hello world"))
}

func login(w http.ResponseWriter,r *http.Request){
    tmpl := template.Must(template.ParseFiles("login.html"))
    tmpl.Execute(w,nil)
}

func signUp(w http.ResponseWriter,r *http.Request){
    tmpl := template.Must(template.ParseFiles("signup.html"))
    tmpl.Execute(w, nil)
}

func hashPassword(password string) (string, error){
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),14)
    return string(hashedPassword), err
}
func checkPasswordsMatch(hashedPassword, enteredPassword string) bool{
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(enteredPassword))
    return err == nil 
}

func signUpPost(w http.ResponseWriter,r *http.Request){
    user := SignUser{} 
    sqlQuery := `SELECT email FROM users WHERE email=$1;`
    fmt.Println("checking email: ",r.FormValue("email"))
    row := Db.QueryRow(sqlQuery,r.FormValue("email"))
    switch err := row.Scan(&user.Email);err{
    case sql.ErrNoRows:
        user.Email = r.FormValue("email")
        user.FirstName = r.FormValue("firstName")
        user.Surname = r.FormValue("surname")
        user.Password, err = hashPassword(r.FormValue("password")) 
        if err != nil {
            log.Fatal("failed to hash password") // need to complete this section
            return
        }
        fmt.Println(user)
        signUpQuery := `INSERT INTO users (first_name,last_name,hash,email) values ($1,$2,$3,$4)` 
        _, err = Db.Exec(signUpQuery,user.FirstName,user.Surname,user.Password,user.Email)
        if err !=nil{
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Failed to post user to db"))
            return
        }
        w.WriteHeader(http.StatusCreated) 
        w.Write([]byte("User was created successfully"))
        return
    case nil:
        w.Write([]byte("User already has an account"))
    default:
        panic(err)
    }
    
}
