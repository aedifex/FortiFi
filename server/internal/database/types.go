package database

type User struct {
	Id 			string	`json:"id" sql:"id"`
 	FirstName 	string	`json:"first_name" sql:"first_name"`
 	LastName 	string 	`json:"last_name" sql:"last_name"`
 	Email 		string	`json:"email" sql:"email"`
 	Password 	string	`json:"password" sql:"password"`
}

type Pi struct {
	Id			string	`json:"id" sql:"id"`
}