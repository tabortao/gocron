package models

import (
	"encoding/json"
	"strconv"
)

type Setting struct {
	Id    int    `gorm:"primaryKey;autoIncrement"`
	Code  string `gorm:"type:varchar(32);not null"`
	Key   string `gorm:"type:varchar(64);not null"`
	Value string `gorm:"type:varchar(4096);not null;default:''"`
}

const slackTemplate = `Task ID: {{.TaskId}}
Task Name: {{.TaskName}}
Status: {{.Status}}
Result: {{.Result}}
Remark: {{.Remark}}`

const emailTemplate = `Task ID: {{.TaskId}}
Task Name: {{.TaskName}}
Status: {{.Status}}
Result: {{.Result}}
Remark: {{.Remark}}`
const webhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}",
  "remark": "{{.Remark}}"
}
`

const serverChan3TitleTemplate = `{{.TaskName}} - {{.StatusZh}}`

const serverChan3DespTemplate = `**任务**：{{.TaskName}}（ID: {{.TaskId}}）

**状态**：{{.StatusZh}}

{{ if .Host }}**节点**：{{.Host}}

{{ end }}**摘要**：{{.ResultSummary}}

{{ if .Remark }}**备注**：{{.Remark}}{{ end }}`

const barkTitleTemplate = `{{.TaskName}} - {{.StatusZh}}`

const barkBodyTemplate = `任务：{{.TaskName}}（ID: {{.TaskId}}）
状态：{{.StatusZh}}
{{ if .Host }}节点：{{.Host}}
{{ end }}摘要：{{.ResultSummary}}
{{ if .Remark }}备注：{{.Remark}}{{ end }}`

const (
	SlackCode        = "slack"
	SlackUrlKey      = "url"
	SlackTemplateKey = "template"
	SlackChannelKey  = "channel"
)

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	WebhookCode        = "webhook"
	WebhookTemplateKey = "template"
	WebhookUrlKey      = "url"
)

const (
	ServerChan3Code             = "serverchan3"
	ServerChan3TitleTemplateKey = "title_template"
	ServerChan3DespTemplateKey  = "desp_template"
	ServerChan3UrlKey           = "url"
)

const (
	BarkCode             = "bark"
	BarkTitleTemplateKey = "title_template"
	BarkBodyTemplateKey  = "body_template"
	BarkUrlKey           = "url"
)

const (
	SystemCode          = "system"
	LogRetentionDaysKey = "log_retention_days"
	LogCleanupTimeKey   = "log_cleanup_time"
	LogFileSizeLimitKey = "log_file_size_limit"
)

// region slack配置

type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (setting *Setting) Slack() (Slack, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", SlackCode).Find(&list).Error
	slack := Slack{}
	if err != nil {
		return slack, err
	}

	setting.formatSlack(list, &slack)

	return slack, err
}

func (setting *Setting) formatSlack(list []Setting, slack *Slack) {
	for _, v := range list {
		switch v.Key {
		case SlackUrlKey:
			slack.Url = v.Value
		case SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, Channel{
				v.Id, v.Value,
			})
		}
	}
}

func (setting *Setting) UpdateSlack(url, template string) error {
	setting.Value = url
	Db.Model(&Setting{}).Where("code = ? AND `key` = ?", SlackCode, SlackUrlKey).Update("value", url)

	setting.Value = template
	Db.Model(&Setting{}).Where("code = ? AND `key` = ?", SlackCode, SlackTemplateKey).Update("value", template)

	return nil
}

// 创建slack渠道
func (setting *Setting) CreateChannel(channel string) (int64, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	result := Db.Create(setting)
	return result.RowsAffected, result.Error
}

func (setting *Setting) IsChannelExist(channel string) bool {
	var count int64
	Db.Model(&Setting{}).Where("code = ? AND `key` = ? AND value = ?", SlackCode, SlackChannelKey, channel).Count(&count)
	return count > 0
}

// 删除slack渠道
func (setting *Setting) RemoveChannel(id int) (int64, error) {
	result := Db.Where("code = ? AND `key` = ? AND id = ?", SlackCode, SlackChannelKey, id).Delete(&Setting{})
	return result.RowsAffected, result.Error
}

// endregion

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mail_users"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// region 邮件配置
func (setting *Setting) Mail() (Mail, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", MailCode).Find(&list).Error
	mail := Mail{MailUsers: make([]MailUser, 0)}
	if err != nil {
		return mail, err
	}

	setting.formatMail(list, &mail)

	return mail, err
}

func (setting *Setting) formatMail(list []Setting, mail *Mail) {
	mailUser := MailUser{}
	for _, v := range list {
		switch v.Key {
		case MailServerKey:
			if v.Value != "" {
				_ = json.Unmarshal([]byte(v.Value), mail)
			}
		case MailUserKey:
			if v.Value != "" {
				_ = json.Unmarshal([]byte(v.Value), &mailUser)
				mailUser.Id = v.Id
				mail.MailUsers = append(mail.MailUsers, mailUser)
			}
		case MailTemplateKey:
			mail.Template = v.Value
		}

	}
}

func (setting *Setting) UpdateMail(config, template string) error {
	Db.Model(&Setting{}).Where("code = ? AND `key` = ?", MailCode, MailServerKey).Update("value", config)
	Db.Model(&Setting{}).Where("code = ? AND `key` = ?", MailCode, MailTemplateKey).Update("value", template)

	return nil
}

func (setting *Setting) CreateMailUser(username, email string) (int64, error) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	mailUser := MailUser{0, username, email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		return 0, err
	}
	setting.Value = string(jsonByte)

	result := Db.Create(setting)
	return result.RowsAffected, result.Error
}

func (setting *Setting) RemoveMailUser(id int) (int64, error) {
	result := Db.Where("code = ? AND `key` = ? AND id = ?", MailCode, MailUserKey, id).Delete(&Setting{})
	return result.RowsAffected, result.Error
}

type WebHook struct {
	WebhookUrls []WebhookUrl `json:"webhook_urls"`
	Template    string       `json:"template"`
}

type WebhookUrl struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (setting *Setting) Webhook() (WebHook, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", WebhookCode).Find(&list).Error
	webHook := WebHook{WebhookUrls: make([]WebhookUrl, 0)}
	if err != nil {
		return webHook, err
	}

	setting.formatWebhook(list, &webHook)

	return webHook, err
}

func (setting *Setting) formatWebhook(list []Setting, webHook *WebHook) {
	webhookUrl := WebhookUrl{}
	for _, v := range list {
		switch v.Key {
		case WebhookUrlKey:
			if v.Value != "" {
				_ = json.Unmarshal([]byte(v.Value), &webhookUrl)
				webhookUrl.Id = v.Id
				webHook.WebhookUrls = append(webHook.WebhookUrls, webhookUrl)
			}
		case WebhookTemplateKey:
			webHook.Template = v.Value
		}
	}
}

func (setting *Setting) UpdateWebHook(template string) error {
	Db.Model(&Setting{}).Where("code = ? AND `key` = ?", WebhookCode, WebhookTemplateKey).Update("value", template)
	return nil
}

func (setting *Setting) CreateWebhookUrl(name, url string) (int64, error) {
	webhookUrl := WebhookUrl{0, name, url}
	jsonByte, err := json.Marshal(webhookUrl)
	if err != nil {
		return 0, err
	}

	newSetting := Setting{
		Code:  WebhookCode,
		Key:   WebhookUrlKey,
		Value: string(jsonByte),
	}

	result := Db.Create(&newSetting)
	return result.RowsAffected, result.Error
}

func (setting *Setting) RemoveWebhookUrl(id int) (int64, error) {
	result := Db.Where("code = ? AND `key` = ? AND id = ?", WebhookCode, WebhookUrlKey, id).Delete(&Setting{})
	return result.RowsAffected, result.Error
}

// endregion

type ServerChan3 struct {
	Urls          []ServerChan3Url `json:"urls"`
	TitleTemplate string           `json:"title_template"`
	DespTemplate  string           `json:"desp_template"`
}

type ServerChan3Url struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (setting *Setting) ServerChan3() (ServerChan3, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", ServerChan3Code).Find(&list).Error
	serverChan3 := ServerChan3{Urls: make([]ServerChan3Url, 0)}
	if err != nil {
		return serverChan3, err
	}

	setting.formatServerChan3(list, &serverChan3)
	return serverChan3, nil
}

func (setting *Setting) formatServerChan3(list []Setting, serverChan3 *ServerChan3) {
	urlItem := ServerChan3Url{}
	for _, v := range list {
		switch v.Key {
		case ServerChan3UrlKey:
			if v.Value != "" {
				_ = json.Unmarshal([]byte(v.Value), &urlItem)
				urlItem.Id = v.Id
				serverChan3.Urls = append(serverChan3.Urls, urlItem)
			}
		case ServerChan3TitleTemplateKey:
			serverChan3.TitleTemplate = v.Value
		case ServerChan3DespTemplateKey:
			serverChan3.DespTemplate = v.Value
		}
	}
}

func (setting *Setting) UpdateServerChan3(titleTemplate, despTemplate string) error {
	Db.Model(&Setting{}).
		Where("code = ? AND `key` = ?", ServerChan3Code, ServerChan3TitleTemplateKey).
		Update("value", titleTemplate)
	Db.Model(&Setting{}).
		Where("code = ? AND `key` = ?", ServerChan3Code, ServerChan3DespTemplateKey).
		Update("value", despTemplate)
	return nil
}

func (setting *Setting) CreateServerChan3Url(name, url string) (int64, error) {
	urlItem := ServerChan3Url{0, name, url}
	jsonByte, err := json.Marshal(urlItem)
	if err != nil {
		return 0, err
	}

	newSetting := Setting{
		Code:  ServerChan3Code,
		Key:   ServerChan3UrlKey,
		Value: string(jsonByte),
	}

	result := Db.Create(&newSetting)
	return result.RowsAffected, result.Error
}

func (setting *Setting) RemoveServerChan3Url(id int) (int64, error) {
	result := Db.Where("code = ? AND `key` = ? AND id = ?", ServerChan3Code, ServerChan3UrlKey, id).Delete(&Setting{})
	return result.RowsAffected, result.Error
}

type Bark struct {
	Urls          []BarkUrl `json:"urls"`
	TitleTemplate string    `json:"title_template"`
	BodyTemplate  string    `json:"body_template"`
}

type BarkUrl struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (setting *Setting) Bark() (Bark, error) {
	list := make([]Setting, 0)
	err := Db.Where("code = ?", BarkCode).Find(&list).Error
	bark := Bark{Urls: make([]BarkUrl, 0)}
	if err != nil {
		return bark, err
	}

	setting.formatBark(list, &bark)
	return bark, nil
}

func (setting *Setting) formatBark(list []Setting, bark *Bark) {
	urlItem := BarkUrl{}
	for _, v := range list {
		switch v.Key {
		case BarkUrlKey:
			if v.Value != "" {
				_ = json.Unmarshal([]byte(v.Value), &urlItem)
				urlItem.Id = v.Id
				bark.Urls = append(bark.Urls, urlItem)
			}
		case BarkTitleTemplateKey:
			bark.TitleTemplate = v.Value
		case BarkBodyTemplateKey:
			bark.BodyTemplate = v.Value
		}
	}
}

func (setting *Setting) UpdateBark(titleTemplate, bodyTemplate string) error {
	Db.Model(&Setting{}).
		Where("code = ? AND `key` = ?", BarkCode, BarkTitleTemplateKey).
		Update("value", titleTemplate)
	Db.Model(&Setting{}).
		Where("code = ? AND `key` = ?", BarkCode, BarkBodyTemplateKey).
		Update("value", bodyTemplate)
	return nil
}

func (setting *Setting) CreateBarkUrl(name, url string) (int64, error) {
	urlItem := BarkUrl{0, name, url}
	jsonByte, err := json.Marshal(urlItem)
	if err != nil {
		return 0, err
	}

	newSetting := Setting{
		Code:  BarkCode,
		Key:   BarkUrlKey,
		Value: string(jsonByte),
	}

	result := Db.Create(&newSetting)
	return result.RowsAffected, result.Error
}

func (setting *Setting) RemoveBarkUrl(id int) (int64, error) {
	result := Db.Where("code = ? AND `key` = ? AND id = ?", BarkCode, BarkUrlKey, id).Delete(&Setting{})
	return result.RowsAffected, result.Error
}

// region 通用配置辅助方法

// getSettingValue 获取配置值的通用方法
func (setting *Setting) getSettingValue(code, key string) (string, error) {
	var s Setting
	err := Db.Where("code = ? AND `key` = ?", code, key).First(&s).Error
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

// updateOrCreateSetting 更新或创建配置的通用方法
func (setting *Setting) updateOrCreateSetting(code, key, value string) error {
	var s Setting
	err := Db.Where("code = ? AND `key` = ?", code, key).First(&s).Error
	if err != nil {
		// 记录不存在，创建新记录
		s.Code = code
		s.Key = key
		s.Value = value
		result := Db.Create(&s)
		return result.Error
	}
	// 记录存在，更新
	result := Db.Model(&Setting{}).Where("code = ? AND `key` = ?", code, key).Update("value", value)
	return result.Error
}

// endregion

// region 系统配置
func (setting *Setting) GetLogRetentionDays() int {
	value, err := setting.getSettingValue(SystemCode, LogRetentionDaysKey)
	if err != nil || value == "" {
		return 0
	}
	days, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return days
}

func (setting *Setting) UpdateLogRetentionDays(days int) error {
	return setting.updateOrCreateSetting(SystemCode, LogRetentionDaysKey, strconv.Itoa(days))
}

func (setting *Setting) GetLogCleanupTime() string {
	value, err := setting.getSettingValue(SystemCode, LogCleanupTimeKey)
	if err != nil || value == "" {
		return "03:00"
	}
	return value
}

func (setting *Setting) UpdateLogCleanupTime(cleanupTime string) error {
	return setting.updateOrCreateSetting(SystemCode, LogCleanupTimeKey, cleanupTime)
}

func (setting *Setting) GetLogFileSizeLimit() int {
	value, err := setting.getSettingValue(SystemCode, LogFileSizeLimitKey)
	if err != nil || value == "" {
		return 0
	}
	size, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return size
}

func (setting *Setting) UpdateLogFileSizeLimit(size int) error {
	return setting.updateOrCreateSetting(SystemCode, LogFileSizeLimitKey, strconv.Itoa(size))
}

// endregion
