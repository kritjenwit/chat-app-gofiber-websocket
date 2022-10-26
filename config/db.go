package config

import "example.com/Chat-app/types"

var DBConfig = map[string]types.DBConfig{
	"ccu": {
		Driver:   "mysql",
		Host:     "210.246.248.22",
		Port:     3306,
		User:     "root",
		Password: "shk,g-hk",
		Database: "mis_db",
	},
}
