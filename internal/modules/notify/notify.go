package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/tabortao/gocron/internal/modules/logger"
)

type Message map[string]interface{}

type Notifiable interface {
	Send(msg Message)
}

var queue = make(chan Message, 100)

func init() {
	go run()
}

// 把消息推入队列
func Push(msg Message) {
	queue <- msg
}

func run() {
	for msg := range queue {
		// 根据任务配置发送通知
		taskType, taskTypeOk := msg["task_type"]
		_, taskReceiverIdOk := msg["task_receiver_id"]
		_, nameOk := msg["name"]
		_, outputOk := msg["output"]
		_, statusOk := msg["status"]
		if !taskTypeOk || !taskReceiverIdOk || !nameOk || !outputOk || !statusOk {
			logger.Errorf("#notify#参数不完整#%+v", msg)
			continue
		}
		enhanceMessage(msg)
		msg["content"] = fmt.Sprintf("============\n============\n============\n任务名称: %s\n状态: %s\n输出:\n %s\n", msg["name"], msg["status"], msg["output"])
		logger.Debugf("%+v", msg)
		switch taskType.(int8) {
		case 0:
			// 邮件
			mail := Mail{}
			go mail.Send(msg)
		case 1:
			// Slack
			slack := Slack{}
			go slack.Send(msg)
		case 2:
			// WebHook
			webHook := WebHook{}
			go webHook.Send(msg)
		}
		time.Sleep(1 * time.Second)
	}
}

func parseNotifyTemplate(notifyTemplate string, msg Message) string {
	enhanceMessage(msg)
	tmpl, err := template.New("notify").Parse(notifyTemplate)
	if err != nil {
		return fmt.Sprintf("解析通知模板失败: %s", err)
	}
	var buf bytes.Buffer
	status := msg["status"]
	statusZh := msg["status_zh"]
	if v, ok := statusZh.(string); ok && strings.TrimSpace(v) == "" {
		statusZh = status
	}
	resultSummary := msg["result_summary"]
	if v, ok := resultSummary.(string); ok && strings.TrimSpace(v) == "" {
		resultSummary = msg["output"]
	}
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"TaskId":            msg["task_id"],
		"TaskName":          msg["name"],
		"Status":            status,
		"StatusZh":          statusZh,
		"IsSuccess":         msg["is_success"],
		"Result":            msg["output"],
		"ResultBody":        msg["result_body"],
		"ResultSummary":     resultSummary,
		"Host":              msg["host"],
		"ResultJsonCode":    msg["result_json_code"],
		"ResultJsonErrno":   msg["result_json_errno"],
		"ResultJsonMessage": msg["result_json_message"],
		"Remark":            msg["remark"],
	}); err != nil {
		return fmt.Sprintf("执行模板失败: %s", err)
	}

	return buf.String()
}

func enhanceMessage(msg Message) {
	status, _ := msg["status"].(string)
	isSuccess := status == "Success"
	msg["is_success"] = isSuccess
	if isSuccess {
		msg["status_zh"] = "成功"
	} else {
		msg["status_zh"] = "失败"
	}

	rawOutput, _ := msg["output"].(string)
	host, body, summary, jsonCode, jsonErrno, jsonMessage := extractOutputInfo(rawOutput)
	msg["host"] = host
	msg["result_body"] = body
	msg["result_summary"] = summary
	msg["result_json_code"] = jsonCode
	msg["result_json_errno"] = jsonErrno
	msg["result_json_message"] = jsonMessage
}

func extractOutputInfo(raw string) (host string, body string, summary string, jsonCode int, jsonErrno int, jsonMessage string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", "", 0, 0, ""
	}
	lines := strings.Split(raw, "\n")
	if len(lines) > 0 {
		first := strings.TrimSpace(lines[0])
		if strings.HasPrefix(first, "Host:") {
			host = strings.TrimSpace(strings.TrimPrefix(first, "Host:"))
			body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		} else {
			body = raw
		}
	} else {
		body = raw
	}

	if strings.HasPrefix(strings.TrimSpace(body), "{") && strings.HasSuffix(strings.TrimSpace(body), "}") {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(body), &payload); err == nil {
			if v, ok := payload["code"].(float64); ok {
				jsonCode = int(v)
			}
			if v, ok := payload["errno"].(float64); ok {
				jsonErrno = int(v)
			}
			if v, ok := payload["message"].(string); ok {
				jsonMessage = v
			}
		}
	}

	if jsonMessage != "" {
		summary = jsonMessage
	} else {
		summary = body
	}
	summary = truncateForNotify(summary, 800)
	body = truncateForNotify(body, 4000)
	return host, body, summary, jsonCode, jsonErrno, jsonMessage
}

func truncateForNotify(s string, max int) string {
	if max <= 0 {
		return s
	}
	if len(s) <= max {
		return s
	}
	return s[:max]
}
