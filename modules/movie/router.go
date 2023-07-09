package movie

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SetupRouter(r *gin.Engine, h HandlerInterface) {
	validate := validator.New()
	validate.RegisterValidation("blacklist", BlacklistValidation)

	MovieRouter := r.Group("/")
	MovieRouter.GET("movies", h.GetMovies)
	MovieRouter.GET("movie/details/:id", h.GetMovieDetails)
}
