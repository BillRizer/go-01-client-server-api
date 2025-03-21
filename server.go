package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func getCotation(w http.ResponseWriter, r *http.Request) {
	response, err := fetch()
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/cotation", getCotation)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
