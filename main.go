package main

import (
	"log"

	"github.com/taofit/GPA-calculator/api"
	"github.com/taofit/GPA-calculator/database"
)

func main() {
	db, err := database.NewDBInstance()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConn()

	server := api.NewAPIServer(":8080", db)

	server.SeedDB()
	server.Run()
}
