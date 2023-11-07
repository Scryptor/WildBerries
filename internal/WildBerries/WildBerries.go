package WildBerries

import (
	"WildBerries/internal/Telegra"
	"context"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Parser struct {
	ja3Worker     Ja3Worker
	originalLink  string
	metrics       metrics
	allErrors     []string
	ctx           context.Context
	CancelWork    context.CancelFunc
	timeStarted   time.Time
	articles      map[int]Advert
	articlesMutex sync.RWMutex
	Telega        *Telegra.Sender
}

type Ja3Worker struct {
	Ja3Client  cycletls.CycleTLS
	Ja3Ciphers string
	UserAgent  string
}

type metrics struct {
	GoodReq atomic.Uint32
	BadReq  atomic.Uint32
	BadJson atomic.Uint32
}
type Advert struct {
	Id          int
	Link        string
	ParsingDate time.Time
	Price       int
	PriceSale   int
	Name        string
	Brand       string
	Images      []string
	Pics        int
}

func NewParser() *Parser {
	ctx, Cancel := context.WithCancel(context.Background())
	P := Parser{
		ja3Worker: Ja3Worker{
			Ja3Client:  cycletls.Init(),
			Ja3Ciphers: "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-49408-49409-49410-49411-49412-49413-49414-65413-129,23-0-27-10-5-45-65281-35-16-13-17513-51-18-11-43-21,29-23-24,0",
			UserAgent:  "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
		},
		allErrors:     make([]string, 0, 50),
		ctx:           ctx,
		CancelWork:    Cancel,
		timeStarted:   time.Now(),
		articles:      map[int]Advert{},
		articlesMutex: sync.RWMutex{},
		Telega:        Telegra.NewSender(ctx),
	}
	go P.updateArticles()
	return &P
}

// Добавляет новый товар в коллекцию для последующего обхода
func (wb *Parser) addMapArticle(adv Advert) {
	imgs := make([]string, 0, adv.Pics)
	link := fmt.Sprintf("https://www.wildberries.ru/catalog/%d/detail.aspx", adv.Id)
	img := wb.getImageUrlById(adv.Id)
	for i := 1; i <= adv.Pics; i++ {
		imgs = append(imgs, strings.ReplaceAll(img, "XXX", fmt.Sprintf("%d", i)))
	}
	adv.Images = imgs
	adv.Link = link
	wb.articlesMutex.Lock()
	defer wb.articlesMutex.Unlock()
	if _, ok := wb.articles[adv.Id]; !ok {
		wb.articles[adv.Id] = adv
	}
}

// Получает url картинки товара
func (wb *Parser) getImageUrlById(id int) string {
	var basket string
	var vol, part int
	vol = id / 100000
	part = id / 1000
	switch {
	case vol >= 0 && vol <= 143:
		basket = "01"
	case vol >= 144 && vol <= 287:
		basket = "02"
	case vol >= 288 && vol <= 431:
		basket = "03"
	case vol >= 432 && vol <= 719:
		basket = "04"
	case vol >= 720 && vol <= 1007:
		basket = "05"
	case vol >= 1008 && vol <= 1061:
		basket = "06"
	case vol >= 1062 && vol <= 1115:
		basket = "07"
	case vol >= 1116 && vol <= 1169:
		basket = "08"
	case vol >= 1170 && vol <= 1313:
		basket = "09"
	case vol >= 1314 && vol <= 1601:
		basket = "10"
	case vol >= 1602 && vol <= 1655:
		basket = "11"
	case vol >= 1656 && vol <= 1919:
		basket = "12"
	case vol >= 1920 && vol <= 2045:
		basket = "13"
	default:
		basket = "14"
	}
	img := fmt.Sprintf("https://basket-%s.wb.ru/vol%d/part%d/%d/images/big/XXX.webp", basket, vol, part, id)
	return img
}

// Обновляет информацию по уже имеющимся товарам
func (wb *Parser) updateMapArticle(adv Advert) (bool, int) {
	wb.articlesMutex.Lock()
	defer wb.articlesMutex.Unlock()
	if art, ok := wb.articles[adv.Id]; ok {
		if art.PriceSale > adv.PriceSale {
			priceOld := art.PriceSale
			wb.articles[adv.Id] = adv
			return true, priceOld
		}
	}
	return false, 0
}

// SetOrigLink Устанавливает ссылку для поиска по категории
func (wb *Parser) SetOrigLink(link string) {
	wb.originalLink = link
}

func (wb *Parser) GetArticlesCount() int {
	wb.articlesMutex.RLock()
	defer wb.articlesMutex.RUnlock()
	return len(wb.articles)
}
