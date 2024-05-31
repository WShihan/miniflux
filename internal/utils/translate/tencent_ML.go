package translate

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"miniflux.app/v2/internal/model"
)

type TencentML struct {
	SecretId  string
	SecretKey string
	ProjectId int64
	Source    string
	Target    string
	Server    *tmt.Client
}

func (tencentML TencentML) GetKey() string {
	return ""
}

func (tencentML TencentML) InitKey(instance *TencentML) {
	credential := common.NewCredential(
		tencentML.SecretId,
		tencentML.SecretKey,
	)

	server, err := tmt.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(err.Error())
	}
	instance.Server = server

}
func (tencentML TencentML) Execute(sem chan struct{}, entry *model.Entry, client *http.Client, wg *sync.WaitGroup, ak string) {
	// wg.Done()
	defer func() {
		wg.Done()
		sem <- struct{}{} // 获取信号量
		<-sem             // 释放信号量
	}()
	msg := entry.Title

	request := tmt.NewTextTranslateRequest()
	request.SourceText = &msg
	request.Source = &tencentML.Source
	request.Target = &tencentML.Target
	request.ProjectId = &tencentML.ProjectId

	response, err := tencentML.Server.TextTranslate(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		slog.Error(err.Error())
		return
	}
	// fmt.Printf("%s", (*response.Response.TargetText))
	content := (*response.Response.TargetText)
	slog.Info(fmt.Sprintf("Translate title:%s, result:%s", entry.Title, content))
	entry.Title += "｜" + content
}
