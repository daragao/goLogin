package db

var users_table_name = "users"
var timestamp_columns = "last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP " +
	"NOT NULL ON UPDATE CURRENT_TIMESTAMP"
var create_users_table = "CREATE TABLE IF NOT EXISTS " + users_table_name +
	" (id int(4) NOT NULL auto_increment,username varchar(32) NOT NULL UNIQUE," +
	"password varchar(32) NOT NULL,PRIMARY KEY (id), " + timestamp_columns +
	") ENGINE=MyISAM;"
