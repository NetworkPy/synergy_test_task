package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func DoRequest() {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		log.Println(err)
		return
	}

	newCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req = req.WithContext(newCtx)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}
