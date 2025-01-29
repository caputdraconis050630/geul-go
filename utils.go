package main

import (
	"log"
	"net/http"
)

func HandleStatusCodeError(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatal("Failed to get the page. Status code: ", resp.StatusCode)
	}
}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
