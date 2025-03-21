package main

import (
	logger "challengerone/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func fetch(ctx context.Context) (*BidResponse, error) {
	const url = "http://localhost:8080/cotation"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error when creating request: %w", err)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error when making request: %w", err)
	}
	defer response.Body.Close()
	var data BidResponse
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
	logger.InitLogger("client.log")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	response, err := fetch(ctx)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	text := fmt.Sprintf("Dólar: %s", response.Bid)
	errSaveFile := saveToFile("cotacao.txt", text)
	if errSaveFile != nil {
		fmt.Println("Erro ao salvar o arquivo:", err)
		return
	}
	fmt.Printf("Dólar: %s", response.Bid)
}
