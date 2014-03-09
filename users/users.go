package users

import (
	"../auth"
	"../db"
	"../logger"
)

type User struct {
	// contains filtered or unexported fields
	Id       int
	Username string
	password []byte
}

func GetUserByID(idArg int) (userRow *User, err error) {
	row, err := db.GetRowBy("users", "id,username,password", "id", idArg)
	userRow = &User{}
	err = row.Scan(&userRow.Id, &userRow.Username, &userRow.password)
	return
}

func GetUserByUsername(usernameArg string) (userRow *User, err error) {
	row, err := db.GetRowBy("users", "id,username,password", "username", usernameArg)
	userRow = &User{}
	err = row.Scan(&userRow.Id, &userRow.Username, &userRow.password)
	return
}

func RegisterUser(username string, password string) (userRow *User, err error) {
	cryptPass, err := auth.Crypt(password)
	if err != nil {
		logger.ERRO.Println("Password couldn't be crypted: " + err.Error())
		return
	}
	err = db.InsertNewUser(username, cryptPass)
	if err == nil {
		userRow, err = GetUserByUsername(username)
	} else {
		userRow = nil
	}
	return
}

func GetAllUsers(limit int, offset int) (users []User, err error) {
	rows, err := db.GetRows("users", "id,username,password", limit, offset)
	users = make([]User, 0)
	for rows.Next() {
		userRow := User{}
		err = rows.Scan(&userRow.Id, &userRow.Username, &userRow.password)
		users = append(users, userRow)
	}
	return
}
