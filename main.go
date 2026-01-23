package main

import (
	"log"

	ticketspb "github.com/Mesailor/study.go-microservices.server-gateway/proto/tickets"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:3082", opts...)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := ticketspb.NewTicketServiceClient(conn)

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
		tickets, err := client.ListTickets(c.Context(), &ticketspb.ListTicketsRequest{})
		if err != nil {
			log.Printf("Error fetching tickets: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tickets"})
		}
		return c.JSON(tickets)
	})

	ticketsGroup.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON("Not implemented yet")
	})

	ticketsGroup.Post("/", func(c *fiber.Ctx) error {
		createTicketData := new(struct {
			Subject     string `json:"subject"`
			Description string `json:"description"`
		})
		if err := c.BodyParser(createTicketData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		newTicket, err := client.CreateTicket(c.Context(), &ticketspb.CreateTicketRequest{
			Subject:     createTicketData.Subject,
			Description: createTicketData.Description,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create ticket"})
		}
		return c.JSON(newTicket)
	})

	ticketsGroup.Put("/:id", func(c *fiber.Ctx) error {
		return c.JSON("Not implemented yet")
	})

	ticketsGroup.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON("Not implemented yet")
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
