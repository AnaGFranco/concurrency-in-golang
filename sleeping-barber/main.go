package main

/* Exemplo do barbeiro adormecido.
Podemos ter um número finito de barbeiros, um número finito de assentos em uma
sala de espera, um período fixo de tempo em que a barbearia fica aberto e clientes
chegando em intervalos (aproximadamente) regulares.

Regras:

- se não houver clientes, o barbeiro adormece na cadeira
- um cliente deve acordar o barbeiro se ele estiver dormindo
- se um cliente chega enquanto o barbeiro está trabalhando, o cliente sai se todas as cadeiras estiverem ocupadas e
senta-se em uma cadeira vazia se estiver disponível
- quando o barbeiro termina de cortar o cabelo, ele inspeciona a sala de espera para ver se há clientes esperando
e adormece se não houver nenhum
- a loja pode deixar de aceitar novos clientes na hora de fechar, mas os barbeiros não podem sair até que a sala de espera esteja
vazio
- depois que a loja fechar e não houver clientes na área de espera, o barbeiro
ir para casa
*/

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arrivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	//gera numeros aleatorios
	rand.Seed(time.Now().UnixNano())

	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day!")

	shop.addBarber("Frank")

	// comece a barbearia como uma goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 1

	go func() {
		for {
			//obter um número aleatório com taxa média de chegada
			randomMillseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillseconds)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	//bloquear até a barbearia fechar
	<-closed
}
