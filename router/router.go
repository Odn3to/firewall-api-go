package router

import (
	"github.com/gofiber/fiber/v2"

	"firewall-api-go/controllers"
)

func Register(app *fiber.App) {

	firewall := app.Group("/firewall")

	// Adicione o middleware às rotas
	firewall.Use(controllers.TokenValidationMiddleware)

	firewall.Get("/regras/:filter?", controllers.GetRegras)
	firewall.Post("/newregra", controllers.CreateRegra)
	firewall.Post("/apply", controllers.Apply)
	firewall.Delete("/regra/:id", controllers.DeleteRegra)
	firewall.Put("/regra/:id", controllers.EditRegra)
}