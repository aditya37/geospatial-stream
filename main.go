package main

import (
	"log"

	"github.com/aditya37/geospatial-stream/service"
)

func main() {
	svc, err := service.NewService()
	if err != nil {
		log.Panic(err)
		return
	}
	svc.Run()
}
