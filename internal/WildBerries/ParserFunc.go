package WildBerries

import (
	"encoding/json"
	"fmt"
	"log"
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

// getFullAdvert Не более 400 артикулов за раз
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
