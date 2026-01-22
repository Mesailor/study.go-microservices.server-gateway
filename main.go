package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	auth := app.Group("/auth")

	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Login successful", "token": "mock_token"})
	})

	auth.Post("/register", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Registration successful", "userId": "mock_user_id"})
	})

	ticketsGroup := app.Group("/api/tickets", authMiddleware)

	ticketsGroup.Get("/", func(c *fiber.Ctx) error {
		return c.JSON([]fiber.Map{
			{"id": "1", "title": "Ticket 1", "description": "Description for ticket 1"},
			{"id": "2", "title": "Ticket 2", "description": "Description for ticket 2"},
		})
	})

	ticketsGroup.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{"id": id, "title": "Mock Ticket", "description": "Mock description for ticket " + id})
	})

	ticketsGroup.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Ticket created successfully", "ticketId": "mock_ticket_id"})
	})

	ticketsGroup.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.JSON(fiber.Map{"message": "Ticket updated successfully", "ticketId": id})
	})

	app.Listen(":3081")
}

func authMiddleware(c *fiber.Ctx) error {
	// Mock authentication check
	token := c.Get("Authorization")
	if token != "Bearer mock_token" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.Next()
}