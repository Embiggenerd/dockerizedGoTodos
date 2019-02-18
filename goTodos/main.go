package main

import (
	"goTodos/models"
	"goTodos/routes"

	_ "github.com/lib/pq"
)

func main() {
	models.Init()
	routes.Init()
}
