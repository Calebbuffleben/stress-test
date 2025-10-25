package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type resultsSummary struct {
	totalRequests   int
	ok200Count      int
	statusHistogram map[int]int
	duration        time.Duration
}

func main() {
	serviceURL := flag.String("url", "", "URL do serviço a ser testado")
	totalRequests := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 0, "Número de chamadas simultâneas")
	flag.Parse()

	if strings.TrimSpace(*serviceURL) == "" {
		fmt.Fprintln(os.Stderr, "erro: parâmetro --url é obrigatório")
		os.Exit(2)
	}
	if _, err := url.ParseRequestURI(*serviceURL); err != nil {
		fmt.Fprintf(os.Stderr, "erro: url inválida: %v\n", err)
		os.Exit(2)
	}
	if *totalRequests <= 0 {
		fmt.Fprintln(os.Stderr, "erro: --requests deve ser maior que 0")
		os.Exit(2)
	}
	if *concurrency <= 0 {
		fmt.Fprintln(os.Stderr, "erro: --concurrency deve ser maior que 0")
		os.Exit(2)
	}
	if *concurrency > *totalRequests {
		*concurrency = *totalRequests
	}

	summary, err := runLoadTest(*serviceURL, *totalRequests, *concurrency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao executar teste de carga: %v\n", err)
		os.Exit(1)
	}

	printReport(summary)
}

func runLoadTest(targetURL string, total int, concurrency int) (resultsSummary, error) {
	start := time.Now()

	transport := &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: concurrency,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	jobs := make(chan int)
	results := make(chan int, total)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for range jobs {
			// Independent context per request to allow cancellation in future if needed
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, targetURL, nil)
			if err != nil {
				results <- -1
				continue
			}
			resp, err := client.Do(req)
			if err != nil {
				results <- -1
				continue
			}
			// Drain and close body to enable connection reuse
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			results <- resp.StatusCode
		}
	}

	// Start workers
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker()
	}

	// Enqueue jobs
	go func() {
		for i := 0; i < total; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Collect results
	statusHistogram := make(map[int]int)
	ok200Count := 0
	for i := 0; i < total; i++ {
		code := <-results
		if code == 200 {
			ok200Count++
			continue
		}
		statusHistogram[code] = statusHistogram[code] + 1
	}

	// Ensure all workers are done
	wg.Wait()

	return resultsSummary{
		totalRequests:   total,
		ok200Count:      ok200Count,
		statusHistogram: statusHistogram,
		duration:        time.Since(start),
	}, nil
}

func printReport(r resultsSummary) {
	fmt.Println("==== Relatório do Teste de Carga ====")
	fmt.Printf("Tempo total: %s\n", r.duration)
	fmt.Printf("Total de requests: %d\n", r.totalRequests)
	fmt.Printf("HTTP 200: %d\n", r.ok200Count)

	// Print other status codes in ascending order (grouping -1 as erros)
	if len(r.statusHistogram) == 0 {
		fmt.Println("Outros status: nenhum")
		return
	}
	codes := make([]int, 0, len(r.statusHistogram))
	for code := range r.statusHistogram {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	fmt.Println("Outros status:")
	for _, code := range codes {
		label := fmt.Sprintf("%d", code)
		if code == -1 {
			label = "erro"
		}
		fmt.Printf("  %s: %d\n", label, r.statusHistogram[code])
	}
}


