package bootstrap

import "runtime/debug"

func Boot() {
	debug.SetGCPercent(10)
	debug.SetMemoryLimit(128 << 20)

	initConf()
	initOrm()
	runMigrate()
	initValidator()
	initHttp()
}
