package main

import (
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// User maps the users table into a structure
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName" db:"first_name" form:"fname"`
	LastName  string `json:"lastName" db:"last_name" form:"lname"`
	Email     string `json:"email" form:"email"`
	Token     string `json:"-" db:"token"`
}

// IsValid makes sure all fields are set and look right
func (user *User) IsValid() bool {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return false
	}
	return true
}

// FullName return the full name of the user
func (user *User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

// TableName defines the expected DB table name that holds data for users
func (user *User) TableName() string {
	return "users"
}

// FindByEmail finds a user by email
func (user *User) FindByEmail(db *dbx.DB, email string) error {
	q := db.NewQuery("select * from users where email={:email} limit 1")
	q.Bind(dbx.Params{"email": email})
	return q.One(user)
}

// FindByToken finds a user by access token
func (user *User) FindByToken(db *dbx.DB, token string) error {
	q := db.NewQuery("select * from users where token={:token} limit 1")
	q.Bind(dbx.Params{"token": token})
	return q.One(user)
}
