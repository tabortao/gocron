package notify

import (
	"html"
	urlpkg "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/httpclient"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/utils"
)

type WebHook struct{}

const webHookAllReceiverId = "-2"

func (webHook *WebHook) Send(msg Message) {
	model := new(models.Setting)
	webHookSetting, err := model.Webhook()
	if err != nil {
		logger.Error("#webHook#从数据库获取webHook配置失败", err)
		return
	}
	if len(webHookSetting.WebhookUrls) == 0 {
		logger.Error("#webHook#webhook地址列表为空")
		return
	}
	for _, key := range []string{
		"name",
		"output",
		"host",
		"result_body",
		"result_summary",
		"result_json_message",
		"remark",
		"status_zh",
	} {
		if value, ok := msg[key].(string); ok {
			msg[key] = utils.EscapeJson(value)
		}
	}
	msg["content"] = parseNotifyTemplate(webHookSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))

	// 获取任务配置的接收者ID列表
	activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)
	if len(activeUrls) == 0 {
		receiverId, _ := msg["task_receiver_id"].(string)
		logger.Warnf("#webHook#未匹配到webhook接收地址#task_id=%v#receiver_id=%s", msg["task_id"], receiverId)
		return
	}

	// 向所有激活的webhook地址发送
	for _, webhookUrl := range activeUrls {
		go webHook.send(msg, webhookUrl.Url)
	}
}

func (webHook *WebHook) getActiveWebhookUrls(webHookSetting models.WebHook, msg Message) []models.WebhookUrl {
	taskReceiverId, _ := msg["task_receiver_id"].(string)
	taskReceiverId = strings.TrimSpace(taskReceiverId)
	if taskReceiverId == "" {
		return webHookSetting.WebhookUrls
	}
	taskReceiverIdsRaw := strings.Split(taskReceiverId, ",")
	taskReceiverIds := make([]string, 0, len(taskReceiverIdsRaw))
	for _, id := range taskReceiverIdsRaw {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		taskReceiverIds = append(taskReceiverIds, id)
	}
	if utils.InStringSlice(taskReceiverIds, webHookAllReceiverId) {
		return webHookSetting.WebhookUrls
	}
	urls := []models.WebhookUrl{}
	for _, v := range webHookSetting.WebhookUrls {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			urls = append(urls, v)
		}
	}
	return urls
}

func (webHook *WebHook) send(msg Message, url string) {
	content, _ := msg["content"].(string)
	timeout := 30
	maxTimes := 3
	for i := 0; i < maxTimes; i++ {
		resp := httpclient.PostJson(url, content, timeout)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return
		}
		if i < maxTimes-1 {
			time.Sleep(2 * time.Second)
		}
		logger.Errorf("#webHook#发送失败#url=%s#attempt=%d/%d#status=%d#resp=%s", redactWebhookUrl(url), i+1, maxTimes, resp.StatusCode, truncateString(resp.Body, 1000))
	}
}

func redactWebhookUrl(raw string) string {
	u, err := urlpkg.Parse(raw)
	if err != nil {
		return "<invalid_url>"
	}
	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if last != "" {
			suffix := last
			if len(suffix) > 6 {
				suffix = suffix[len(suffix)-6:]
			}
			parts[len(parts)-1] = "***" + suffix
			path = strings.Join(parts, "/")
		}
	}
	if u.Scheme == "" {
		return u.Host + path
	}
	return u.Scheme + "://" + u.Host + path
}

func truncateString(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max]
}
