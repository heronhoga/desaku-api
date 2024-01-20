package main

import (
	"desaku-api/databases"
	"desaku-api/routes"
)

func main() {
	databases.ConnectDatabase()
	routes.Route()
}