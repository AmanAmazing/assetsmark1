package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth"
)

var DB *sql.DB
var TokenAuth *jwtauth.JWTAuth

type User struct {
    Id int `json:"id"`
	FirstName string `json:"firstName"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
    CreatedAt time.Time `json:"createdAt"` 
    EditedAt  time.Time `json:"editedAt"`
}

type Organisation struct{
    Id int `json:"orgId"`
    Name string `json:"orgName"`
}

type UserOrg struct {
    Id int `json:"id"` 
    UserId int `json:"userId"`
    OrgId int `json:"orgId"`
}

func CreateTables() error{
    ctx, cancelfunc := context.WithTimeout(context.Background(),5*time.Second)
    defer cancelfunc() 
    
    // create EXTENSION citext;
    usersQuery := `create table if not exists users(
	userId SERIAL primary key, 
	firstName varchar(255) not null, 
	lastName varchar(255) not null, 
	hash text not null, 
	email citext unique not null, 
	createdAt timestamp not null default now(), 
	edited_at timestamp )` 
    
    _,err := DB.ExecContext(ctx,usersQuery)
    if err!= nil {
        return fmt.Errorf("CreateTables -> usersQuery error:  %w", err)
    }


    orgQuery := `create table if not exists organisations(
	orgId serial, 
	orgName varchar(255) not null, 
	primary key (orgId) )`
    _,err = DB.ExecContext(ctx,orgQuery)
    if err!= nil {
        return fmt.Errorf("CreateTables -> orgQuery error:  %w", err)
    }
    
    userOrgQuery := `create table if not exists userOrgs (
	userOrgId serial, 
	userId int not null, 
	orgId int not null, 
	foreign key (userId) references users(userId), 
	foreign key (orgId) references organisations(orgId),
    unique (userId, orgId),
	primary key (userOrgId))`
    _,err = DB.ExecContext(ctx,userOrgQuery)
    if err!= nil {
        return fmt.Errorf("CreateTables -> userOrgQuery error:  %w", err)
    }

    return nil 
}


