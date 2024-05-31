package translate

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/model"
)

type Translater interface {
	Execute(ch chan struct{}, entry *model.Entry, client *http.Client, wg *sync.WaitGroup, ak string)
	GetKey() string
}

func PostProcessEntriesTitle(feed *model.Feed, entries *model.Entries) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if !ok {
				slog.Error(fmt.Errorf("pkg: %v", r).Error())
				slog.Error(err.Error())
			}
			return
		}
	}()
	if !feed.Translatable {
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}

	var wg sync.WaitGroup
	var translatehandler Translater

	translateURL := config.Opts.TranslateURL()
	if translateURL == "" {
		return
	}

	translateOpts := strings.Split(translateURL, "@")
	which := strings.ToLower(translateOpts[0])
	// refer to translate api max request times per second
	numWorkers := 5

	switch which {
	case "chatgpt":
		translatehandler = &ChatGPT{
			URL:   translateOpts[1],
			Key:   translateOpts[2],
			Model: translateOpts[3],
			To:    translateOpts[4],
		}
		numWorkers = 20
	case "baidu_ml":
		baiduMl := BaiduML{
			Appid:  translateOpts[1],
			Secret: translateOpts[2],
			From:   "auto",
			To:     translateOpts[3],
		}
		baiduMl.InitKey()
		translatehandler = &baiduMl
	case "tencent_ml":
		projID, _ := strconv.ParseInt(translateOpts[3], 10, 64)
		tencentMl := TencentML{
			SecretId:  translateOpts[1],
			SecretKey: translateOpts[2],
			ProjectId: projID,
			Target:    translateOpts[4],
			Source:    "auto",
		}
		tencentMl.InitKey(&tencentMl)
		translatehandler = &tencentMl
		numWorkers = 10
	default:
		slog.Error("Unsupported translation service: " + which)
	}
	sem := make(chan struct{}, numWorkers)
	for _, entry := range *entries {
		wg.Add(1)
		go translatehandler.Execute(sem, entry, client, &wg, translatehandler.GetKey())
	}
	wg.Wait()
	slog.Info("finish")
}
