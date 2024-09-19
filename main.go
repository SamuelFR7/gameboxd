package main

import (
	"gameboxd/api/api"
)

func main() {
	app := api.NewServer()

	app.Listen(":3000")
}
