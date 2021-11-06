package golangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println("background:", background)

	todo := context.TODO()
	fmt.Println("todo:", todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA) // context.Background
	fmt.Println(contextB) // context.Background.WithValue(type string, val B)
	fmt.Println(contextC) // context.Background.WithValue(type string, val C)
	fmt.Println(contextD) // context.Background.WithValue(type string, val B).WithValue(type string, val D)
	fmt.Println(contextE) // context.Background.WithValue(type string, val B).WithValue(type string, val E)
	fmt.Println(contextF) // context.Background.WithValue(type string, val C).WithValue(type string, val F)
}

func TestContextGetValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA.Value("b")) // nil
	fmt.Println(contextB.Value("a")) // nil
	fmt.Println(contextC.Value("a")) // nil
	fmt.Println(contextD.Value("b")) // B
	fmt.Println(contextE.Value("b")) // B
	fmt.Println(contextF.Value("b")) // nil
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	// goroutine leak
	// go func() {
	// 	defer close(destination)

	// 	for counter := 1; ; counter++ {
	// 		destination <- counter
	// 	}
	// }()

	// no goroutine leak
	go func() {
		defer close(destination)

		for counter := 1; ; counter++ {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
		// Counter:  6
		// Counter:  7
		// Counter:  8
		// Counter:  9
		// Counter: 10

		if n == 10 {
			break
		}
	}

	cancel()

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
	// When Goroutine Leak: Total Goroutine 3
	// When Goroutine No Leak: Total Goroutine 2
}

func CreateSlowCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)

		for counter := 1; ; counter++ {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				time.Sleep(time.Second)
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2
}
