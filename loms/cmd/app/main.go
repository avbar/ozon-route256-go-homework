package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

func main() {
	http.Handle("/createOrder", srvwrapper.New(createorder.New().Handle))
	http.Handle("/listOrder", srvwrapper.New(listorder.New().Handle))
	http.Handle("/orderPayed", srvwrapper.New(orderpayed.New().Handle))
	http.Handle("/cancelOrder", srvwrapper.New(cancelorder.New().Handle))
	http.Handle("/stocks", srvwrapper.New(stocks.New().Handle))

	log.Println("listening http at", port)
	err := http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
