package main

import (
	"context"
	"fmt"

	"github.com/fabiuhp/orders-api/app"
)

func main() {
	app := app.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Erro ao iniciar: ", err)
	}
}
