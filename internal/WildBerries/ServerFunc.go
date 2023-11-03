package WildBerries

import (
	"context"
	"fmt"
	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"net/http"
	"time"
)

func (wb *Parser) GetDataFromWb(link string) (cycletls.Response, error) {
	headers := map[string]string{}
	headers["Accept"] = "application/json, text/plain, */*"
	headers["Accept-Language"] = "es,ru;q=0.9,en;q=0.8"
	headers["Sec-Fetch-Dest"] = "empty"
	headers["Sec-Fetch-Mode"] = "cors"
	headers["Sec-Fetch-Site"] = "same-site"
	headers["X-DeviceOS"] = "0"
	headers["sec-ch-ua"] = "\"Chromium\";v=\"116\", \"Not)A;Brand\";v=\"24\", \"YaBrowser\";v=\"23\""
	headers["sec-ch-ua-mobile"] = "?1"
	headers["sec-ch-ua-platform"] = "\"Android\""
	headers["Cache-Control"] = "no-cache, no-store, must-revalidate"
	headers["Pragma"] = "no-cache"
	headers["Expires"] = "0"

	opts := cycletls.Options{
		Body:      "",
		Ja3:       wb.ja3Worker.Ja3Ciphers,
		UserAgent: wb.ja3Worker.UserAgent,
		Headers:   headers,
		Timeout:   2,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	response, err := execute(ctx, wb, link, opts, "GET")
	if err != nil {
		wb.allErrors = append(wb.allErrors, fmt.Sprintf("Request Failed: %s", err.Error()))
		wb.metrics.BadReq.Add(1)
		return response, err
	}

	if response.Status != http.StatusOK {
		wb.allErrors = append(wb.allErrors, fmt.Sprintf("Status is not ok: %d", response.Status))
		wb.metrics.BadReq.Add(1)
		return response, err
	}
	return response, nil

}

// Для обхода ошибки таймаута в либе
func execute(ctx context.Context, P *Parser, avUrl string, opts cycletls.Options, Method string) (cycletls.Response, error) {
	ch := make(chan ResponseError, 1)

	go func(Aprs *Parser) {
		resp, err := P.ja3Worker.Ja3Client.Do(avUrl, opts, Method)
		rErr := ResponseError{
			resp,
			err,
		}
		ch <- rErr
	}(P)

	select {
	case res := <-ch:
		return res.Response, res.error
	case <-ctx.Done():
		return cycletls.Response{}, ctx.Err()
	}
}

type ResponseError struct {
	cycletls.Response
	error
}

// Задержка перед отправкой
func (wb *Parser) canSend() bool {
	cS := time.Now().After(wb.timeStarted.Add(70 * time.Second))
	return cS
}
