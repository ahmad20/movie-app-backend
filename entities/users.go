package entities

import "time"

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Name       string    `json:"name"`
	Age        int       `json:"age"`
	Balance    int       `json:"balance"`
	Ticket     []Ticket  `json:"ticket"`
	Created_at time.Time `json:"-"`
	Updated_at time.Time `json:"-"`
}
