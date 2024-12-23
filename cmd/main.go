package main

import (
	"github.com/safayildirim/asset-management-service/app"
	"github.com/safayildirim/asset-management-service/pkg/log"
)

func main() {
	err := app.New().Run()
	if err != nil {
		log.Logger.Sugar().Fatal(err)
	}
}
