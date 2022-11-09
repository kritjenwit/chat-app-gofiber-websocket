package database

import (
	"database/sql"
	"errors"
	"time"

	"example.com/Chat-app/types"
)

var (
	strSQL string
	row    *sql.Row
	rows   *sql.Rows
	result sql.Result
)

func setup() (*sql.DB, error) {
	db, ok, err := ConnectDB("ccu")

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if !ok {
		return nil, errors.New("could not connect to database")
	}

	return db, nil
}

func CreateRoom(userId int) (bool, error) {
	db, err := setup()
	if err != nil {
		return false, err
	}

	strSQL = "select user_id, create_time from room_online where user_id = ?"
	row := db.QueryRow(strSQL, userId)

	if err := row.Err(); err != nil {
		return false, err
	}

	var roomOnline types.RoomOnline

	row.Scan(&roomOnline.UserID, &roomOnline.CreateTime)

	if roomOnline.UserID != 0 || roomOnline.CreateTime != "" {
		return false, errors.New("room found")
	}

	strSQL = "insert into room_online (user_id, create_time) values (?, ?)"
	result, err = db.Exec(strSQL, userId, time.Now())
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affectedRows > 0 {
		return true, nil
	}
	return false, errors.New("cannot create room")
}

func GetAllRooms() (map[int]types.RoomOnline, error) {
	db, err := setup()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("select user_id, create_time from room_online order by create_time desc")
	if err != nil {
		return nil, err
	}

	result := make(map[int]types.RoomOnline)

	for rows.Next() {
		var roomOnline types.RoomOnline
		var userID int
		var createTime string
		err = rows.Scan(&userID, &createTime)
		if err != nil {
			return nil, err
		}
		roomOnline.UserID = userID
		roomOnline.CreateTime = createTime

		result[roomOnline.UserID] = roomOnline
	}

	return result, nil
}

func GetChatLog() ([]types.ChatLog, error) {
	db, err := setup()
	if err != nil {
		return nil, err
	}

	var chatLogs []types.ChatLog

	rows, err := db.Query("select * from chat_log where room_id = ? order by create_time asc")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var chatLog types.ChatLog
		err = rows.Scan(&chatLog.RoomID, &chatLog.UserID, &chatLog.Text, &chatLog.CreateTime)
		if err != nil {
			return nil, err
		}

		chatLogs = append(chatLogs, chatLog)
	}

	return chatLogs, err
}

func InsertChatLog(roomId int, userId int, text string) error {
	db, err := setup()
	if err != nil {
		return err
	}

	strSQL = "insert into chat_log (room_id, user_id, text, create_time) values (?, ?, ?, ?)"
	result, err = db.Exec(strSQL, roomId, userId, text, time.Now())
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows > 0 {
		return nil
	}

	result.LastInsertId()

	return err
}
