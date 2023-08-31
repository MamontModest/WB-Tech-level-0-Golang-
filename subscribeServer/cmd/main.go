package main

import (
	"context"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	"log"
	"net/http"
	"os"
	cash2 "subscribe/internal/cash"
	"subscribe/internal/nadlers"
	"subscribe/internal/subscriber"
	"subscribe/pkg/db"
)

func main() {
	sub := subscriber.New()
	sub.ConnectToSubscribe()

	d := db.NewSDatabase()
	_, err := d.ConnWith(context.Background())
	if err != nil {
		panic(err)
	}
	cash := cash2.NewInit()
	router := routing.New()
	router.Use(
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	address := fmt.Sprintf(":%v", "8080")
	nadlers.RegisterHandlers(router.Group(""), cash)
	hs := &http.Server{
		Addr:    address,
		Handler: router,
	}

	log.Println(fmt.Sprintf("server listen on address localhost%s", address))

	err = hs.ListenAndServe()
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
