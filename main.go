package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"pintap/config"
	"pintap/utils"
)

func main() {
	var port string

	flag.StringVar(&port, "port", os.Getenv("PORT"), "port of the service")

	db := utils.GetDBConnection()
	defer db.Close()

	routes := &config.Routes{DB:db}

	routes.Setup(port)

}