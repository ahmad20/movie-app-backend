package repositories

import (
	"movie-app-go/entities"
)

func GenerateSeats() []entities.Seat {
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	seats := make([]entities.Seat, 0)

	for _, row := range rows {
		for i := 1; i <= 8; i++ {
			seat := entities.Seat{
				Row:    row,
				Number: i,
				Booked: false,
			}
			seats = append(seats, seat)
		}
	}

	return seats
}
