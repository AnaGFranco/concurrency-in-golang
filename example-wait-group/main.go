package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done() // diminuir wg em um após a conclusão da função

	fmt.Println(s)
}

func main() {
	var wg sync.WaitGroup

	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"pi",
		"zeta",
		"eta",
		"theta",
		"epsilon",
	}

	wg.Add(len(words))

	for i, x := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, x), &wg)
	}

	// o programa irá pausar neste ponto, até que wg seja 0
	wg.Wait()

	// adicionando wg para não dar erro na chamada desta função
	wg.Add(1)
	printSomething("This is the second thing to be printed!", &wg)
}

/* IMPORTANTE: WaitGroup não garante a ordena */
