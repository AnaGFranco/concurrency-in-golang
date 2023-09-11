package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

/* O problema Dining Philosophers é bem conhecido nos círculos da ciência da computação.
Cinco filósofos, numerados de 0 a 4, moram em uma casa onde a mesa está posta para eles;
cada filósofo tem seu próprio lugar à mesa. Sua única dificuldade – além da filosofia –
é que o prato servido é um tipo de espaguete muito difícil que deve ser comido com dois
garfos. Existem dois garfos ao lado de cada prato, de modo que não apresenta dificuldade.
Como consequência, no entanto, isso significa que não há dois vizinhos podem estar comendo
simultaneamente, pois são cinco filósofos e cinco garfos.

Esta é uma implementação simples da solução da Dijkstra para o "Dining
O dilema dos filósofos. */

// Philosopher é uma estrutura que armazena informações sobre um filósofo.
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// philosophers é uma lista de todos filosofos.
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

var hunger = 3                  // quantas vezes um filósofo come
var eatTime = 1 * time.Second   // quanto tempo leva para comer
var thinkTime = 3 * time.Second // quanto tempo um filósofo pensa
var sleepTime = 1 * time.Second // quanto tempo esperar ao imprimir as coisas

var orderMutex sync.Mutex
var orderFinished []string // a ordem em que os filósofos terminam de jantar e saem

func main() {
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("The table is empty.")

	time.Sleep(sleepTime)

	// começar a refeição
	dine()

	fmt.Println("The table is empty.")

	time.Sleep(sleepTime)
	fmt.Printf("Order finished: %s.\n", strings.Join(orderFinished, ", "))

}

func dine() {
	// eatTime = 0 * time.Second
	// sleepTime = 0 * time.Second
	// thinkTime = 0 * time.Second

	/*wg é o WaitGroup que controla quantos filósofos ainda estão na mesa. Quando
	chega a zero, todos terminaram de comer e foram embora. Adicionamos 5 (o número de filósofos) ao wait group.*/
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// Queremos que todos estejam sentados antes de começarem a comer, então crie um WaitGroup para isso e defina-o como 5
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	/* forks é um mapa de todos os 5 forks. As bifurcações são atribuídas usando os campos leftFork e rightFork no Philosopher
	tipo. Cada bifurcação, então, pode ser encontrada usando o índice (um inteiro), e cada bifurcação tem um mutex único.*/
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// Comece a refeição repetindo nossa fatia de filosofos.
	for i := 0; i < len(philosophers); i++ {
		// dispare uma goroutine para o filósofo atual
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	//Espere os filósofos terminarem. Isso bloqueia até que o grupo de espera seja 0.
	wg.Wait()
}

/*
	DiningProblem é a função disparada como uma goroutine para cada um de nossos filósofos. é preciso um

filósofo, nosso WaitGroup para determinar quando todos terminam, um mapa contendo os mutexes para cada
fork na mesa e um WaitGroup usado para pausar a execução de cada instância desta goroutine
até que todos estejam sentados à mesa.
*/
func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// sentar o filósofo à mesa
	fmt.Printf("%s is seated at the table.\n", philosopher.name)

	// Diminuir o WaitGroup sentado em um.
	seated.Done()

	// Espere até que todos estejam sentados.
	seated.Wait()

	// Faça com que este filósofo comaTempo e penseTempo em tempos de "fome" (3).
	for i := hunger; i > 0; i-- {
		/* Obtenha um bloqueio nos garfos esquerdo e direito. Temos que escolher primeiro a bifurcação numerada inferior para
		para evitar uma condição de corrida lógica, que não é detectada pelo sinalizador -race nos testes; se não fizermos isso,
		temos potencial para um impasse, já que dois filósofos esperarão indefinidamente pela mesma bifurcação.
		Observe que a goroutine irá bloquear (pausar) até obter um bloqueio nas bifurcações direita e esquerda.*/
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
		}

		// No momento em que chegamos a esta linha, o filósofo tem um bloqueio (mutex) em ambos os garfos.
		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		// O filósofo começa a pensar, mas ainda não larga os garfos.
		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		// Desbloqueie os mutexes para ambos os garfos.
		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}

	// O filósofo terminou de comer, então imprima uma mensagem.
	fmt.Println(philosopher.name, "is satisified.")
	fmt.Println(philosopher.name, "left the table.")

	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	orderMutex.Unlock()
}
