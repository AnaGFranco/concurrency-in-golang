package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	// variável para saldo bancário
	var bankBalance int
	var balance sync.Mutex

	// imprimir valores iniciais
	fmt.Printf("Initial account balance: $%d.00", bankBalance)
	fmt.Println()

	// definir receita semanal
	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Gifts", Amount: 10},
		{Source: "Part time job", Amount: 50},
		{Source: "Investments", Amount: 100},
	}

	wg.Add(len(incomes))

	// faça um loop por 52 semanas e imprima quanto é feito; manter um total em execução
	for i, income := range incomes {

		go func(i int, income Income) {
			defer wg.Done()

			for week := 1; week <= 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()

				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}
		}(i, income)
	}

	wg.Wait()

	// imprimir saldo final
	fmt.Printf("Final bank balance: $%d.00", bankBalance)
	fmt.Println()
}

/*
	RUN:  go run -race .
	IMPORTANTE: execução simultanea nao garante a ordem da execução
	O Mutex é importante para nao dar race conditional, pois bloqueia
	o recurso enquanto ele está sendo utilizado
*/
