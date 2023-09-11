package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	stock = 5
	wg    sync.WaitGroup
)

func customer(name string, quantity int) {
	defer wg.Done()

	fmt.Printf("%s chegou e quer comprar/devolver %d unidades.\n", name, quantity)

	// Simulação de processamento
	time.Sleep(3 * time.Second)

	if stock >= quantity {
		stock -= quantity
		fmt.Printf("%s comprou %d unidades e o estoque atualizou para: %d\n", name, quantity, stock)
	} else {
		fmt.Printf("%s não comprou \n", name)
	}

}

func main() {
	wg.Add(3)
	go customer("Cliente A", 2) // Compra de 2 unidades
	go customer("Cliente B", 3) // Compra de 3 unidades
	go customer("Cliente C", 1) // Compra de 1 unidade
	wg.Wait()

	fmt.Println("Estoque final:", stock)
}
