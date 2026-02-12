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

type ServerChan3 struct{}

const serverChan3AllReceiverId = "-3"
const serverChan3TypedPrefix = "c"

func (serverChan3 *ServerChan3) Send(msg Message) {
	model := new(models.Setting)
	setting, err := model.ServerChan3()
	if err != nil {
		logger.Error("#serverchan3#从数据库获取配置失败", err)
		return
	}
	if len(setting.Urls) == 0 {
		logger.Error("#serverchan3#地址列表为空")
		return
	}

	activeUrls := serverChan3.getActiveUrls(setting, msg)
	if len(activeUrls) == 0 {
		receiverId, _ := msg["task_receiver_id"].(string)
		logger.Warnf("#serverchan3#未匹配到接收地址#task_id=%v#receiver_id=%s", msg["task_id"], receiverId)
		return
	}

	title := parseNotifyTemplate(setting.TitleTemplate, msg)
	desp := parseNotifyTemplate(setting.DespTemplate, msg)
	title = html.UnescapeString(title)
	desp = html.UnescapeString(desp)

	for _, u := range activeUrls {
		go serverChan3.send(u.Url, title, desp)
	}
}

func (serverChan3 *ServerChan3) getActiveUrls(setting models.ServerChan3, msg Message) []models.ServerChan3Url {
	raw, _ := msg["task_receiver_id"].(string)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return setting.Urls
	}

	typed, legacy := parseReceiverTokens(raw)
	if values, ok := typed[serverChan3TypedPrefix]; ok && len(values) > 0 {
		if containsWildcard(values) || utils.InStringSlice(values, serverChan3AllReceiverId) {
			return setting.Urls
		}
		idSet := toIntSet(values)
		urls := make([]models.ServerChan3Url, 0, len(setting.Urls))
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
	if utils.InStringSlice(taskReceiverIds, serverChan3AllReceiverId) {
		return setting.Urls
	}

	urls := make([]models.ServerChan3Url, 0, len(setting.Urls))
	for _, v := range setting.Urls {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			urls = append(urls, v)
		}
	}
	return urls
}

func (serverChan3 *ServerChan3) send(apiUrl, title, desp string) {
	timeout := 30
	maxTimes := 3

	params := url.Values{}
	params.Set("title", title)
	params.Set("desp", desp)

	for i := 0; i < maxTimes; i++ {
		resp := httpclient.PostParams(apiUrl, params.Encode(), timeout)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return
		}
		if i < maxTimes-1 {
			time.Sleep(2 * time.Second)
		}
		logger.Errorf("#serverchan3#发送失败#status=%d#resp=%s", resp.StatusCode, truncate(resp.Body, 1000))
	}
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max]
}
