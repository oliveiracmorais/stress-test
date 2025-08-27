// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")

	flag.Parse()

	// Validação dos parâmetros
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		log.Fatal("Erro: --url, --requests e --concurrency são obrigatórios e devem ser válidos")
	}

	start := time.Now()

	// Canais para jobs e resultados
	jobs := make(chan struct{}, *requests)
	results := make(chan int, *requests)

	// Preenche o canal de jobs
	for i := 0; i < *requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	var wg sync.WaitGroup
	var totalRequests int
	var mu sync.Mutex

	// Cliente HTTP com timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Cria a barra de progresso
	bar := progressbar.NewOptions(*requests,
		progressbar.OptionSetDescription("Executando testes..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
	)

	// Workers concorrentes
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				resp, err := client.Get(*url)
				if err != nil {
					mu.Lock()
					totalRequests++
					mu.Unlock()
					results <- 0
					bar.Add(1) // avança a barra
					continue
				}

				// Lê e descarta o corpo
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				mu.Lock()
				totalRequests++
				mu.Unlock()
				results <- resp.StatusCode
				bar.Add(1) // avança a barra
			}
		}()
	}

	// Fechar canal de resultados quando todos os workers terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coleta resultados
	statusCodes := make(map[int]int)
	for code := range results {
		statusCodes[code]++
	}

	elapsed := time.Since(start)

	// Relatório final
	fmt.Println("=== RELATÓRIO DE TESTE DE CARGA ===")
	fmt.Printf("URL: %s\n", *url)
	fmt.Printf("Total de Requests: %d\n", *requests)
	fmt.Printf("Concorrência: %d\n", *concurrency)
	fmt.Printf("Tempo Total: %v\n", elapsed)
	fmt.Println("-----------------------------------")

	if count, ok := statusCodes[200]; ok {
		fmt.Printf("Status 200 (OK): %d requests\n", count)
		delete(statusCodes, 200)
	} else {
		fmt.Printf("Status 200 (OK): 0 requests\n")
	}

	if len(statusCodes) > 0 {
		fmt.Println("Outros códigos de status:")
		for code, count := range statusCodes {
			if code == 0 {
				fmt.Printf("  Erros de rede/conexão: %d requests\n", count)
			} else {
				fmt.Printf("  %d: %d requests\n", code, count)
			}
		}
	} else {
		fmt.Println("Outros códigos de status: Nenhum")
	}

	fmt.Println("===================================")
}
