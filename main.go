package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct { // Define a struct
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello, World!")

	app := fiber.New() // Create a new Fiber instance

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	// Get all todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Create a new todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // Create a new instance of our struct with default values {id: 0, completed: false, body: ""}

		// Check if the body is valid JSON, then parse it into the `todo` struct
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1     // Generate a new ID
		todos = append(todos, *todo) // Append the new todo to the todos slice

		return c.Status(201).JSON(todo)

	})

	// Update a todo (Use put or patch)
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") // Get the id from the URL

		for i, todo := range todos { // Loop through all todos, and find the one with the matching ID, then update it
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"}) // If no todo with the ID was found, return a 404
	})

	// Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") // Get the id from the URL

		for i, todo := range todos { // Loop through all todos, and find the one with the matching ID, then remove it
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...) // Remove the todo from the slice
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"}) // If no todo with the ID was found, return a 404
	})

	log.Fatal(app.Listen(":" + PORT))
}
