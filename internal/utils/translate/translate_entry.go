package translate

import (
	"net/http"
	"strings"
	"sync"

	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/model"
)

type Translater interface {
	Execute(entry *model.Entry, client *http.Client, wg *sync.WaitGroup, ak string)
	GetKey() string
}

func PostProcessEntriesTitle(feed *model.Feed, entries *model.Entries) {
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
	switch which {
	case "chatgpt":
		translatehandler = &ChatGPT{
			URL:   translateOpts[1],
			Key:   translateOpts[2],
			Model: translateOpts[3],
			To:    translateOpts[4],
		}
	case "baidu_ml":
		baiduMl := BaiduML{
			Appid:  translateOpts[1],
			Secret: translateOpts[2],
			From:   "auto",
			To:     translateOpts[3],
		}
		baiduMl.InitKey()
		translatehandler = &baiduMl
	default:
		panic("Unsupported translation service: " + which)
	}

	for _, entry := range *entries {
		wg.Add(1)
		go translatehandler.Execute(entry, client, &wg, translatehandler.GetKey())
	}
	wg.Wait()
}
