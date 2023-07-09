package user

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"movie-app-go/entities"
	"movie-app-go/modules/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	userUseCase UseCaseInterface
	auth        auth.AuthInterface
}

type HandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	BuyTicket(c *gin.Context)
	CancelTicket(c *gin.Context)
	GetUser(c *gin.Context)
	TopUp(c *gin.Context)
	Withdraw(c *gin.Context)
}

func NewHandler(userUseCase UseCaseInterface, auth auth.AuthInterface) HandlerInterface {
	return &handler{
		userUseCase: userUseCase,
		auth:        auth,
	}
}

func (h handler) Register(c *gin.Context) {
	var req Register
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Handler.Register.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "INVALID_REQUEST",
			Data:    err.Error(),
		})
		return
	}
	UUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Handler.Register.02", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_GENERATE_UUID",
			Data:    err.Error(),
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Handler.Register.03", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_HASH_PASSWORD",
			Data:    err.Error(),
		})
		return
	}
	newUser := entities.User{
		ID:         UUID.String(),
		Username:   req.Username,
		Password:   string(hashedPassword),
		Name:       req.Name,
		Age:        req.Age,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	if err := h.userUseCase.Create(newUser); err != nil {
		log.Println("Handler.Register.04", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_USECASE",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusCreated,
		Message: "CREATED_USER",
		Data:    UUID.String(),
	})
}

func (h handler) GetUser(c *gin.Context) {
	authInfo, _ := c.Get("AuthInfo")
	user, err := h.userUseCase.GetUser(authInfo.(auth.AuthInfo).Username)
	if err != nil {
		log.Println("Handler.GetUser.01", err.Error())

		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    err.Error(),
		})
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    user,
	})
}
func (h handler) Login(c *gin.Context) {
	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Handler.Login.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "INVALID_REQUEST",
			Data:    err.Error(),
		})
		return
	}

	user, err := h.userUseCase.GetUser(req.Username)
	if err != nil {
		log.Println("Handler.Login.02", err.Error())

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "USERNAME_NOT_FOUND",
			Data:    err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("Handler.Login.03", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "PASSWORD_NOT_MATCH",
			Data:    err.Error(),
		})
		return
	}

	token, err := h.auth.GenerateToken(user.Username)
	if err != nil {
		log.Println("Handler.Login.04", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_GENERATE_TOKEN",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusCreated,
		Message: "CREATED_USER",
		Data:    token,
	})
}

func (h handler) BuyTicket(c *gin.Context) {
	var req BuyTicketRequest
	movieId := c.Param("movie_id")
	if movieId == "" {
		log.Println("Handler.BuyTicket.01", errors.New("BAD_REQUEST"))

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "INVALID_REQUEST",
			Data:    errors.New("BAD_REQUEST"),
		})
		return
	}
	mid, _ := strconv.Atoi(movieId)
	authInfo, _ := c.Get("AuthInfo")
	user, err := h.userUseCase.GetUser(authInfo.(auth.AuthInfo).Username)
	if err != nil {
		log.Println("Handler.BuyTicket.02", err.Error())

		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    err.Error(),
		})
		return
	}

	movie, err := h.userUseCase.GetMovie(mid)
	if err != nil {
		log.Println("Handler.BuyTicket.03", err.Error())

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "NOT_FOUND",
			Data:    err.Error(),
		})
		return
	}
	// Check age
	if h.userUseCase.NotPermitted(movie, user) {
		log.Println("Handler.BuyTicket.04", errors.New("AGE_RESTRICTION"))

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "AGE_RESTRICTION",
			Data:    errors.New("AGE_RESTRICTION"),
		})
		return
	}
	// Select Seats
	if err := c.ShouldBindJSON(&req); err != nil {
		if err != nil {
			log.Println("Handler.BuyTicket.05", err.Error())

			c.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Message: "FAILED_BINDING_REQUEST",
				Data:    err.Error(),
			})
			return
		}
	}
	// Check available seats
	seats, err := h.userUseCase.CheckAvailability(req.Seats, &movie)
	if err != nil {
		log.Println("Handler.BuyTicket.06", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_BOOKED",
			Data:    err.Error(),
		})
		return
	}

	// Check Balance
	costs := movie.Ticket_price * len(seats)
	if err := h.userUseCase.CheckBalance(user, costs); err != nil {
		log.Println("Handler.BuyTicket.07", err)

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "NOT_ENOUGH_BALANCE",
			Data:    err,
		})
		return
	}
	if err := h.userUseCase.Withdraw(&user, costs); err != nil {
		log.Println("Handler.BuyTicket.08", err)

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "NOT_ENOUGH_BALANCE",
			Data:    err,
		})
		return
	}
	for i1, s1 := range seats {
		for i2, s2 := range movie.Seats {
			if s1 == s2 {
				seats[i1].Booked = true
				movie.Seats[i2].Booked = true
			}
		}
	}

	UUID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Handler.BuyTicket.09", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_GENERATE_UUID",
			Data:    err.Error(),
		})
		return
	}

	newTicket := entities.Ticket{
		ID:      UUID.String(),
		MovieID: movie.ID,
		UserID:  user.ID,
		Seats:   seats,
		Cost:    costs,
	}
	if err := h.userUseCase.BuyTicket(user, newTicket); err != nil {
		log.Println("Handler.BuyTicket.10", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_CREATE_TICKET",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "TICKET_CREATED",
		Data:    newTicket,
	})

}

func (h handler) TopUp(c *gin.Context) {
	var req Balance
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Handler.TopUp.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "BAD_REQUEST",
			Data:    err.Error(),
		})
		return
	}
	authInfo, _ := c.Get("AuthInfo")
	user, err := h.userUseCase.GetUser(authInfo.(auth.AuthInfo).Username)
	if err != nil {
		log.Println("Handler.TopUp.02", err.Error())

		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    err.Error(),
		})
		return
	}
	if err := h.userUseCase.TopUp(&user, req.Amount); err != nil {
		log.Println("Handler.TopUp.03", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_TOP_UP",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS_TOP_UP",
		Data:    user.ID,
	})
}

func (h handler) Withdraw(c *gin.Context) {
	var req Balance
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Handler.Withdraw.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "BAD_REQUEST",
			Data:    err.Error(),
		})
		return
	}
	authInfo, _ := c.Get("AuthInfo")
	user, err := h.userUseCase.GetUser(authInfo.(auth.AuthInfo).Username)
	if err != nil {
		log.Println("Handler.Withdraw.02", err.Error())

		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    err.Error(),
		})
		return
	}
	if err := h.userUseCase.Withdraw(&user, req.Amount); err != nil {
		log.Println("Handler.Withdraw.03", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_WITHDRAW",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS_WITHDRAW",
		Data:    user.ID,
	})
}
func (h handler) CancelTicket(c *gin.Context) {
	authInfo, _ := c.Get("AuthInfo")
	user, err := h.userUseCase.GetUser(authInfo.(auth.AuthInfo).Username)
	if err != nil {
		log.Println("Handler.CancelTicket.01", err.Error())

		c.JSON(http.StatusUnauthorized, Response{
			Code:    http.StatusUnauthorized,
			Message: "UNAUTHORIZED",
			Data:    err.Error(),
		})
		return
	}
	var req TicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Handler.CancelTicket.02", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "BAD_REQUEST",
			Data:    err.Error(),
		})
		return
	}
	ticket, err := h.userUseCase.GetTicket(req.ID)
	if err != nil {
		log.Println("Handler.CancelTicket.03", err.Error())

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "NOT_FOUND",
			Data:    err.Error(),
		})
		return
	}
	seats := ticket.Seats
	movie, _ := h.userUseCase.GetMovie(ticket.MovieID)
	for i1, s1 := range seats {
		for i2, s2 := range movie.Seats {
			if s1 == s2 {
				seats[i1].Booked = false
				movie.Seats[i2].Booked = false
			}
		}
	}
	if err := h.userUseCase.CancelTicket(user, ticket); err != nil {
		log.Println("Handler.CancelTicket.03", err.Error())

		c.JSON(http.StatusNotFound, Response{
			Code:    http.StatusNotFound,
			Message: "NOT_FOUND",
			Data:    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    nil,
	})
}
