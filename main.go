package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errorRequestFailed = errors.New("Request Failed")

func main() {
	urls := []string{
		"https://www.airbnb.co.kr/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://nomadcoders.co/",
	}

	for _, url := range urls {
		hitURL(url)
	}
}

func hitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return errorRequestFailed
	}
	return nil
}