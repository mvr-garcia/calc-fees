package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

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

	for msg := range out {
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
		fmt.Println(outputDTO)
	}
}
