package main

import (
	"terrabutler/cmd"
	"terrabutler/logger"
)

func main() {

	logger.InitLogger()
	defer logger.Log.Sync()

	cmd.Execute()
}
