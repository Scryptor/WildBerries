package main

import (
	"WildBerries/internal/WildBerries"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var wb map[int64]*WildBerries.Parser

func main() {
	wb = make(map[int64]*WildBerries.Parser)
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	fmt.Println("Сервер запущен по адресу http://127.0.0.1:4445")
	go func() {
		log.Fatal("Ошибка: ", http.ListenAndServe(":4445", mux))
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(tgHandler),
	}

	b, err := bot.New("6917956424:AAE0CG0KfslaxAQt3QwnOINHYOSxzqwGtj8", opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
