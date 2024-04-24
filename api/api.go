package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"openred/openred-agent/sender"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// curl --unix-socket /tmp/echo.sock http://localhost/send -X POST -d '{"message": "Hello, world!"}'

func Start() {

	// Path to the file you want to remove
	filePath := "/tmp/echo.sock"

	// Attempt to remove the file
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error removing file: %v\n", err)
		return
	}

	listener, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	e.Listener = listener

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "BYE, World!")
	})

	e.POST("/send", sendPayload)

	server := new(http.Server)
	if err := e.StartServer(server); err != nil {
		log.Fatal(err)
	}

}

func sendPayload(c echo.Context) error {

	type ExampleRequest struct {
		Value1 string `json:"value" form:"value" query:"value"`
	}

	exampleRequest := new(ExampleRequest)
	if err := c.Bind(exampleRequest); err != nil {
		return err
	}

	senderClient := sender.New()
	sender.CreateIndex(senderClient)
	sender.Send(senderClient, "{\"value\":\"field\"}")

	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf("%s", "All good!"),
		),
	)
}
