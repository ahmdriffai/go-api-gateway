package main

import (
	"github.com/gofiber/fiber/v2"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Belajar Go", Done: true},
	{ID: 2, Title: "Belajar Go", Done: false},
}

func getTask(c *fiber.Ctx) error {
	return c.JSON(tasks)
}

func createTask(c *fiber.Ctx) error {
	return c.SendString("create task")
}

func main() {
	app := fiber.New()

	tasks := app.Group("/tasks")

	tasks.Get("/", getTask)
	tasks.Get("/create", createTask)

	app.Listen(":3000")
}
