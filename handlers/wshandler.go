package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/Chat-app/database"
	"example.com/Chat-app/types"
	"example.com/Chat-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	nsec int64
	mt   int
	msg  []byte
	err  error
	data types.WebSocketEvent
)

func Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WsHandler(c *websocket.Conn) {

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

		userId := data.Data["userId"]

		fmt.Println(userId)

		if data.EventName == "connected" {
			rooms, err := database.GetAllRooms()
			if err != nil {
				log.Println(err)
				panic(err)
			}

			if rooms == nil {
				rooms = make(map[int]types.RoomOnline, 0)
			}

			response := fiber.Map{
				"emitName": data.EventName,
				"data": fiber.Map{
					"rooms": utils.MapValues(rooms),
				},
			}

			if msg, err = json.Marshal(response); err != nil {
				log.Println(err)
				panic(err)
			}

			utils.WriteMessage(c, mt, msg)
		}
	}
}
