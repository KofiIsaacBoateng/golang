package main

import (
	"fmt"
	"time"
)

// Simulates produce prices generation...
func produducePrices(priceChannel chan<- float64) {
	prices := []float64{5.34, 234.34, 231.12, 88.32, 232.11, 11.00};

	for _, price := range prices {
		fmt.Printf("[Producer] Generating price tickets $%.2f\n", price)
		priceChannel <- price; // sent price to channel
		time.Sleep(150 * time.Millisecond)// simulate network latency
	}

	fmt.Println("[Producer] Generation complete! Closing channel...!")
	close(priceChannel);
}

// process incoming data from stream
func consumePrices(priceChannel <-chan float64, doneChannel chan<- bool) {

	for price := range priceChannel{ // loops automatically until channel is drained
		price = price * 1.02 // applying 2% tax on every price index
		fmt.Printf("\t[Consumer] Processed price (with fee): $%.2f\n", price)
	}

	fmt.Println("[Consumer] Channel drained completely!")
	doneChannel <- true; // signal main that processing is done!
}

func main() {
	// create unbuffered channels
	priceChannel := make(chan float64)
	doneChannel := make(chan bool)

	fmt.Println("Pipleline operation initiating...")

	go produducePrices(priceChannel)
	go consumePrices(priceChannel, doneChannel)

	// global timeout channel to simulate a network timeout
	timeout := time.After(2 * time.Second)

	select {
	case <-doneChannel:
		fmt.Println("Pipeline operation executed successfully!")
	
	case <-timeout:
		fmt.Println("CRITICAL ERROR: Operation timed out!") 
	}
}