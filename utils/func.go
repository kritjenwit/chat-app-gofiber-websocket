package utils

import (
	"log"
	"strconv"

	"github.com/gofiber/websocket/v2"
)

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func ToInt(s string) int {
	intVar, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return intVar
}

func WriteMessage(c *websocket.Conn, mt int, msg []byte) (s bool, err error) {
	var status bool = true
	if err := c.WriteMessage(mt, msg); err != nil {
		log.Println("write:", err)
		status = false
	}
	return status, err
}
