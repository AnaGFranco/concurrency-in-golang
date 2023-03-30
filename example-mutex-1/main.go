package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

/* Aqui temos um problema de concorrencia de recurso,
pois temos duas chamadas de go routines que utilizam
do mesmo nivel da variavel msg.
Com o Mutex ele consegue bloquear o recurso  para que
ninguem mais use até que a execução finalize*/

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = s
	m.Unlock()
}

func main() {
	msg = "Hello, world!"

	var mutex sync.Mutex

	wg.Add(2)
	go updateMessage("Hello, universe!", &mutex)
	go updateMessage("Hello, cosmos!", &mutex)
	wg.Wait()

	fmt.Println(msg)
}

/* Importante: com o uso do mutex não dara mais data race,
mas não conseguimos garantir a qual será a saida ao
executar go run -race*/