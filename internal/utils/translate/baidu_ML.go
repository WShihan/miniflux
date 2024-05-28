package translate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"miniflux.app/v2/internal/model"
)

type BaiduMLAuthorization struct {
	Access_Token string
	Expires_in   int64
}

type TransResult struct {
	Dst string
	Src string
}
type BaiduMLResult struct {
	From         string
	To           string
	Trans_result []TransResult
}

type BaiduMLResponse struct {
	Log_id int64
	Result BaiduMLResult
}

type BaiduML struct {
	Appid         string
	Secret        string
	From          string
	To            string
	Authorization BaiduMLAuthorization
}

func (bdml *BaiduML) InitKey() (auth BaiduMLAuthorization, ok bool) {
	client := &http.Client{}
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", bdml.Appid, bdml.Secret)
	payload := strings.NewReader("")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonErr := json.Unmarshal(body, &auth)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return
	}

	ok = true
	bdml.Authorization = auth
	return

}
func (bdml BaiduML) GetKey() string {
	return bdml.Authorization.Access_Token
}

func (bdml *BaiduML) Execute(entry *model.Entry, client *http.Client, wg *sync.WaitGroup, ak string) {
	defer wg.Done()
	url := "https://aip.baidubce.com/rpc/2.0/mt/texttrans/v1?access_token=" + ak
	var ResultRes BaiduMLResponse
	// data := fmt.Sprintf(`{
	// 	"q": "%s",
	// 	"from": "%s",
	// 	"to: "%s",
	// 	"termIds": ""
	// }`, entry.Title, bdml.From, bdml.To)
	// payload := strings.NewReader(data)

	data := fmt.Sprintf(`{
		"q": "%s",
		"from": "%s",
		"to": "%s",
		"termIds": ""
	}`, entry.Title, bdml.From, bdml.To)
	payload := strings.NewReader(data)

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonErr := json.Unmarshal(body, &ResultRes)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return
	}
	if len(ResultRes.Result.Trans_result) > 0 {
		entry.Title += "｜" + ResultRes.Result.Trans_result[0].Dst
	}

}
