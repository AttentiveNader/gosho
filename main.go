package main

import (
	"./models"
	"./routes"
	"./utils"
	"net/http"
)

func main() {
	port := os.Getenv("PORT")
	utils.LoadTemplates("templates/*.gohtml")
	models.Dbc.InitializeDatabase()
	r := routes.InitRouter()
	//models.Dbc.InsertUrl("12341", "herokuapp2", "googlesearch alot of this kidadadand of shit")
	//port := os.Getenv("port")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
