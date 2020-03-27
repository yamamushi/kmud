package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/utils"
)

func main() {

	log.Println("Checking config")
	conf, err := config.GetConfig("frontend.conf")
	if err != nil {
		utils.HandleError(err)
	}

	log.Println("Creating endpoint handlers")
	svc := stringService{}
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	log.Println("Registering endpoint handlers")
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)

	log.Println("Listening for connections...")
	err = http.ListenAndServe(conf.Server.Interface+":"+conf.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
