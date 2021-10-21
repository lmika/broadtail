package main

import (
	"flag"
	"fmt"
	"github.com/lmika/broadtail/handlers"
	"log"
	"net/http"
)

func main() {
	flagBindAddr := flag.String("b", "localhost", "bind address")
	flagPort := flag.Int("p", 3690, "port")
	flag.Parse()

	server, err := handlers.Server()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Listening on %v:%v", *flagBindAddr, *flagPort)
	log.Println(http.ListenAndServe(fmt.Sprintf("%v:%v", *flagBindAddr, *flagPort), server))
}
