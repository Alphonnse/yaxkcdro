package main

import (
	"log"

	"github.com/Alphonnse/yaxkcdro/internal/app"
)

func main(){
	// myapp := app.NewApp()
	// myapp.RunApp()

	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.RunApp()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
