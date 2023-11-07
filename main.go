package main

import (
	"WildBerries/internal/WildBerries"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

var wb *WildBerries.Parser

func main() {
	wb = WildBerries.NewParser()
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	fmt.Println("Сервер запущен по адресу http://127.0.0.1:4445")
	go func() {
		log.Fatal("Ошибка: ", http.ListenAndServe(":4445", mux))
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New("6917956424:AAE0CG0KfslaxAQt3QwnOINHYOSxzqwGtj8", opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if strings.Contains(update.Message.Text, "wildberries.ru") || strings.Contains(update.Message.Text, "wb.ru") {
		if strings.Contains(update.Message.Text, "search.wb.ru") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Обнаружена ссылка на поиск по категории, ищу список доступных объявлений"),
			})
			artCount := wb.GetWBAdverts(update.Message.Text)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Найдено %d товаров, они добавлены в базу для последующего обхода", artCount),
			})
		} else if strings.Contains(update.Message.Text, "detail.aspx") {
			articleS := strings.ReplaceAll(update.Message.Text, "https://www.wildberries.ru/catalog/", "")
			articleS = strings.ReplaceAll(articleS, "/detail.aspx", "")
			article, err := strconv.Atoi(articleS)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Не корректная ссылка"),
				})
				return
			}

			wb.AddArticles([]string{strconv.Itoa(article)})
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Артикул %d добавлен\nТоваров в работе: %d", article, wb.GetArticlesCount()),
			})
		}

	}

}
