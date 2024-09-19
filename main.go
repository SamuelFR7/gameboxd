package main

import (
	"gameboxd/api"
)

func main() {
	app := api.NewServer()

	app.Listen(":3000")
}
