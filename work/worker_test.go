package work_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/justsocialapps/assert"
	"github.com/justsocialapps/justlib/work"
)

func Example() {

	// create a worker with 3 goroutines that can process jobs in parallel
	worker := work.NewWorker(3, func(p work.Payload) interface{} {
		// this is our worker function that is called for every new job
		fmt.Println(p.Data)
		time.Sleep(200 * time.Millisecond)
		return nil
	}, false)

	// create 100 jobs and dispatch them to the worker.
	for i := 0; i < 100; i++ {
		// this call will block when all 3 goroutines are currently busy.
		worker.Dispatch(work.Payload{Data: strconv.Itoa(i)})
	}

	// this call makes sure that the worker stops all goroutines as soon as
	// they have processed all remaining jobs.
	worker.Quit()

}

func TestWorkerShouldWorkSequentiallyWithOnlyOneGoroutine(t *testing.T) {
	var result string

	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		return fmt.Sprintf("%s.", p.Data)
	}, true)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range worker.Completions() {
			result += v.Output.(string)
		}
	}()

	// create 100 jobs and dispatch them to the worker.
	for i := 0; i < 100; i++ {
		worker.Dispatch(work.Payload{Data: strconv.Itoa(i)})
	}

	// this call makes sure that the worker stops all goroutines as soon as
	// they have processed all remaining jobs.
	worker.Quit()
	wg.Wait()

	a := assert.NewAssert(t)
	a.Equal(result, "0.1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17.18.19.20.21.22.23.24.25.26.27.28.29.30.31.32.33.34.35.36.37.38.39.40.41.42.43.44.45.46.47.48.49.50.51.52.53.54.55.56.57.58.59.60.61.62.63.64.65.66.67.68.69.70.71.72.73.74.75.76.77.78.79.80.81.82.83.84.85.86.87.88.89.90.91.92.93.94.95.96.97.98.99.", "Jobs were completed in wrong order or incompletely")

}