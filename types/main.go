package types

type Chat struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

type WebSocketEvent struct {
	EventName string                 `json:"eventName"`
	Data      map[string]interface{} `json:"data"`
}

type DBConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Database string `json:"database"`
}

type Job struct {
	Id          int    `json:"id"`
	Game_id     int    `json:"game_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
