package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	
	"primeServer/internal/handlers"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var port = "8080"

func init() {
	flag.String("port", port, "port for server to listen on")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	providedPort := viper.GetString("port")
	port = ":" + providedPort
}

func main() {

	http.HandleFunc("/", handlers.PrimeHandler)

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, os.Interrupt)

	server := http.Server{
		Addr: port,
	}

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
