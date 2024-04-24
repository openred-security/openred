package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

// echo '{"value": "Hello, world!"}' | ./plugin-sample

func main() {
	fmt.Println("Hello, World!")

	unixSocket := "/tmp/echo.sock"
	// Create a Go's http.Transport so we can set it in resty.
	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", unixSocket)
		},
	}
	// Create a Resty Client
	client := resty.New()

	// Set the previous transport that we created, set the scheme of the communication to the
	// socket and set the unixSocket as the HostURL.
	client.SetTransport(&transport).SetScheme("http").SetBaseURL(unixSocket)

	// No need to write the host's URL on the request, just the path.
	/*
	   resp, err := client.R().
	       EnableTrace().
	       Get("http://localhost/")

	   // Explore response object
	   fmt.Println("Response Info:")
	   fmt.Println("  Error      :", err)
	   fmt.Println("  Status Code:", resp.StatusCode())
	   fmt.Println("  Status     :", resp.Status())
	   fmt.Println("  Proto      :", resp.Proto())
	   fmt.Println("  Time       :", resp.Time())
	   fmt.Println("  Received At:", resp.ReceivedAt())
	   fmt.Println("  Body       :\n", resp)
	   fmt.Println()

	   }

	   resp, err := client.R().
	         SetBody(User{Username: "testuser", Password: "testpass"}).
	         SetResult(&AuthSuccess{}).    // or SetResult(AuthSuccess{}).
	         SetError(&AuthError{}).       // or SetError(AuthError{}).
	         Post("https://myapp.com/login")

	   }
	*/

	// Listen to stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Input is being piped in
		var inputData map[string]interface{}
		decoder := json.NewDecoder(os.Stdin)
		err := decoder.Decode(&inputData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding JSON: %v\n", err)
			os.Exit(1)
		}

		// Print the received JSON data (optional)
		fmt.Println("Received JSON:", inputData)

		// Send JSON data via POST request
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(inputData).
			Post("http://localhost/send") // Replace with your endpoint

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending POST request: %v\n", err)
			os.Exit(1)
		}

		// Print response
		fmt.Println("Response:", resp.String())
	} else {
		fmt.Println("No input piped through stdin")
	}
}
