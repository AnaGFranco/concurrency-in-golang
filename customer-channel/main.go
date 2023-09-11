package main

import (
	"fmt"
	"sync"
	"time"
)

type Purchase struct {
	Name     string
	Quantity int
}

var (
	stock     = 5
	stockLock sync.Mutex
	wg        sync.WaitGroup
)

func customer(name string, quantity int, purchaseCh chan<- Purchase, doneCh <-chan bool) {
	defer wg.Done()

	fmt.Printf("%s chegou e quer comprar %d unidades.\n", name, quantity)

	time.Sleep(1 * time.Second)

	purchase := Purchase{Name: name, Quantity: quantity}
	purchaseCh <- purchase // Envia pedido de compra para o canal

	confirmed := <-doneCh // Recebe confirmação de compra do canal
	if confirmed {
		fmt.Printf("%s comprou %d unidades. Estoque atual: %d\n", name, quantity, stock)
	} else {
		fmt.Printf("%s não pôde comprar %d unidades. Estoque insuficiente.\n", name, quantity)
	}
}

func stockManager(purchaseCh <-chan Purchase, doneCh chan<- bool) {
	for purchase := range purchaseCh {
		stockLock.Lock()
		if stock >= purchase.Quantity {
			stock -= purchase.Quantity
			doneCh <- true // Confirmação de compra bem-sucedida
		} else {
			doneCh <- false // Confirmação de compra sem estoque suficiente
		}
		stockLock.Unlock()
	}
	close(doneCh)
}

func main() {
	purchaseCh := make(chan Purchase) // Canal para pedidos de compra
	doneCh := make(chan bool)         // Canal para confirmação de compra

	go stockManager(purchaseCh, doneCh)

	wg.Add(3)
	go customer("Cliente A", 2, purchaseCh, doneCh)
	go customer("Cliente B", 3, purchaseCh, doneCh)
	go customer("Cliente C", 1, purchaseCh, doneCh)
	wg.Wait()

	close(purchaseCh)
	close(doneCh)
	fmt.Println("Estoque final:", stock)
}
