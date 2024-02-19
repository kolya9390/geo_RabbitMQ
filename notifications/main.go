package main

import (
	"log"

	rpcserver "notific/rps_server"
)

func main() {

	rpcServer := rpcserver.NewGeoServis()
	err := rpcServer.StartServer()

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
