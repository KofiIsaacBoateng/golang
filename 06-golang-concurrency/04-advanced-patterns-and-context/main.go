package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type SharedFacts struct {
	mu sync.Mutex

	Facts []string

	TotalBytes int
}

func NewSharedFacts () *SharedFacts {
	return &SharedFacts{
		Facts: make([]string, 0),
		TotalBytes: 0,
	}
}


func (facts *SharedFacts) AddFacts(fact string, totalBytes int){
	facts.mu.Lock()
	defer facts.mu.Unlock()


	facts.Facts = append(facts.Facts, fact)
	facts.TotalBytes += totalBytes
}


func worker(workerId int, taskQueue <-chan int, sharedFacts *SharedFacts, limiter <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initialize a reusable HTTP Client for connection efficiency
	client := &http.Client{};

	for task := range taskQueue{
		<-limiter; 	// Wait for the rate-limiter channel ticker token to arrive

		startTime := time.Now()
		fmt.Printf("[Worker %d] Dispatching Request %d at %s\n", workerId, task, startTime.Format("15:04:05.000"))
		


		// Enforce a strict 2-second timeout SLA for this network call using context
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		fact, bytesRead, err := fetchCatFact(ctx, client)
		cancel() // Destroy the context timer resource immediately after the network call finishes

		if err != nil {
			fmt.Printf("\t[Worker %d] ! HTTP ERROR on Request %d: %v\n", workerId, task, err)
			continue
		}

		// Save results into our thread-safe shared stats object
		sharedFacts.AddFacts(fact, bytesRead)
		
		fmt.Printf("\t[Worker %d] ✓ Request %d Complete in %v\n", workerId, task, time.Since(startTime))
		
	}

}

func fetchCatFact(ctx context.Context, client *http.Client) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://catfact.ninja/fact", nil);
	if(err != nil) {
		return "", 0, err;
	}

	res, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close();

	var fact CatFact
	if err := json.NewDecoder(res.Body).Decode(&fact); err != nil {
		return "", 0, err
	}

	return fact.Fact, fact.Length, nil
}


func main() {
	taskCount := 6;
	taskQueue := make(chan int, taskCount)
	sharedFacts := NewSharedFacts(); // state manager instance


	// Rate limiter - throttle to two tasks per second (1 tick every 500ms)
	limiter := time.Tick(500*time.Millisecond);

	var wg sync.WaitGroup;

	// Spawn a fixed pool of 3 workers 
	for i := range 3{
		wg.Add(1)
		go worker(i + 1, taskQueue, sharedFacts, limiter, &wg)
	}


	// Dispatch tasks to buffered task channel for workers to operate
	for task := range taskCount {
		taskQueue <- task + 1
	}

	close(taskQueue) // Close channel so workers know no more tasks are coming


	// Wait for all background worker routines to drain the queue and exit
	wg.Wait()


	// Safely output the compiled final state metrics
	fmt.Printf("\n--- Scraping Operation Complete ---\n")
	fmt.Printf("Total Payload Data Collected: %d bytes\n", sharedFacts.TotalBytes)
	fmt.Println("Captured Facts:")
	for idx, fact := range sharedFacts.Facts {
		fmt.Printf("  %d. %s\n", idx+1, fact)
	}

}