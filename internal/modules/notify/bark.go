package notify

import (
	"html"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/httpclient"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/utils"
)

type Bark struct{}

const barkAllReceiverId = "-4"
const barkTypedPrefix = "b"

func (bark *Bark) Send(msg Message) {
	model := new(models.Setting)
	setting, err := model.Bark()
	if err != nil {
		logger.Error("#bark#从数据库获取配置失败", err)
		return
	}
	if len(setting.Urls) == 0 {
		logger.Error("#bark#地址列表为空")
		return
	}

	activeUrls := bark.getActiveUrls(setting, msg)
	if len(activeUrls) == 0 {
		receiverId, _ := msg["task_receiver_id"].(string)
		logger.Warnf("#bark#未匹配到接收地址#task_id=%v#receiver_id=%s", msg["task_id"], receiverId)
		return
	}

	title := parseNotifyTemplate(setting.TitleTemplate, msg)
	body := parseNotifyTemplate(setting.BodyTemplate, msg)
	title = html.UnescapeString(title)
	body = html.UnescapeString(body)

	for _, u := range activeUrls {
		go bark.send(u.Url, title, body)
	}
}

func (bark *Bark) getActiveUrls(setting models.Bark, msg Message) []models.BarkUrl {
	raw, _ := msg["task_receiver_id"].(string)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return setting.Urls
	}

	typed, legacy := parseReceiverTokens(raw)
	if values, ok := typed[barkTypedPrefix]; ok && len(values) > 0 {
		if containsWildcard(values) || utils.InStringSlice(values, barkAllReceiverId) {
			return setting.Urls
		}
		idSet := toIntSet(values)
		urls := make([]models.BarkUrl, 0, len(setting.Urls))
		for _, v := range setting.Urls {
			if _, ok := idSet[v.Id]; ok {
				urls = append(urls, v)
			}
		}
		return urls
	}

	taskReceiverIds := make([]string, 0, len(legacy))
	for _, id := range legacy {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		taskReceiverIds = append(taskReceiverIds, id)
	}
	if utils.InStringSlice(taskReceiverIds, barkAllReceiverId) {
		return setting.Urls
	}

	urls := make([]models.BarkUrl, 0, len(setting.Urls))
	for _, v := range setting.Urls {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			urls = append(urls, v)
		}
	}
	return urls
}

func (bark *Bark) send(apiUrl, title, body string) {
	timeout := 30
	maxTimes := 3

	finalUrl := strings.TrimRight(apiUrl, "/") + "/" + url.PathEscape(title) + "/" + url.PathEscape(body)

	for i := 0; i < maxTimes; i++ {
		resp := httpclient.Get(finalUrl, timeout)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return
		}
		if i < maxTimes-1 {
			time.Sleep(2 * time.Second)
		}
		logger.Errorf("#bark#发送失败#status=%d#resp=%s", resp.StatusCode, truncate(resp.Body, 1000))
	}
}

