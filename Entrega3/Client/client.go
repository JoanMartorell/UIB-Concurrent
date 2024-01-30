package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

var (
	tresorerQueue = "Diposits"
	balancesQueue = "Balances"
)

func main() {

	// Obtenir el nom del client des de la línia de comanda
	if len(os.Args) < 2 {
		fmt.Println("Proporciona el nom del client com a argument.")
		return
	}
	clientName := os.Args[1]

	// Inicia el client
	client(clientName)
}

func client(name string) {
	// Inicia la connexió amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Crea el canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Informa del client i les operacions que vol fer
	operacions := rand.Intn(5) + 1

	fmt.Printf("Hola el meu nom és: %s\n", name)
	fmt.Printf("%s vol fer %d opracions\n", name, operacions)

	// Bucle de les operacions del client
	for i := 1; i <= operacions; i++ {

		// Decideix l'operació
		operacio := rand.Intn(15) - 10
		if rand.Intn(2) == 0 {
			operacio = operacio * -1
		}

		// Mostra l'operació
		fmt.Printf("%s operació %d: %d\n", name, i, operacio)

		// Publica el missatge a la cua Dipòsits
		publishOperation(ch, tresorerQueue, name, operacio)
		fmt.Println("Operació sol·licitada...")

		// Simula l'espera de la resposta
		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

		// Consulta el balanç
		msgs, err := ch.Consume(
			balancesQueue,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Processa els missatges de balanç
		for msg := range msgs {
			fmt.Printf("%s\n", string(msg.Body))
			balance, _ := strconv.Atoi(string(msg.Body))

			switch {
			case balance == 0:
				fmt.Println("EL TESORER SEMBLA SOSPITÓS")
			case balance < 0:
				fmt.Println("NO HI HA SALDO AL COMPTE")
			case operacio >= 0:
				fmt.Println("INGRÉS CORRECTE")
			case operacio < 0 && balance+operacio >= 0:
				fmt.Println("ES FARÀ EL REINTEGRE SI HI HA SALDO")
			}

			// Mostra el balanç i el resultat de l'operació
			fmt.Printf("Balanç actual: %d\n", balance)

			break
		}

		fmt.Printf("%d---------------------------\n", i)
	}

}

func publishOperation(ch *amqp.Channel, tresorerQueueName, clientName string, operacio int) {
	// Converteix l'operació a una cadena
	operacioStr := strconv.Itoa(operacio)

	// Publica l'operació a la cua Dipòsits
	err := ch.Publish(
		"",
		tresorerQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(clientName + " " + operacioStr),
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
