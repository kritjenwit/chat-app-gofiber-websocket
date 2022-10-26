package types

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
