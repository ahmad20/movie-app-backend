package entities

type Ticket struct {
	ID      string
	UserID  string
	MovieID int
	Seats   []Seat
	Cost    int
}
