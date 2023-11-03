package WildBerries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// AddArticle Добавляет товары по одному
func (wb *Parser) AddArticle(art string) {
	wb.getFullAdvert([]string{art})
}

// AddArticles Добавляет пачку товаров
func (wb *Parser) AddArticles(arts []string) {
	wb.getFullAdvert(arts)
}

// updateArticles Обновляет данные по сохраненным товарам
func (wb *Parser) updateArticles() {
	t := time.NewTicker(30 * time.Minute)
	ctx, _ := context.WithCancel(wb.ctx)
	for {
		select {
		case <-t.C:
			wb.articlesMutex.Lock()

			ArtsByParts400 := make([][]string, 0, len(wb.articles)/400+1)
			counter := 0
			arrIndex := 0
			for key := range wb.articles {
				if counter >= 400 {
					counter = 0
					arrIndex++
					ArtsByParts400 = append(ArtsByParts400, []string{})
				}
				ArtsByParts400[arrIndex] = append(ArtsByParts400[0], strconv.Itoa(key))
				counter++
			}

			wb.articlesMutex.Unlock()
			for _, articles400 := range ArtsByParts400 {
				wb.updateFullAdvert(articles400)
			}
		case <-ctx.Done():
			log.Println("Останавливаю обновление товаров")
			return
		}
	}

}

// Получает товары по артикулам. Не более 400 артикулов за раз
func (wb *Parser) getFullAdvert(articles []string) {
	artString := strings.Join(articles, ";")
	link := fmt.Sprintf("https://card.wb.ru/cards/v1/detail?appType=1&curr=rub&spp=30&nm=%s", artString)
	response, err := wb.GetDataFromWb(link)
	if err != nil {
		return
	}
	var Arts FullAdvertJson
	err = json.Unmarshal([]byte(response.Body), &Arts)
	if err != nil {
		wb.allErrors = append(wb.allErrors, fmt.Sprintf("Json Unmarshal Failed: %s", err.Error()))
		wb.metrics.BadJson.Add(1)
		return
	}
	wb.metrics.GoodReq.Add(1)
	for _, val := range Arts.Data.Products {
		adv := Advert{
			Id:          val.Id,
			Link:        fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", val.Id),
			ParsingDate: time.Now(),
			Price:       val.PriceU,
			PriceSale:   val.SalePriceU,
			Name:        val.Name,
			Brand:       val.Brand,
		}
		wb.addMapArticle(adv)
	}
}

// Обновляем товары. Не более 400 артикулов за раз
func (wb *Parser) updateFullAdvert(articles []string) {
	artString := strings.Join(articles, ";")
	link := fmt.Sprintf("https://card.wb.ru/cards/v1/detail?appType=1&curr=rub&spp=30&nm=%s", artString)
	response, err := wb.GetDataFromWb(link)
	if err != nil {
		return
	}
	var Arts FullAdvertJson
	err = json.Unmarshal([]byte(response.Body), &Arts)
	if err != nil {
		wb.allErrors = append(wb.allErrors, fmt.Sprintf("Json Unmarshal Failed: %s", err.Error()))
		wb.metrics.BadJson.Add(1)
		return
	}
	wb.metrics.GoodReq.Add(1)
	for _, val := range Arts.Data.Products {
		adv := Advert{
			Id:          val.Id,
			Link:        fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", val.Id),
			ParsingDate: time.Now(),
			Price:       val.PriceU,
			PriceSale:   val.SalePriceU,
			Name:        val.Name,
			Brand:       val.Brand,
		}
		isNewPrice, oldPrice := wb.updateMapArticle(adv)
		if isNewPrice {
			wb.Telega.AddAdvert(struct {
				Id          int
				Link        string
				ParsingDate time.Time
				Price       int
				PriceSale   int
				Name        string
				Brand       string
				OldPrice    int
			}{Id: adv.Id, Link: adv.Link, ParsingDate: adv.ParsingDate, Price: adv.Price, PriceSale: adv.PriceSale, Name: adv.Name, Brand: adv.Brand, OldPrice: oldPrice})
		}
	}
}

// GetWBAdverts Получает список товаров в категории
func (wb *Parser) GetWBAdverts(link string) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			response, err := wb.GetDataFromWb(link)
			if err != nil {
				continue
			}
			var WBAdverts WbAdvertsJson

			err = json.Unmarshal([]byte(response.Body), &WBAdverts)
			if err != nil {
				wb.allErrors = append(wb.allErrors, fmt.Sprintf("Json Unmarshal Failed: %s", err.Error()))
				wb.metrics.BadJson.Add(1)
				return
			}

			wb.metrics.GoodReq.Add(1)
			for _, val := range WBAdverts.Data.Products {
				if _, ok := wb.allAdverts[val.Id]; !ok {

					wb.allAdverts[val.Id] = Advert{
						Name:  val.Name,
						Brand: val.Brand,
						Price: val.SalePriceU,
						Link:  fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", val.Id),
					}
					log.Println(val.Name)
					if wb.canSend() {
						log.Println(val)
					}

				}
			}
		case <-wb.ctx.Done():
			return
		}
	}
}
