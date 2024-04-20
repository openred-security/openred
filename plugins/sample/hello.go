package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
)

func main() {
    fmt.Println("Hello, World!")

	unixSocket := "/var/run/my_socket.sock"
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

/*
resp, err := client.R().
      SetBody(User{Username: "testuser", Password: "testpass"}).
      SetResult(&AuthSuccess{}).    // or SetResult(AuthSuccess{}).
      SetError(&AuthError{}).       // or SetError(AuthError{}).
      Post("https://myapp.com/login")

}
*/