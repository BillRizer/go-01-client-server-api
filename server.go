package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func fetch() (*Response, error) {
	const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching data: %w", err)
	}
	defer response.Body.Close()
	var data Response
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}
	return &data, nil
}

func main() {
	fmt.Println("server cotation dolar")
	response, err := fetch()
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	fmt.Printf("High: %s\n", response.Usdbrl.High)
	fmt.Printf("Low: %s\n", response.Usdbrl.Low)
}
