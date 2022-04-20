package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func InvokeAPI(wg *sync.WaitGroup, queue chan int) []string {
	defer wg.Done()
	resp, err := http.Get("http://localhost:4000/athleteSchedule/findScheduleByName?name=Green&region=US")
	for job := range queue {
		csvWriter([]string{strconv.Itoa(job), strconv.Itoa(resp.StatusCode), http.StatusText(resp.StatusCode)})
	}
	if resp == nil {
		return []string{"nil", "nil"}
	}
	if err != nil {
		fmt.Println("Error:", err)
	}

	res := []string{strconv.Itoa(resp.StatusCode), http.StatusText(resp.StatusCode)}
	return res
}

func csvWriter(yourSliceGoesHere []string) {
	f, err := os.OpenFile("perf-test.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	w := csv.NewWriter(f)
	w.Write(yourSliceGoesHere)
	w.Flush()
	f.Close()
}

func main() {
	start := time.Now()
	const max = 200
	queue := make(chan int, max)
	var wg sync.WaitGroup
	for i := 0; i < max; i++ {
		wg.Add(1)
		go InvokeAPI(&wg, queue)
	}
	for i := 0; i < 500000; i++ {
		queue <- i
	}
	close(queue)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Print("Total elapsed time: ", elapsed)
}
