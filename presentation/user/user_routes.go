package user

import (
	"GoSyntaxDoc/infrastructure/database"
	"GoSyntaxDoc/infrastructure/repositories"
	"GoSyntaxDoc/services/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	userRepo := repositories.NewUserRepository(&database.Database)
	userService := user.NewUserService(userRepo)
	app.Get("/users/:id", userService.FetchUserByIdHandler)
	app.Post("/users/create", userService.CreateUserHandler)
}
