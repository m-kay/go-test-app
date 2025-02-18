package main

import "github.com/m-kay/test-app/src/app"

func main() {
	app := app.NewApplication()
	err := app.Start(3000)

	if err != nil {
		panic(err)
	}
}
