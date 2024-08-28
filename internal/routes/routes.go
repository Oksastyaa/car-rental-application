package routes

import (
	"car-rental-application/internal/controller"
	"car-rental-application/internal/middleware"
	"car-rental-application/internal/repository"
	"car-rental-application/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Controller struct {
	UserController controller.UserController
}

func RegisterRoutes(e *echo.Echo, c *Controller, jwtSecret []byte) {
	// user routes
	r := e.Group("/api/v1/users")
	r.POST("/register", c.UserController.RegisterUser)
	r.POST("/login", c.UserController.LoginUser)
	r.POST("/deposit", c.UserController.TopUpBalance, middleware.JWTMiddleware(jwtSecret))

}

func InitController(db *gorm.DB, jwtSecret string) *Controller {
	return &Controller{
		UserController: controller.NewUserController(service.NewUserService(repository.NewUserRepo(db), jwtSecret)),
	}
}
