package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	logger "challengerone/logger"

	_ "github.com/mattn/go-sqlite3"
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
type BidResponse struct {
	Bid string `json:"bid"`
}

func saveCotation(ctx context.Context, data Response) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS cotations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		bid TEXT,
		ask TEXT,
		create_date TEXT
	);`
	_, err = db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("erro ao criar a tabela: %v", err)
	}
	insertSQL := `INSERT INTO cotations (code,codein,bid,ask,create_date) VALUES (?, ?, ?, ?, ?)`
	_, err = db.ExecContext(ctx, insertSQL, data.Usdbrl.Code, data.Usdbrl.Codein, data.Usdbrl.Bid, data.Usdbrl.Ask, data.Usdbrl.CreateDate)
	if err != nil {
		return fmt.Errorf("erro ao inserir os dados: %v", err)
	}
	return nil
}

func fetch(ctx context.Context) (*Response, error) {
	const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}
	response, err := http.DefaultClient.Do(req)
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
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	response, err := fetch(ctx)
	if err != nil {
		log.Printf("Erro ao buscar a cotação: %v\n", err)
		return
	}
	ctxDb, cancelDb := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDb()
	errSave := saveCotation(ctxDb, *response)
	if errSave != nil {
		log.Printf("Erro ao salvar a cotação: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(BidResponse{
		Bid: response.Usdbrl.Bid,
	})
}

func main() {
	logger.InitLogger("server.log")

	log.Printf("Starting server...")
	http.HandleFunc("/cotation", getCotation)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
