package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
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

func fetch(ctx context.Context) (*Response, error) {
	const url = "http://localhost:3333/cotation"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error when creating request: %w", err)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error when making request: %w", err)
	}
	defer response.Body.Close()
	var data Response
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %w", err)
	}
	return &data, nil
}

func saveToFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	response, err := fetch(ctx)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	text := fmt.Sprintf("Dólar: %s", response.Usdbrl.Bid)
	errSaveFile := saveToFile("cotacao.txt", text)
	if errSaveFile != nil {
		fmt.Println("Erro ao salvar o arquivo:", err)
		return
	}
	fmt.Print(response)
	fmt.Println("Hello, World!")
	fmt.Println("This is a Go program")
}
