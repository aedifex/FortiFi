package database

import "time"

type User struct {
	Id 			string	`json:"id" sql:"id"`
 	FirstName 	string	`json:"first_name" sql:"first_name"`
 	LastName 	string 	`json:"last_name" sql:"last_name"`
 	Email 		string	`json:"email" sql:"email"`
 	Password 	string	`json:"password" sql:"password"`
}

type Token struct {
	Token		string
	FK_UserId 	string
    expires		time.Time
}