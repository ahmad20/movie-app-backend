package movie

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type handler struct {
	movieUseCase UseCaseInterface
	validate     *validator.Validate
}

type HandlerInterface interface {
	GetMovies(c *gin.Context)
	GetMovieDetails(c *gin.Context)
}

var (
	ErrMovieMissing   = errors.New("movie does not exist")
	ErrInternalServer = errors.New("internal server error")
	ErrEmptyList      = errors.New("list empty")
	ErrBadRequest     = errors.New("bad request")
)

func NewHandler(movieUseCase UseCaseInterface) HandlerInterface {
	validation := validator.New()
	validation.RegisterValidation("blacklist", BlacklistValidation)
	return &handler{
		movieUseCase: movieUseCase,
		validate:     validation,
	}
}

func (h handler) GetMovies(c *gin.Context) {
	var req Request

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Println("Handler.GetMovies.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_BIND_QUERY",
			Data:    err.Error(),
		})
		return
	}
	if err := h.validate.Struct(req); err != nil {
		log.Println("Handler.GetMovies.02", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "VALIDATION_ERROR",
			Data:    err.Error(),
		})
		return
	}

	movies, err := h.movieUseCase.GetAll()
	if err != nil {
		log.Println("Handler.GetMovies.02", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_USECASE",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    movies,
	})
}

func (h handler) GetMovieDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Handler.GetMovieDetails.01", err.Error())

		c.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Message: "FAILED_CONVERT_ID",
			Data:    err.Error(),
		})
		return
	}
	movies, err := h.movieUseCase.GetById(id)
	if err != nil {
		log.Println("Handler.GetMovieDetails.02", err.Error())

		c.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "FAILED_USECASE",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "SUCCESS",
		Data:    movies,
	})
}
