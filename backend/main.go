package main

import (
	"log"
)

func main() {
	config := LoadConfig()
	service := NewSurveyService(NewMissionStore())
	log.Printf("drone survey ops listening on :%s", config.Port)
	log.Fatal(serveAddress(":"+config.Port, NewRouter(service)))
}
