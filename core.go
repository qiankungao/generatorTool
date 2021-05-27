package main

import "database/sql"

var (
	Conn    *sql.DB  //连接对象
	DbConn  DBConfig //db config
	formats []string //format
)
