package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example.com/Chat-app/database"
	"example.com/Chat-app/handlers"
	"example.com/Chat-app/types"
	"example.com/Chat-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
	_ "github.com/joho/godotenv/autoload"
)

var (
	chatHistorys map[int64]fiber.Map
)

func initServer(engine *html.Engine) *fiber.App {
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New())
	app.Use("/ws", handlers.Upgrade)
	app.Static("/style/", "./public/style/")
	app.Static("/js/", "./public/js/")
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
	engine := html.New("./views", ".html")
	app := initServer(engine)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "My World",
		})
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		var (
			nsec int64
			mt   int
			msg  []byte
			err  error
			data types.WebSocketEvent
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
					"chatHistorys": utils.MapValues(chatHistorys),
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
