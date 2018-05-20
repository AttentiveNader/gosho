package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type mysqlDB struct { //a struct for controling the database and error handling
	Db   *sql.DB
	Name string
	Err  error
}

var Dbc *mysqlDB = &mysqlDB{}

func (dbc *mysqlDB) InitializeDatabase() {
	dbc.Name = os.Getenv("DATABASE_NAME")
	user := os.Getenv("DATABASE_USER")
	server := os.Getenv("DATABASE_SERVER")
	pass := os.Getenv("DATABASE_PASSWORD")
	dbc.Db, dbc.Err = sql.Open("mysql",
		user+":"+pass+"@tcp"+"("+server+")"+"/"+dbc.Name)
	if dbc.Err != nil {
		panic(dbc.Err)
	}
}

func (dbc *mysqlDB) InsertUrl(id, shorturl, longurl string) {

	stmt, err := dbc.Db.Prepare("INSERT INTO " + dbc.Name + "." + "urls" + " (id,long_url , short_url) VALUES (?,?,?) ")
	dbc.Err = err
	if dbc.Err != nil {
		fmt.Println("Insert URL dbc.Err >>>>>>>>>>>>>>> 1", dbc.Err)
		return
	}
	defer stmt.Close()

	_, dbc.Err = stmt.Exec(id, longurl, shorturl)

	if dbc.Err != nil {
		fmt.Println("Insert URL dbc.Err >>>>>>>>>>>>>>> 2", dbc.Err)
	}
}

func (dbc *mysqlDB) GetUrl(id string) string {

	row := dbc.Db.QueryRow("SELECT long_url FROM `" + dbc.Name + "`." + "`urls`" + " WHERE id=" + `"` + id + `"`)
	var longurl string

	dbc.Err = row.Scan(&longurl)
	if dbc.Err != nil {
		longurl = ""
	}

	return longurl
}
func (dbc *mysqlDB) GetShort(longurl string) string {

	row := dbc.Db.QueryRow("SELECT short_url FROM `" + dbc.Name + "`." + "`urls`" + " WHERE long_url=" + `"` + longurl + `"`)
	var short string

	dbc.Err = row.Scan(&short)
	if dbc.Err != nil {
		dbc.Err = nil
		short = ""
	}
	return short
}
