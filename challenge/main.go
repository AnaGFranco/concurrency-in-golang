package main

import (
	"fmt"
)

var msg string

func updateMessage(s string) {
	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

func main() {

	/*desafio: modifique este código para que as chamadas para updateMessage() nas linhas
	27, 30 e 33 são executados como goroutines e implementam grupos de espera para que
	o programa é executado corretamente e imprime três mensagens diferentes.
	Em seguida, escreva um teste para todas as três funções neste programa: updateMessage(),
	printMessage() e main().*/

	msg = "Hello, world!"

	updateMessage("Hello, universe!")
	printMessage()

	updateMessage("Hello, cosmos!")
	printMessage()

	updateMessage("Hello, world!")
	printMessage()
}
