package main

import (
	"fmt"
)

func printSomething(s string) {
	fmt.Println(s)
}

func main() {
	/*  Com a go routines programa (provavelmente) apenas imprimirá a segunda mensagem,
	pois o programa termina a execução antes da go routines iniciar. Obs: O comando go (goroutine)
	não tem hora exata de iniciar */

	//go printSomething("This is the first thing to be printed!")

	printSomething("This is the first thing to be printed!")

	/* Para dar tempo para a goroutine finalizar, não é uma boa solução, mas pode ser
	adicionado um tempo de espera para que de tempo da execução */
	//time.Sleep(1 * time.Second)

	printSomething("This is the second thing to be printed!")
}
