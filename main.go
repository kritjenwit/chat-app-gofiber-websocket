package main

import (
	"fmt"

	"example.com/Chat-app/database"
	"example.com/Chat-app/handlers"
	"example.com/Chat-app/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	_ "github.com/joho/godotenv/autoload"
)

var (
	chatHistorys map[int64]fiber.Map
)

func setup() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use("/ws", handlers.Upgrade)
	return app
}

func main() {
	app := setup()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<h1>Chat APP | Backend</h1>")
	})
	app.Get("/ws", websocket.New(handlers.WsHandler))
	app.Listen(":3000")
}

func TestDBConnection() {

	ccuDb, ok, err := database.ConnectDB("ccu")
	if !ok {
		panic(err)
	}

	var strSQL string = "select id, game_id, title, description from jobs order by id desc limit 5"

	var jobs []types.Job
	result, err := ccuDb.Query(strSQL)
	// err = result.Scan()

	if err != nil {
		panic(err)
	}

	for result.Next() {
		var job types.Job
		err = result.Scan(&job.Id, &job.Game_id, &job.Title, &job.Description)
		if err != nil {
			panic(err)
		}
		jobs = append(jobs, job)
	}

	fmt.Printf("%v", jobs)
}
