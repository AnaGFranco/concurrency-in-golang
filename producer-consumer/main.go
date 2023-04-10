package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

/* Producer é um tipo de struct que contém dois canais: umara pizzas,
com todos informações para um determinado pedido de pizza, incluindo
se foi feito com sucesso, e outro para lidar com o fim do processamento
(quando saímos do channel)*/

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

/* PizzaOrder é um tipo de struct que descreve um determinado pedido
de pizza. Tem a ordem número, uma mensagem indicando o que aconteceu
com o pedido e um booleano indicando se o pedido foi concluído com
sucesso.*/

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

/* Fechar é simplesmente um método de fechar o channel quando terminamos
com ele (ou seja, algo é enviado para o channel de saída) */

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

/*makePizza tenta fazer uma pizza. Geramos um número aleatório de 1 a 12,
e colocar em dois casos em que não podemos fazer a pizza a tempo. De outra
forma, fazemos a pizza sem problemas. Para tornar as coisas interessantes,
cada pizzalevará um tempo diferente para produzir (algumas pizzas são mais
duras do que outras).*/

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds....\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

/* pizzaria é uma goroutine que roda em segundo plano e chama makePizza para
tentar fazer um pedido cada vez que itera o loop for. Ele executa até receber
algo ao sair do channel. O channel de saída não recebe nada até que o consumidor
envia (quando o número de pedidos é maior ou igual ao NumberOfPizzas constante). */

func pizzeria(pizzaMaker *Producer) {
	// acompanhar qual pizza estamos fazendo
	var i = 0

	/* executar para sempre ou até receber uma notificação de encerramento
	tentar fazer pizzas */
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// tentamos fazer uma pizza (enviamos algo para os dados channel -- a chan PizzaOrder)
			case pizzaMaker.data <- *currentPizza:

			// queremos sair, então envie pizzMaker.quit para o quitChan (um erro chan)
			case quitChan := <-pizzaMaker.quit:
				// fechar channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {
	// gerador de números aleatórios
	rand.Seed(time.Now().UnixNano())

	// imprimir uma mensagem colorida
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	//  criar producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// executar o produtor em segundo plano
	go pizzeria(pizzaJob)

	// criar e executar consumidor
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	// imprimir a mensagem final
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day....")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day!")
	default:
		color.Green("It was a great day!")
	}
}
