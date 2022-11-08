package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"example.com/Chat-app/database"
	"example.com/Chat-app/types"
	"example.com/Chat-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	mt        int
	msg       []byte
	err       error
	userId    int
	strUserId string
	data      types.WebSocketEvent
	ok        bool
)

var pl = fmt.Println

func Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WsHandler(c *websocket.Conn) {
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
			onConnected(c, data)
		}
	}
}

func onConnected(c *websocket.Conn, data types.WebSocketEvent) {

	strUserId, ok = data.Data["userId"].(string)
	if !ok {
		log.Println("userId is not a string")
		return
	}

	userId, err = strconv.Atoi(strUserId)
	if err != nil {
		log.Println("userId is not a number")
		return
	}

	ok, err := database.CreateRoom(userId)

	if err != nil {
		log.Println(err.Error())
		return
	}

	if !ok {
		return
	}

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
