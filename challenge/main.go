package main

import (
	"fmt"
	"sync"
)

var msg string

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

var wg sync.WaitGroup

func main() {

	/*desafio: modifique este código para que as chamadas para updateMessage() nas linhas
	27, 30 e 33 são executados como goroutines e implementam grupos de espera para que
	o programa é executado corretamente e imprime três mensagens diferentes.
	Em seguida, escreva um teste para todas as três funções neste programa: updateMessage(),
	printMessage() e main().*/

	msg = "Hello, world!"

	wg.Add(1)
	go updateMessage("Hello, universe!")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello, cosmos!")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello, world!")
	wg.Wait()
	printMessage()
}
