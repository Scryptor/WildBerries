package Telegra

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Sender struct {
	Client      *http.Client
	Chat        int64
	TgId        int64
	Token       string
	AdvertsList []Advert
	ctx         context.Context
	rwMute      sync.RWMutex
	Metrics
	SearchName string
}

type Metrics struct {
	SendGood atomic.Uint32
	SendBad  atomic.Uint32
}

type Advert struct {
	Id          int
	Link        string
	ParsingDate time.Time
	Price       int
	PriceSale   int
	Name        string
	Brand       string
	OldPrice    int
}

func NewSender(ctx context.Context, tgid, chat int64) *Sender {
	cctx, _ := context.WithCancel(ctx)
	sndr := Sender{
		Client:      &http.Client{Timeout: 10 * time.Second},
		Token:       "6917956424:AAE0CG0KfslaxAQt3QwnOINHYOSxzqwGtj8",
		Chat:        chat,
		AdvertsList: make([]Advert, 0, 50),
		ctx:         cctx,
		TgId:        tgid,
	}
	go sndr.StartSending()
	return &sndr
}

// AddAdvert Добавляет товар для последующей отправки в телеграм
func (Tgs *Sender) AddAdvert(adv Advert) {
	Tgs.rwMute.Lock()
	defer Tgs.rwMute.Unlock()
	Tgs.AdvertsList = append(Tgs.AdvertsList, adv)
}

// StartSending Запускает тикер отправки
func (Tgs *Sender) StartSending() {
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			if len(Tgs.AdvertsList) > 0 {
				err := Tgs.sendWithPhotoIfNoFullAdvert(Tgs.AdvertsList[len(Tgs.AdvertsList)-1])
				if err != nil {
					log.Println(err)
					Tgs.SendBad.Add(1)
				} else {
					Tgs.SendGood.Add(1)
				}
				Tgs.rwMute.Lock()
				//Tgs.AdvertsList = append(Tgs.AdvertsList[:0], Tgs.AdvertsList[1:]...)
				Tgs.AdvertsList = Tgs.AdvertsList[:len(Tgs.AdvertsList)-1]
				Tgs.rwMute.Unlock()
			}
		case <-Tgs.ctx.Done():
			log.Println("Отправка объявлений отменена")
			return
		}

	}
}

func (Tgs *Sender) sendWithPhotoIfNoFullAdvert(fAdvert Advert) error {
	// Если нет фото
	img := "https://mywork2.ru/nophotos.jpg"

	mdPrice := fmt.Sprintf("%d", fAdvert.PriceSale/100)

	caption := fmt.Sprintf("*%s*\n💵 *%s*\n\n%s  \n\n`%s`\n",
		MarkdownCorrector(fAdvert.Name),
		MarkdownCorrector(mdPrice),
		MarkdownCorrector(fAdvert.Link),
		MarkdownCorrector(Tgs.SearchName),
	)
	err := Tgs.sendPostMessage(caption, img)
	if err != nil {
		return err
	}
	return nil
}

func (Tgs *Sender) sendPostMessage(message string, photoUrl string) error {
	data := url.Values{}

	data.Add("chat_id", fmt.Sprintf("%d", Tgs.Chat))
	data.Add("text", message)
	var sendMode string

	sendMode = "sendPhoto"
	data.Add("photo", photoUrl)
	data.Add("caption", message)
	data.Add("parse_mode", "MarkdownV2")

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/%s", Tgs.Token, sendMode), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", " gzip, deflate")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Tgs.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprint("status is not ok", resp.StatusCode, string(body)))
	}
	return nil
}

// SendHelloMessage Отправляет приветственное сообщение при новом поиске
func (Tgs *Sender) SendHelloMessage(message, link, name string) error {
	mln := fmt.Sprintf("*%s*\n\n`%s`\n\n%s", MarkdownCorrector(message), MarkdownCorrector(name), MarkdownCorrector(link))

	data := url.Values{}
	data.Add("text", mln)
	data.Add("chat_id", fmt.Sprintf("%d", Tgs.Chat))
	var sendMode string

	sendMode = "sendMessage"
	data.Add("parse_mode", "MarkdownV2")
	data.Add("disable_web_page_preview", "True")

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/%s", Tgs.Token, sendMode), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", " gzip, deflate")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Tgs.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprint("status is not ok", resp.StatusCode, string(body)))
	}
	return nil
}

// MarkdownCorrector Корректировка разметки при отправке в телегу, без этого не отправит
func MarkdownCorrector(text string) string {
	text = strings.ReplaceAll(text, "\\", "")
	text = strings.ReplaceAll(text, "-", "\\-")
	text = strings.ReplaceAll(text, "_", "\\_")
	text = strings.ReplaceAll(text, "*", "\\*")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "~", "\\~")
	text = strings.ReplaceAll(text, "`", "\\`")
	text = strings.ReplaceAll(text, ">", "\\>")
	text = strings.ReplaceAll(text, "#", "\\#")
	text = strings.ReplaceAll(text, "+", "\\+")
	text = strings.ReplaceAll(text, "=", "\\=")
	text = strings.ReplaceAll(text, "|", "\\|")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, ".", "\\.")
	text = strings.ReplaceAll(text, "!", "\\!")
	return text
}
