package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func fetchTitle(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done() // Marca a goroutine como concluída no WaitGroup

	// Faz uma requisição HTTP GET
	req, err := http.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Erro ao acessar %s: %v", url, err)
		return
	}
	defer req.Body.Close()

	// Verifica se a requisição foi bem-sucedida
	if req.StatusCode != 200 {
		results <- fmt.Sprintf("Erro ao acessar %s: status %d %s", url, req.StatusCode, req.Status)
		return
	}

	// Carrega documento HTML
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		results <- fmt.Sprintf("Erro ao carregar documento de %s: %v", url, err)
		return
	}

	// Seleciona o título da página e retorna
	title := doc.Find("title").Text()
	results <- fmt.Sprintf("Título de %s: %s", url, title)
}

func main() {
	urls := []string{
		"https://www.google.com/",
		"https://www.github.com/",
		"https://www.golang.org/",
	}

	var wg sync.WaitGroup
	results := make(chan string, len(urls)) // Canal para armanzenar os resultados

	for _, url := range urls {
		wg.Add(1)
		go fetchTitle(url, &wg, results)
	}

	// Aguarda todas as goroutines terminarem
	wg.Wait()
	close(results) // Fecha o acesso ao canal

	// Imprime os resultados
	for result := range results {
		fmt.Println(result)
	}
}
