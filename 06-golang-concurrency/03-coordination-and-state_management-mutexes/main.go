package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCache struct {
	mu sync.Mutex

	Cache map[string]int
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		Cache: make(map[string]int),
	}
}

// Set safely updates the cache using an exclusive Lock
func (sc *SafeCache) Set(key string, value int) {
	sc.mu.Lock() // Stop all other goroutines from accessing the map
	defer sc.mu.Unlock() // Crucial: Unlock when the function completes

	sc.Cache[key] = value;
}

// Get safely reads from the cache using a Lock
func (sc *SafeCache) Get(key string) (int, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	value, exists := (sc.Cache[key]);

	return value, exists
}

func worker(workerId int, store *SafeCache, wg *sync.WaitGroup) {
	wg.Add(1);
	defer wg.Done();

	// make each worker change the value of the same keys
	for i := range 3{
		key := fmt.Sprintf("key_%d", i + 1);
		store.Set(key, workerId * 10 + i);
		time.Sleep(10 * time.Millisecond); // Simulate network latency
	}
}


func main() {
	store := NewSafeCache();

	// initialize a wait group
	var wg sync.WaitGroup;

	fmt.Println("Launching 5 concurrent cache writers...")

	// Spawn 5 concurrent worker groups to update the store
	for i := range 5 {
		go worker(i + 1, store, &wg);
	}

	wg.Wait(); // Block main until all 5 workers call wg.Done();

	fmt.Println("\nWriting complete! Reading final values securely:")
	for i := range 3 {
		key := fmt.Sprintf("key_%d", i + 1)
		val, _ := store.Get(key)
		fmt.Printf(" -> Cache Key: %s | Settled Value: %d\n", key, val)
	}
}