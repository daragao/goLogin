package db

import (
	"../logger"
	"database/sql"
	_ "github.com/ziutek/mymysql/godrv"
	"strings"
)

var dbVars = struct {
	database string
	user     string
	password string
}{"goGeteon", "goUser", "goUser"}

func connectDB() (con *sql.DB, err error) {
	con, err = sql.Open("mymysql", dbVars.database+"/"+dbVars.user+"/"+dbVars.password)
	if err != nil {
		logger.ERRO.Println("Failed to connect to DB: " + err.Error())
	}
	return
}

func CreateTables() {
	con, err := connectDB()
	defer con.Close()
	_, err = con.Exec(create_users_table)
	if err != nil {
		logger.ERRO.Println("Failed to create DB: " + err.Error())
	}
}

func insertRowStr(tablename string, columns ...string) (sqlStatement string) {
	colValuesArr := make([]string, len(columns))
	colNames := strings.Join(columns, ",")
	for i, _ := range columns {
		colValuesArr[i] = "?"
	}
	colValues := strings.Join(colValuesArr, ",")
	sqlStatement = "INSERT IGNORE INTO " + tablename + " (" + colNames +
		") VALUES (" + colValues + ")"
	return
}

func InsertNewUser(username string, password []byte) (err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := insertRowStr(users_table_name, "username", "password")
	_, err = con.Exec(sqlStatement, username, password)
	return
}

func getRowBy(tablename string, attribName string, id int) (row *sql.Row, err error) {
	con, err := connectDB()
	defer con.Close()
	row = con.QueryRow("SELECT * FROM " + tablename + " WHERE " + attribName + " = ?")
	return
}

func GetUserByUsername(usernameArg string) (id int, username string, password []byte) {
	row, _ := getRowBy(users_table_name, "username", id)
	err := row.Scan(&id, &username, &password)
	if err != nil {
		logger.ERRO.Println("Failed to scan DB row: " + err.Error())
	}
	return
}

func GetUserByID(idArg int) (id int, username string, password []byte) {
	row, _ := getRowBy(users_table_name, "id", idArg)
	err := row.Scan(&id, &username, &password)
	if err != nil {
		logger.ERRO.Println("Failed to scan DB row: " + err.Error())
	}
	return
}
