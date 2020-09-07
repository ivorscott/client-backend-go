package user

import (
	"time"
)

// Users represent user data from the database
type User struct {
	ID        string    `db:"id" json:"id" `
	Auth0ID   string    `db:"auth0Id" json:"auth0Id" `
	Email     string    `db:"email" json:"email"`
	FirstName *string    `db:"firstName" json:"firstName"`
	LastName  *string    `db:"lastName" json:"lastName"`
	Picture   *string    `db:"picture" json:"picture"`
	Created   time.Time `db:"created" json:"created"`
}

type NewUser struct {
	Auth0ID   string  `json:"auth0Id" `
	Email     string  `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Picture   *string `json:"picture"`
}

type UpdateUser struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Picture   *string `json:"picture"`
}
