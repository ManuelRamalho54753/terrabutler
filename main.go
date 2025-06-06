package main

import (
	"github.com/montblu/terrabutler/internal/cmd"
	"github.com/montblu/terrabutler/internal/logger"
)

func main() {

	logger.InitLogger()
	defer logger.Log.Sync()

	cmd.Execute()
}
