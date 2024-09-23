package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var store1 = 10
var store2 = 11
var store3 = 12

var product1 = 10
var product2 = 11
var product3 = 12
var product4 = 13
var product5 = 14
var product6 = 15
var product7 = 16
var product8 = 17
var product9 = 18

func main() {
	var wg sync.WaitGroup
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store1, product1)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store1, product2)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store1, product3)
	}()
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store1, 66)
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store2, product4)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store2, product5)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store2, product6)
	}()
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store2, 66)
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store3, product7)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store3, product8)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store3, product9)
	}()
	go func ()  {
		defer wg.Done()
		runGettingJob(ctx, store3, 66)
	}()

	// deleting
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store1, product1)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store1, product2)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store1, product3)
	}()
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store1, 66)
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store2, product4)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store2, product5)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store2, product6)
	}()
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store2, 66)
	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store3, product7)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store3, product8)
	}()
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store3, product9)
	}()
	go func ()  {
		defer wg.Done()
		runDeletingItemJob(ctx, store3, 66)
	}()

	wg.Wait()
}

func runGettingJob(ctx context.Context, storeId, productId int) {
	url := fmt.Sprintf("http://arch.homework/store-service/api/v1/stores/%d/products/%d", storeId, productId)
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		count++
		_, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(5 * time.Second)
	}
}

func runDeletingItemJob(ctx context.Context, storeId, productId int) {
	client := &http.Client{}
	url := fmt.Sprintf("http://arch.homework/store-service/api/v1/stores/%d/products/item/%d", storeId, productId)
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		count++

		// create a new DELETE request
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalln(err)
		}

		// send the request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		time.Sleep(5 * time.Second)
	}
}
