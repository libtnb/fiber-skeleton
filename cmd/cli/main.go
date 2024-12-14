package main

import (
	_ "time/tzdata"
)

func main() {
	cli, err := initCli()
	if err != nil {
		panic(err)
	}

	if err = cli.Run(); err != nil {
		panic(err)
	}
}
