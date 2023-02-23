package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	mainCtx := context.Background()

	newCtx, _ := context.WithDeadline(mainCtx, time.Now().Add(time.Second*10))

	go firstOne(newCtx)

	<-mainCtx.Done()

	fmt.Println("Finished main")

}

func firstOne(ctx context.Context) {
	go secondOne(ctx)

	<-ctx.Done()
	fmt.Println("Finished first")
}

func secondOne(ctx context.Context) {
	go thirdOne(ctx)

	<-ctx.Done()
	fmt.Println("Finished second")
}

func thirdOne(ctx context.Context) {
	<-ctx.Done()
	fmt.Println("Finished third")
}
