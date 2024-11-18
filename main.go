package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fabiuhp/orders-api/app"
)

func main() {
	app := app.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("Erro ao iniciar: ", err)
	}
}
