package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type AgifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func main() {
	name := "timur"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	apiURL := "https://api.agify.io/?name=" + url.QueryEscape(name)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		log.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}
	var data AgifyResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}
	fmt.Println("Response from public API:")
	fmt.Println("name:", data.Name)
	fmt.Println("predicted age:", data.Age)
	fmt.Println("count:", data.Count)
}