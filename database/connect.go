package database

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/Chat-app/config"
	_ "github.com/go-sql-driver/mysql"
)

var DbConnect map[string]*sql.DB

func ConnectDB(name string) (db *sql.DB, ok bool, err error) {

	if DbConnect == nil {
		DbConnect = make(map[string]*sql.DB)
	}
	if DbConnect[name] != nil {
		return DbConnect[name], true, nil
	}

	dbConfig := config.DBConfig[name]
	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err = sql.Open(dbConfig.Driver, dns)
	if err != nil {
		return nil, false, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	DbConnect[name] = db
	return DbConnect[name], true, nil
}
