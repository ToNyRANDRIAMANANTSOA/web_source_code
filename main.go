package main

import (
	"back/database"
	"back/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	// Student routes
	app.Post("/student", routes.CreateStudent)
	app.Get("/student", routes.GetStudents)
	app.Get("/student/:matricule", routes.GetStudent)
	app.Get("/student/niveau/:niveau", routes.ByNiveau)
	app.Delete("/student/:matricule", routes.DeleteStudent)
	app.Put("/student/:matricule", routes.UpdateStudent)

	app.Post("/ec", routes.CreateEC)
	app.Get("/ec", routes.GetECs)
	app.Get("/ec/:codeEC", routes.GetEC)
	app.Put("/ec/:codeEC", routes.UpdateEC)

	app.Post("/note", routes.CreateNote)
	app.Get("/note/student/:matricule", routes.GetNotesByMatricule)
	app.Put("/note/:id", routes.UpdateNote)
	app.Delete("/note/:id", routes.DeleteNote)
}
func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Content-Type",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(logger.New())
	database.Connect()

	setupRoutes(app)

	app.Listen(":3000")
}
