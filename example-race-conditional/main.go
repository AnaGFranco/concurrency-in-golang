package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	msg = "Hello, world!"

	//gera duas intancias do wg
	wg.Add(2)
	go updateMessage("Hello, universe!")
	go updateMessage("Hello, cosmos!")
	wg.Wait()

	fmt.Println(msg)
}

/* Como não temos uma ordem/tempo de execução com a
go routines, pode ser que o programa finalize e apenas
uma go routines termine a execussão.
Importante: se executado com go run -race, podemos ver
que a um problema de codição de corrida, pois utilizam
do mesmo recurso   */
