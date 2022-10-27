package config

import (
	"os"

	"example.com/Chat-app/constants"
	"example.com/Chat-app/helpers"
	"example.com/Chat-app/types"
)

var (
	DBConfig map[string]types.DBConfig
)

func setup() {
	DBConfig["ccu"] = types.DBConfig{
		Driver:   constants.MYSQL_DRIVER,
		Host:     os.Getenv("DBHOST"),
		Port:     helpers.ToInt(os.Getenv("DBPORT")),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Database: os.Getenv("DBNAME"),
	}
}

func GetDBConfig(s string) types.DBConfig {
	if DBConfig == nil {
		DBConfig = make(map[string]types.DBConfig)
	}
	setup()
	return DBConfig[s]
}
