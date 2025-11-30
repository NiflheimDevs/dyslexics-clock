package main

import (
	"log"
	"net/http"

	"github.com/NiflheimDevs/dyslexics-clock/bootstrap"
	"github.com/NiflheimDevs/dyslexics-clock/internal/delivery/route"
	"github.com/NiflheimDevs/dyslexics-clock/wire"
)

func main() {
	di := bootstrap.Get()

	app ,err :=wire.InitApp(di)
	if err != nil {
		panic (err)
	}
	
	log.Printf("Application is running on port%s", di.Const.Port)
	
	server := &http.Server{
		Addr:    di.Const.Port,
		Handler: route.RouteInit(app),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
