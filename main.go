package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"testTask/handlers"
)

func main() {

	http.HandleFunc("/", handlers.PrimeHandler)

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)

	server := http.Server{}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	select {
	case <-shutdown:
		fmt.Println("Shutting down ...")
		server.Close()
		return
	}
}
