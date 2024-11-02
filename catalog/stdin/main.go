package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// echo '{"value": "Hello, world!"}' | ./plugin-sample

func main() {

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

	} else {
		fmt.Println("No input piped through stdin")
	}
}
