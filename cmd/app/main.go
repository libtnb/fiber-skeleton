package main

import (
	_ "time/tzdata"
)

func main() {
	app, err := initApp()
	if err != nil {
		panic(err)
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
