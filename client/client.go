package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)


type BidDolar struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bid  string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var ColetaBid BidDolar // Aqui executa a função para salvar no arquivo
	err = json.Unmarshal(body, &ColetaBid)
	if err != nil {
		panic(err)
	}

	err = GravarCotacao(ColetaBid.Bid)
	if err != nil {
		println("Erro ao gravar cotacao")
		panic(err)
	}
}

func GravarCotacao(bid string) error { // Grava a cotação no arquivo
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Cotacao: %s\n", bid))
	if err != nil {
		return err
	}
	return nil
}
