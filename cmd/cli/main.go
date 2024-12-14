package main

import (
	"runtime/debug"
	_ "time/tzdata"
)

func main() {
	debug.SetGCPercent(10)

	cli, err := initCli()
	if err != nil {
		panic(err)
	}

	if err = cli.Run(); err != nil {
		panic(err)
	}
}
