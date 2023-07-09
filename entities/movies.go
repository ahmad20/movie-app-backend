package entities

import "time"

type Movie struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Release_date string    `json:"release_date"`
	Age_rating   int       `json:"age_rating"`
	Poster_url   string    `json:"poster_url"`
	Ticket_price int       `json:"ticket_price"`
	Seats        []Seat    `json:"seats"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}
