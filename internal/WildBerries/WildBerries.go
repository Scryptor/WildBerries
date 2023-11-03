package WildBerries

import (
	"WildBerries/internal/Telegra"
	"context"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"sync"
	"sync/atomic"
	"time"
)

type Parser struct {
	ja3Worker     Ja3Worker
	originalLink  string
	metrics       metrics
	allAdverts    map[int]Advert
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
}

func NewParser() *Parser {
	ctx, Cancel := context.WithCancel(context.Background())
	P := Parser{
		ja3Worker: Ja3Worker{
			Ja3Client:  cycletls.Init(),
			Ja3Ciphers: "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53-49408-49409-49410-49411-49412-49413-49414-65413-129,23-0-27-10-5-45-65281-35-16-13-17513-51-18-11-43-21,29-23-24,0",
			UserAgent:  "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36",
		},
		allAdverts:    map[int]Advert{},
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
	wb.articlesMutex.Lock()
	defer wb.articlesMutex.Unlock()
	if _, ok := wb.articles[adv.Id]; !ok {
		wb.articles[adv.Id] = adv
	}
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
