package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mvr-garcia/calc-fees/internal/order/infra/database"
	usecase "github.com/mvr-garcia/calc-fees/internal/order/use_case"
	"github.com/mvr-garcia/calc-fees/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "./orders.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := database.NewOrderRepository(db)
	uc := usecase.CalculateFinalPriceUseCase{OrderRepository: repo}

	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	out := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, out)

	qtdWorkers := 30
	for i := 1; i <= qtdWorkers; i++ {
		go worker(out, &uc, i)
	}

	http.HandleFunc(
		"/total",
		func(w http.ResponseWriter, r *http.Request) {
			uc := usecase.GetTotalUseCase{OrderRepository: repo}

			total, err := uc.Execute()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			json.NewEncoder(w).Encode(total)
		},
	)

	http.ListenAndServe(":8080", nil)
}

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workerID int) {
	for msg := range deliveryMessage {
		var inputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &inputDTO)
		if err != nil {
			panic(err)
		}
		outputDTO, err := uc.Execute(inputDTO)
		if err != nil {
			panic(outputDTO)
		}
		msg.Ack(false)
		fmt.Printf("Worker ID: %d has processed Order ID: %s\n", workerID, outputDTO.ID)
		time.Sleep(1 * time.Second)
	}
}
