package main

import (
	"encoding/json"
	"log"
	"time"

	"example.com/Chat-app/handlers"
	"example.com/Chat-app/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type Chat struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

type WebSocketEvent struct {
	EventName string                 `json:"eventName"`
	Data      map[string]interface{} `json:"data"`
}

var (
	chatHistorys map[int64]fiber.Map
)

func initServer() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Static("/", "./public/")
	app.Use("/ws", handlers.Upgrade)
	return app
}

func send(c *websocket.Conn, mt int, msg []byte) (s bool, err error) {
	var status bool = true
	if err := c.WriteMessage(mt, msg); err != nil {
		log.Println("write:", err)
		status = false
	}

	return status, err

}

func main() {
	chatHistorys = make(map[int64]fiber.Map)
	app := initServer()
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			nsec int64
			mt   int
			msg  []byte
			err  error
			data WebSocketEvent
		)
		nsec = time.Now().UnixNano()

		for {

			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println("unmarshal:", err)
				break
			}

			if data.EventName == "connected" {
				if msg, err = json.Marshal(fiber.Map{
					"eventName":    data.EventName,
					"chatHistorys": helpers.MapValues(chatHistorys),
				}); err != nil {
					log.Println("marshal:", err)
					break
				}
				send(c, mt, msg)
			} else if data.EventName == "chat" {
				chatHistorys[nsec] = fiber.Map{
					"username": data.Data["username"].(string),
					"text":     data.Data["text"].(string),
				}
				send(c, mt, msg)
			}
		}
	}))
	app.Listen(":3000")
}
