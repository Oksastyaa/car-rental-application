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
	UserController        controller.UserController
	CarController         controller.CarController
	RentalController      controller.RentalController
	TransactionController controller.TransactionController
}

func RegisterRoutes(e *echo.Echo, c *Controller, jwtSecret []byte) {
	// user routes
	r := e.Group("/api/v1/users")
	r.POST("/register", c.UserController.RegisterUser)
	r.POST("/login", c.UserController.LoginUser)
	r.POST("/deposit", c.UserController.TopUpBalance, middleware.JWTMiddleware(jwtSecret))

	// car routes
	r = e.Group("/api/v1/users/cars", middleware.JWTMiddleware(jwtSecret))
	r.GET("", c.CarController.GetAllCars)
	r.GET("/:id", c.CarController.GetCarByID)

	// admin routes
	r = e.Group("/api/v1/admin", middleware.RoleMiddleware("admin", jwtSecret))
	//car routes
	r.POST("/cars", c.CarController.CreateCar)
	r.PUT("/cars/:id", c.CarController.UpdateCar)
	r.DELETE("/cars/:id", c.CarController.DeleteCar)
	//transaction routes
	r.GET("/transactions", c.TransactionController.GetAllTransactions)
	r.GET("/transactions/:id", c.TransactionController.GetTransactionByID)

	// rental routes
	r = e.Group("/api/v1/users/rentals", middleware.JWTMiddleware(jwtSecret))
	r.POST("", c.RentalController.BookCar)
	r.GET("/:id", c.RentalController.GetRentalByID)
	r.GET("/history", c.RentalController.GetAllRentals)

	//webhook routes
	r = e.Group("/api/v1/webhook")
	r.POST("/xendit", c.TransactionController.UpdateTransactionStatus)
}

func InitController(db *gorm.DB, jwtSecret string) *Controller {
	return &Controller{
		UserController:        controller.NewUserController(service.NewUserService(repository.NewUserRepo(db), jwtSecret)),
		CarController:         controller.NewCarController(service.NewCarService(repository.NewCarRepo(db))),
		RentalController:      controller.NewRentalController(service.NewRentalService(repository.NewRentalRepository(db), repository.NewTransactionRepository(db), repository.NewUserRepo(db))),
		TransactionController: controller.NewTransactionController(service.NewTransactionService(repository.NewTransactionRepository(db))),
	}
}
