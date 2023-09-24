package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/AmanAmazing/assetsMark1/helper"
	"github.com/AmanAmazing/assetsMark1/models"
)

func HelloWorld(w http.ResponseWriter,r *http.Request){
    w.Write([]byte("hello world"))
}

func SignUp(w http.ResponseWriter,r *http.Request){
    tmpl := template.Must(template.ParseFiles("assets/signup.html"))
    tmpl.Execute(w, nil)
}



func SignUpPost(w http.ResponseWriter,r *http.Request){
    user := models.User{} 
    sqlQuery := `SELECT email FROM users WHERE email=$1;`
    fmt.Println("checking email: ",r.FormValue("email"))
    row := models.DB.QueryRow(sqlQuery,r.FormValue("email"))
    switch err := row.Scan(&user.Email);err{
    case sql.ErrNoRows:
        user.Email = r.FormValue("email")
        user.FirstName = r.FormValue("firstName")
        user.Surname = r.FormValue("surname")
        user.Password, err = helper.HashPassword(r.FormValue("password")) 
        if err != nil {
            log.Fatal("failed to hash password") // need to complete this section
            return
        }
        fmt.Println(user)
        signUpQuery := `INSERT INTO users (first_name,last_name,hash,email) values ($1,$2,$3,$4)` 
        _, err = models.DB.Exec(signUpQuery,user.FirstName,user.Surname,user.Password,user.Email)
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

func Login(w http.ResponseWriter,r *http.Request){
    tmpl := template.Must(template.ParseFiles("assets/login.html"))
    tmpl.Execute(w,nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request){
    var credentials models.User
    credentials.Email = r.FormValue("email")
    credentials.Password = r.FormValue("password")
    var hashFromDatabase string 
    sqlQuery := `SELECT hash FROM users WHERE email=$1;`
    row := models.DB.QueryRow(sqlQuery,credentials.Email)
    switch err := row.Scan(&hashFromDatabase);err{
    case sql.ErrNoRows:
        w.WriteHeader(http.StatusNotFound)
        return 
    } 
    fmt.Println("Entered hashing password stage")
    fmt.Println(hashFromDatabase)
    if helper.CheckPasswordsMatch(hashFromDatabase, credentials.Password) == false{
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    _,tokenString, _ := models.TokenAuth.Encode(map[string]interface{}{"email":credentials.Email})
    
    http.SetCookie(w, &http.Cookie{
        HttpOnly: true,
        Expires: time.Now().Add(7*24*time.Hour),
        //uncomment below for https: 
        // Secure:true, 
        Name: "jwt",
        Value: tokenString,
    })

    http.Redirect(w,r, "/admin",http.StatusSeeOther)

}
