package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"example.com/Chat-app/constants"
	"example.com/Chat-app/database"
	"example.com/Chat-app/types"
	"example.com/Chat-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	mt           int
	msg          []byte
	err          error
	userId       int
	data         types.WebSocketEvent
	roomUserJoin map[int]map[int]int
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

		if data.EventName == constants.WS_EVENT_CONNECTED {
			onConnected(c, data)
		} else if data.EventName == constants.WS_EVENT_JOIN_ROOM {
			onJoinRoom(c, data)
		} else if data.EventName == constants.WS_EVENT_SEND_MESSAGE {
			onSendMessage(c, data)
		} else if data.EventName == constants.WS_EVENT_CREATE_ROOM {
			onCreateRoom(c, data)
		}
	}
}

func onConnected(c *websocket.Conn, data types.WebSocketEvent) {

	rooms, err := database.GetAllRooms()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if rooms == nil {
		rooms = make(map[int]types.RoomOnline, 0)
	}

	response := fiber.Map{
		"eventName": data.EventName,
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

func onJoinRoom(c *websocket.Conn, data types.WebSocketEvent) {

	if roomUserJoin == nil {
		roomUserJoin = make(map[int]map[int]int)
	}

	strRoomId, ok := data.Data["roomId"].(string)
	if !ok {
		log.Println("roomId is not a string")
		return
	}

	roomId, err := strconv.Atoi(strRoomId)
	if err != nil {
		log.Println("roomId is not a number")
		return
	}

	strUserId, ok := data.Data["userId"].(string)
	if !ok {
		log.Println("userId is not a string")
		return
	}

	userId, err := strconv.Atoi(strUserId)
	if err != nil {
		log.Println("userId is not a number")
		return
	}

	if roomUserJoin[roomId] == nil {
		roomUserJoin[roomId] = make(map[int]int, 0)
	}

	roomUserJoin[roomId][userId] = userId

	response := fiber.Map{
		"eventName": data.EventName,
		"data": fiber.Map{
			"status": "ok",
			"roomId": roomId,
			"userId": userId,
			"users":  roomUserJoin,
		},
	}

	if msg, err = json.Marshal(response); err != nil {
		log.Println(err)
		panic(err)
	}

	utils.WriteMessage(c, mt, msg)
}

func onCreateRoom(c *websocket.Conn, data types.WebSocketEvent) {
	if roomUserJoin == nil {
		roomUserJoin = make(map[int]map[int]int, 0)
	}

	strRoomId, ok := data.Data["roomId"].(string)
	if !ok {
		log.Println("userId is not a string")
		return
	}

	roomId, err := strconv.Atoi(strRoomId)
	if err != nil {
		log.Println("userId is not a number")
		return
	}

	_, err = database.CreateRoom(roomId)

	if err != nil {
		log.Println(err.Error())
	}

	response := fiber.Map{
		"eventName": data.EventName,
		"data": fiber.Map{
			"status": "ok",
			"roomId": roomId,
		},
	}

	if msg, err = json.Marshal(response); err != nil {
		log.Println(err)
		panic(err)
	}

	if roomUserJoin[roomId] == nil {
		roomUserJoin[roomId] = make(map[int]int, 0)
	}

	roomUserJoin[roomId][userId] = userId
	// if roomUserJoin[roomId][userId] == nil {
	// 	make(roomUserJoin[roomId][userId])
	// }

	utils.WriteMessage(c, mt, msg)
}

func onSendMessage(c *websocket.Conn, data types.WebSocketEvent) {
	strUserId, ok := data.Data["userId"].(string)
	if !ok {
		log.Println("userId is not a string")
		return
	}

	userId, err = strconv.Atoi(strUserId)
	if err != nil {
		log.Println("userId is not a number")
		return
	}

	strRoomId, ok := data.Data["roomId"].(string)
	if !ok {
		log.Println("roomId is not a string")
		return
	}

	roomId, err := strconv.Atoi(strRoomId)
	if err != nil {
		log.Println("roomId is not a number")
		return
	}

	text, ok := data.Data["text"].(string)
	if !ok {
		log.Println("text is not a string")
		return
	}

	err = database.InsertChatLog(roomId, userId, text)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	utils.WriteMessage(c, mt, []byte("ok"))
}
