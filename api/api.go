package api

import (
	"net/http"
	"net"
	"github.com/labstack/echo/v4"
	"openred/openred-agent/sender"
	"log"

)

func Start() {

	client = sender.New()
	sender.CreateIndex(client)
	
	listener, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal(err)
	}
	
	e := echo.New()
	e.Listener = listener

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "BYE, World!")
	})

	server := new(http.Server)
	if err := e.StartServer(server); err != nil {
		log.Fatal(err)
	}


	



}