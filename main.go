package main

import (
	"encoding/json"
	"fmt"

	"example.com/Chat-app/handlers"
	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

type DataEvent struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

var (
	app     *fiber.App
	ok      bool
	err     error
	clients map[string]string
	rooms   map[string]string
)

var pl = fmt.Println

func setup() *fiber.App {
	app = fiber.New()
	app.Use(cors.New())
	app.Use("/ws", handlers.Upgrade)
	return app
}

func chatSocketHandler() {
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {

		var data DataEvent
		err = json.Unmarshal(ep.Data, &data)

		if err != nil {
			pl(err.Error())
			return
		}

		if data.Type == "chat" {
			response := fiber.Map{
				"type":       data.Type,
				"clientData": data.Data,
				"response": fiber.Map{
					"userId":  data.Data["userId"],
					"message": data.Data["message"],
				},
			}

			msg, err := json.Marshal(response)
			if err != nil {
				pl(err.Error())
			}

			ep.Kws.Broadcast(msg, false)
		}
	})
}

func main() {
	app = setup()
	app.Static("/", "./public")

	clients = make(map[string]string)

	chatSocketHandler()

	app.Get("/ws", ikisocket.New(func(kws *ikisocket.Websocket) {

		userId := kws.Params("userId")
		if userId != "" {
			kws.SetAttribute("userId", userId)
		}

		clients[userId] = kws.UUID

	}))

	app.Listen(":3000")
}
