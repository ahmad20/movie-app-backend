package user

import "movie-app-go/entities"

type RequestInterface interface {
	Validate() interface{}
}

type (
	Register struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Age      int    `json:"age" binding:"required,number"`
	}
	Login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	BuyTicketRequest struct {
		Seats []entities.Seat `json:"seats" binding:"required"`
	}
	Balance struct {
		Amount int `json:"amount" binding:"required,number"`
	}
	TicketRequest struct {
		ID string `json:"id" binding:"required"`
	}
	Response struct {
		Code    int    `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
		Data    any    `json:"data" binding:"required"`
	}
)
