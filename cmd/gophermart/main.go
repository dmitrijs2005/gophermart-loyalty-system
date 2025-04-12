package main

import "github.com/dmitrijs2005/gophermart-loyalty-system/internal/app"

func main() {
	app := app.NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}

}
