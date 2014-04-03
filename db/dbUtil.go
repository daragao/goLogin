package db

import (
	"database/sql"
	"fmt"
	"github.com/daragao/goLogin/logger"
	_ "github.com/lib/pq"
	"reflect"
	"strconv"
	"strings"
)

func connectDB() (con *sql.DB, err error) {
	con, err = sql.Open("postgres", connection_string)
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
		logger.ERRO.Println("Failed to create DB: " + err.Error() + "\n\t" + create_users_table)
	}
}

func insertRowStr(tablename string, columns ...string) (sqlStatement string) {
	colValuesArr := make([]string, len(columns))
	colNames := strings.Join(columns, ",")
	for i, _ := range columns {
		colValuesArr[i] = "$" + strconv.Itoa(i+1)
	}
	colValues := strings.Join(colValuesArr, ",")
	sqlStatement = "INSERT INTO " + tablename + " (" + colNames +
		") VALUES (" + colValues + ")"
	return
}

/* SHOULD FIND A WAY TO MAKE THIS MORE GENERIC AND PUT IT IN THE USERS PACKAGE*/
func InsertNewUser(username string, password []byte) (err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := insertRowStr(users_table_name, "username", "password")
	_, err = con.Exec(sqlStatement, username, password)
	if err != nil {
		logger.ERRO.Println("Failed to insert user: " + err.Error() + "\n\t" + sqlStatement)
	}
	return
}

func GetRowBy(tablename string,
	colNames string,
	attribName string,
	id interface{}) (row *sql.Row, err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := "SELECT " + colNames + " FROM " + tablename + " WHERE " + attribName + " = $1"
	row = con.QueryRow(sqlStatement, id)
	if err != nil {
		logger.ERRO.Println("Failed to get row: " + err.Error() + "\n\t" + sqlStatement)
	}
	return
}

func GetRows(tablename string, colNames string, limit int, offset int) (rows *sql.Rows, err error) {
	con, err := connectDB()
	defer con.Close()
	sqlStatement := "SELECT " + colNames + " FROM " + tablename + " LIMIT $1 OFFSET $2"
	rows, err = con.Query(sqlStatement, limit, offset)
	if err != nil {
		logger.ERRO.Println("Failed to get row: " + err.Error() + "\n\t" + sqlStatement)
	}
	return
}

func printRow(rows *sql.Rows) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Create interface set
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Scan for arbitrary values
	err = rows.Scan(scanArgs...)
	if err == nil {

		// Print data
		for i, value := range values {
			switch value.(type) {
			default:
				fmt.Printf("%s :: %s :: %+v\n", columns[i], reflect.TypeOf(value), value)
			}
		}
	} else {
		panic(err)
	}
}
