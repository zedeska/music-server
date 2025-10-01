package main

import (
	"encoding/json"
	"fmt"
	"music-server/deezer"
)

var API_URL string = "https://api.deezer.com/"

func main() {
	results, err := deezer.Search("Deep Blue 9mm parabellum bullet")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))

}
