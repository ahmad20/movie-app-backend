package user

import (
	"errors"
	"log"
	"movie-app-go/entities"
	"movie-app-go/repositories"
)

type useCase struct {
	userRepo   repositories.UserRepositoryInterface
	movieRepo  repositories.MovieRepositoryInterface
	ticketRepo repositories.TicketRepositoryInterface
}

type UseCaseInterface interface {
	Create(user entities.User) error
	GetMovie(id int) (entities.Movie, error)
	GetTicket(id string) (entities.Ticket, error)
	GetUser(username string) (entities.User, error)
	BuyTicket(u entities.User, t entities.Ticket) error
	CancelTicket(u entities.User, t entities.Ticket) error
	NotPermitted(m entities.Movie, u entities.User) bool
	CheckBalance(u entities.User, p int) error
	CheckAvailability(seat []entities.Seat, m *entities.Movie) ([]entities.Seat, error)
	TopUp(u *entities.User, n int) error
	Withdraw(u *entities.User, n int) error
}

func NewUseCase(userRepo repositories.UserRepositoryInterface,
	movieRepo repositories.MovieRepositoryInterface,
	ticketRepo repositories.TicketRepositoryInterface) UseCaseInterface {
	return &useCase{
		userRepo:   userRepo,
		movieRepo:  movieRepo,
		ticketRepo: ticketRepo,
	}
}
func (usecase *useCase) Create(user entities.User) error {
	if err := usecase.userRepo.Create(user); err != nil {
		return err
	}
	return nil
}

func (usecase *useCase) GetMovie(id int) (entities.Movie, error) {
	movie, err := usecase.movieRepo.Read(id)
	if err != nil {
		return entities.Movie{}, err
	}
	return movie, nil
}
func (usecase *useCase) GetTicket(id string) (entities.Ticket, error) {
	ticket, err := usecase.ticketRepo.Read(id)

	return ticket, err
}

func (usecase *useCase) GetUser(username string) (entities.User, error) {
	user, err := usecase.userRepo.GetByUsername(username)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (usecase *useCase) BuyTicket(u entities.User, t entities.Ticket) error {
	if err := usecase.ticketRepo.Update(t); err != nil {
		if err := usecase.ticketRepo.Create(t); err != nil {
			return err
		}
	}
	u.Ticket = append(u.Ticket, t)
	if err := usecase.userRepo.Update(u); err != nil {
		return err
	}
	return nil
}

func (usecase *useCase) NotPermitted(m entities.Movie, u entities.User) bool {
	return m.Age_rating > u.Age
}

func (usecase *useCase) CheckBalance(u entities.User, p int) error {
	if u.Balance < p {
		return errors.New("BALANCE_INSUFFICIENT")
	}
	return nil
}
func (usecase *useCase) CheckAvailability(seat []entities.Seat, m *entities.Movie) ([]entities.Seat, error) {
	var (
		// countTicket int
		seats    []entities.Seat
		err      error
		prevSeat *entities.Seat
	)
	if len(seat) > 6 {
		return nil, errors.New("MAX_TICKET_REACH")
	}

	for _, s := range seat {
		var foundSeat *entities.Seat

		for i := range m.Seats {
			if s.Row == m.Seats[i].Row && s.Number == m.Seats[i].Number {
				foundSeat = &m.Seats[i]
				break
			}
		}

		if foundSeat == nil || foundSeat.Booked {
			err = errors.New("PANIC")
			// countTicket = 0
			seats = nil
			if prevSeat != nil {
				prevSeat.Booked = false
			}
			return seats, err
		}

		// foundSeat.Booked = true
		prevSeat = foundSeat
		// countTicket++
		seats = append(seats, *foundSeat)
		log.Printf("Booked seat: %v", foundSeat)
	}
	if len(seat) > len(seats) || seats == nil {
		err = errors.New("TICKET_UNAVAILABLE")
	}
	return seats, err
}

func (usecase *useCase) TopUp(u *entities.User, n int) error {
	if n < 0 {
		return errors.New("NEGATIVE_VALUE")
	}
	u.Balance += n
	if err := usecase.userRepo.Update(*u); err != nil {
		return err
	}
	return nil
}

func (usecase *useCase) Withdraw(u *entities.User, n int) error {
	if n < 0 {
		return errors.New("NEGATIVE_VALUE")
	}
	if n > u.Balance {
		return errors.New("BALANCE_INSUFFICIENT")
	}
	u.Balance -= n
	if err := usecase.userRepo.Update(*u); err != nil {
		return err
	}
	return nil
}

func (usecase *useCase) CancelTicket(u entities.User, t entities.Ticket) error {
	for i, t2 := range u.Ticket {
		if t.ID == t2.ID {
			u.Ticket = append(u.Ticket[:i], u.Ticket[i+1:]...)
			break
		}
	}

	if err := usecase.userRepo.Update(u); err != nil {
		return err
	}
	if err := usecase.ticketRepo.Delete(t); err != nil {
		return err
	}
	return nil
}
