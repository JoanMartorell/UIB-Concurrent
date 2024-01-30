package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

const MinimBoti = 20

var (
	wg            sync.WaitGroup
	tresorerQueue = "Diposits"
	balancesQueue = "Balances"
	stopOffice    = make(chan bool)
)

func main() {

	// Inicia el tresorer
	go tresorer()

	// Espera a un senyal per acabar
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Tanca l'oficina
	stopOffice <- true
	wg.Wait()
}

func tresorer() {
	// Inicia la connexió amb RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Crea els canals
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Declara les cues
	dipositsQueue, err := ch.QueueDeclare(
		tresorerQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	balancesQueue, err := ch.QueueDeclare(
		balancesQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Informa del balanç inicial
	balance := 0
	fmt.Printf("El tresorer és al despatx. El botí mínim és: %d\n", MinimBoti)

	// Bucle principal del tresorer
	for {
		select {
		case <-stopOffice:
			// Tanca l'oficina
			fmt.Printf("El Tresorer s'en va\n\n")
			fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			closeOffice(ch, dipositsQueue.Name, balancesQueue.Name)
			return
		default:
			// Rep missatges de la cua Dipòsits
			msgs, err := ch.Consume(
				dipositsQueue.Name,
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

			for msg := range msgs {

				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

				// Processa els missatges
				valors := strings.Split(string(msg.Body), " ")
				clientName := valors[0]
				amount, _ := strconv.Atoi(valors[1])
				fmt.Printf("Operació rebuda: %d del client: %s\n", amount, clientName)

				// Incrementa o decrementa el balanç segons l'operació
				if amount >= 0 {
					balance += amount
				} else if amount < 0 && balance+amount >= 0 {
					balance += amount
				} else {
					fmt.Println("OPERACIÓ NO PERMESA NO HI HA FONS")
				}

				// Publica el balanç a la cua Balances
				publishBalance(ch, balancesQueue.Name, balance)

				fmt.Printf("Balanç: %d\n", balance)

				// Comprova si es pot robar el dipòsit i tancar l'oficina
				if balance >= MinimBoti {
					fmt.Println("El Tresorer decideix robar el dipòsit i tancar el despatx")
					closeOffice(ch, dipositsQueue.Name, balancesQueue.Name)
					return
				}

				fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			}
		}
	}
}

func closeOffice(ch *amqp.Channel, dipositsQueueName, balancesQueueName string) {
	fmt.Printf("Cua '%s' esborrada amb èxit.\n", dipositsQueueName)
	ch.QueueDelete(dipositsQueueName, false, false, false)

	fmt.Printf("Cua '%s' esborrada amb èxit.\n", balancesQueueName)
	ch.QueueDelete(balancesQueueName, false, false, false)
}

func publishBalance(ch *amqp.Channel, balancesQueueName string, balance int) {
	// Converteix el balanç a una cadena
	balanceStr := strconv.Itoa(balance)

	// Publica el balanç a la cua de Balances
	err := ch.Publish(
		"",
		balancesQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(balanceStr),
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
