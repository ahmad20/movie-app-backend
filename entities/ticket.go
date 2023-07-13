package entities

import "time"

type Ticket struct {
	ID         string
	UserID     string
	Movie      Movie
	Seats      []Seat
	Cost       int
	Created_At time.Time
	Updated_At time.Time
}
