package db

var dbVars = struct {
	database string
	user     string
	password string
}{"goGeteon", "goUser", "goUser"}

//var connection_string = "host=localhost user=postgres password='flipflop' " +
	"dbname=mydb"

var connection_string = "host=ec2-54-225-255-208.compute-1.amazonaws.com user=ekbsjwmuusakrv password='Rhiuimvu1Yy3xXVhhPFXEX2CHx' " +
	"dbname=d62qrfrl9jaggm"
var users_table_name = "users"

var timestamp_columns = ", date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP " //+
//	"NOT NULL ON UPDATE CURRENT_TIMESTAMP"

var create_users_table = "CREATE TABLE IF NOT EXISTS " + users_table_name +
	" (id BIGSERIAL PRIMARY KEY,username varchar(32) NOT NULL UNIQUE," +
	"password varchar(60) NOT NULL, is_admin boolean DEFAULT false" + timestamp_columns + ");"
