package main

import (
	"WildBerries/internal/WildBerries"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tgid := r.URL.Query().Get("tgid")
	hash := r.URL.Query().Get("hash")

	if tgid == "" || hash == "" {
		w.Write([]byte("Нет данных для входа в личный кабинет"))
		return
	}

	tgInt, err := strconv.ParseInt(tgid, 10, 64)
	if err != nil {
		w.Write([]byte("Не корректный tgId"))
		return
	}

	wbw, ok := wb[tgInt]
	if !ok {
		w.Write([]byte("Нет такого пользователя"))
		return
	} else {
		if hash != wbw.Hash {
			w.Write([]byte("Не корректные данные для входа"))
			return
		}
	}

	t, err := template.New("articles").Parse(articlesTemplate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	err = t.Execute(w, wbw.GetAllArticles())
	if err != nil {
		log.Println(err)
		return
	}

}
func tgHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if wbWorker, ok := wb[update.Message.From.ID]; !ok {
		wb[update.Message.From.ID] = WildBerries.NewParser(update.Message.From.ID, update.Message.Chat.ID)
		wwb := wb[update.Message.From.ID]
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: fmt.Sprintf("Создан новый пользователь, добавляйте ссылки на товары сюда, в телеграм"+
				"\nОтредактировать товары после добавления, можно в личном кабинете: \n"+
				"%s", wwb.GetLKLink()),
		})
	} else {
		if strings.Contains(update.Message.Text, "wildberries.ru") || strings.Contains(update.Message.Text, "wb.ru") {
			if strings.Contains(update.Message.Text, "search.wb.ru") {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Обнаружена ссылка на поиск по категории, ищу список доступных объявлений"),
				})
				artCount := wbWorker.GetWBAdverts(update.Message.Text)
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

				wbWorker.AddArticles([]string{strconv.Itoa(article)})
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Артикул %d добавлен\nТоваров в работе: %d", article, wbWorker.GetArticlesCount()),
				})
			}

		}
	}

}
