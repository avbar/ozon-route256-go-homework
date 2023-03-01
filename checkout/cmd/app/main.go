package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := loms.New(config.ConfigData.Services.Loms)
	productClient := productservice.New(config.ConfigData.Services.ProductService)
	businessLogic := domain.New(lomsClient, productClient)

	http.Handle("/addToCart", srvwrapper.New(addtocart.New(businessLogic).Handle))
	http.Handle("/deleteFromCart", srvwrapper.New(deletefromcart.New(businessLogic).Handle))
	http.Handle("/listCart", srvwrapper.New(listcart.New(businessLogic).Handle))
	http.Handle("/puchase", srvwrapper.New(purchase.New(businessLogic).Handle))

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
