package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type CotacaoDolar struct {
	ID         string `json:"id"`
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
}

func NewCotacaoDolar(name string, bid string) *CotacaoDolar {
	return &CotacaoDolar{
		ID:   uuid.New().String(),
		Name: name,
		Bid:  bid,
	}
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Pagina nao encontrada")
		return
	}
	w.Write([]byte("Bem vindo ao servidor"))
	log.Println("Pagina carregada com sucesso")
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := BuscaCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro ao buscar cotacao")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Erro ao salvar cotacao")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Cotacao salva com sucesso"))
	log.Println("Cotacao salva com sucesso")
}

func BuscaCotacao() (*CotacaoDolar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cotacao CotacaoDolar
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func GravarCotacao(dado *CotacaoDolar) error {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/cotacaodolar?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = db.ExecContext(ctx, "INSERT INTO cotacaoDolar (id, name, bid) VALUES (?, ?, ?)", dado.ID, dado.Name, dado.Bid)
	if err != nil {
		return err
	}
	return nil
}