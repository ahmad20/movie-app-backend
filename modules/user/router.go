package user

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, h HandlerInterface, middleware gin.HandlerFunc) {
	UserRouter := r.Group("/user")
	UserRouter.POST("/register", h.Register)
	UserRouter.POST("/login", h.Login)

	UserRouter.POST("/buy-ticket/:movie_id", middleware, h.BuyTicket)
	UserRouter.POST("/cancel-ticket", middleware, h.CancelTicket)

	UserRouter.GET("/me", middleware, h.GetUser)
	UserRouter.POST("/top-up", middleware, h.TopUp)
	UserRouter.POST("/withdraw", middleware, h.Withdraw)
}
