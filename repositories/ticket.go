package repositories

import (
	"errors"
	"movie-app-go/entities"
)

type TicketRepository struct {
	data []entities.Ticket
}
type TicketRepositoryInterface interface {
	Create(ticket entities.Ticket) error
	Read(id string) (entities.Ticket, error)
	Update(ticket entities.Ticket) error
	Delete(ticket entities.Ticket) error
}

func NewTicketRepository(data []entities.Ticket) TicketRepositoryInterface {
	return &TicketRepository{
		data: data,
	}
}

func (repo *TicketRepository) Create(ticket entities.Ticket) error {
	for _, existingTicket := range repo.data {
		if existingTicket.ID == ticket.ID {
			return errors.New("ticket with the same ID already exists")
		}
	}
	repo.data = append(repo.data, ticket)
	return nil
}

func (repo *TicketRepository) Read(id string) (entities.Ticket, error) {
	for _, existingTicket := range repo.data {
		if existingTicket.ID == id {
			return existingTicket, nil
		}
	}
	return entities.Ticket{}, errors.New("NOT_FOUND")
}

func (repo *TicketRepository) Update(ticket entities.Ticket) error {
	for i, existingTicket := range repo.data {
		if existingTicket.ID == ticket.ID {
			repo.data[i] = ticket
		}
	}
	return errors.New("NOT_FOUND")
}

func (repo *TicketRepository) Delete(ticket entities.Ticket) error {
	for i, existingTicket := range repo.data {
		if existingTicket.ID == ticket.ID {
			repo.data = append(repo.data[:i], repo.data[i+1:]...)
			return nil
		}
	}
	return errors.New("NOT_FOUND")
}
