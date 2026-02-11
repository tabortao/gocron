package notify

import (
	"testing"

	"github.com/tabortao/gocron/internal/models"
)

// TestWebHook_getActiveWebhookUrls 测试根据任务接收者ID筛选webhook地址
func TestWebHook_getActiveWebhookUrls(t *testing.T) {
	webHook := &WebHook{}

	tests := []struct {
		name             string
		webhookUrls      []models.WebhookUrl
		taskReceiverIds  string
		expectedCount    int
		expectedUrlNames []string
	}{
		{
			name: "single receiver",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
				{Id: 2, Name: "Webhook 2", Url: "https://webhook2.example.com"},
				{Id: 3, Name: "Webhook 3", Url: "https://webhook3.example.com"},
			},
			taskReceiverIds:  "2",
			expectedCount:    1,
			expectedUrlNames: []string{"Webhook 2"},
		},
		{
			name: "multiple receivers",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
				{Id: 2, Name: "Webhook 2", Url: "https://webhook2.example.com"},
				{Id: 3, Name: "Webhook 3", Url: "https://webhook3.example.com"},
			},
			taskReceiverIds:  "1,3",
			expectedCount:    2,
			expectedUrlNames: []string{"Webhook 1", "Webhook 3"},
		},
		{
			name: "no matching receivers",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
				{Id: 2, Name: "Webhook 2", Url: "https://webhook2.example.com"},
			},
			taskReceiverIds:  "99",
			expectedCount:    0,
			expectedUrlNames: []string{},
		},
		{
			name: "empty receiver ids",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
			},
			taskReceiverIds:  "",
			expectedCount:    0,
			expectedUrlNames: []string{},
		},
		{
			name: "all receivers",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
				{Id: 2, Name: "Webhook 2", Url: "https://webhook2.example.com"},
			},
			taskReceiverIds:  "1,2",
			expectedCount:    2,
			expectedUrlNames: []string{"Webhook 1", "Webhook 2"},
		},
		{
			name: "receiver ids with spaces",
			webhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
				{Id: 2, Name: "Webhook 2", Url: "https://webhook2.example.com"},
			},
			taskReceiverIds:  " 1 , 2 ",
			expectedCount:    2,
			expectedUrlNames: []string{"Webhook 1", "Webhook 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webHookSetting := models.WebHook{
				WebhookUrls: tt.webhookUrls,
			}

			msg := Message{
				"task_receiver_id": tt.taskReceiverIds,
			}

			activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)

			if len(activeUrls) != tt.expectedCount {
				t.Errorf("expected %d active urls, got %d", tt.expectedCount, len(activeUrls))
			}

			// 验证返回的webhook名称
			foundNames := make(map[string]bool)
			for _, url := range activeUrls {
				foundNames[url.Name] = true
			}

			for _, expectedName := range tt.expectedUrlNames {
				if !foundNames[expectedName] {
					t.Errorf("expected to find webhook %s, but not found", expectedName)
				}
			}
		})
	}
}

// TestWebHook_getActiveWebhookUrls_EdgeCases 测试边界情况
func TestWebHook_getActiveWebhookUrls_EdgeCases(t *testing.T) {
	webHook := &WebHook{}

	t.Run("empty webhook urls", func(t *testing.T) {
		webHookSetting := models.WebHook{
			WebhookUrls: []models.WebhookUrl{},
		}

		msg := Message{
			"task_receiver_id": "1,2,3",
		}

		activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)

		if len(activeUrls) != 0 {
			t.Errorf("expected 0 active urls, got %d", len(activeUrls))
		}
	})

	t.Run("invalid receiver id format", func(t *testing.T) {
		webHookSetting := models.WebHook{
			WebhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
			},
		}

		msg := Message{
			"task_receiver_id": "abc,def",
		}

		activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)

		if len(activeUrls) != 0 {
			t.Errorf("expected 0 active urls for invalid ids, got %d", len(activeUrls))
		}
	})

	t.Run("duplicate receiver ids", func(t *testing.T) {
		webHookSetting := models.WebHook{
			WebhookUrls: []models.WebhookUrl{
				{Id: 1, Name: "Webhook 1", Url: "https://webhook1.example.com"},
			},
		}

		msg := Message{
			"task_receiver_id": "1,1,1",
		}

		activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)

		// 实际实现：遍历webhookUrls，对每个URL检查其ID是否在receiverIds中
		// 所以即使receiverIds有重复，每个webhook也只会被添加一次
		if len(activeUrls) != 1 {
			t.Errorf("expected 1 active url (no duplicates in result), got %d", len(activeUrls))
		}
	})
}

// TestWebHook_getActiveWebhookUrls_LargeDataset 测试大数据集
func TestWebHook_getActiveWebhookUrls_LargeDataset(t *testing.T) {
	webHook := &WebHook{}

	// 创建100个webhook地址
	webhookUrls := make([]models.WebhookUrl, 100)
	for i := 0; i < 100; i++ {
		webhookUrls[i] = models.WebhookUrl{
			Id:   i + 1,
			Name: "Webhook " + string(rune(i+1)),
			Url:  "https://webhook.example.com/" + string(rune(i+1)),
		}
	}

	webHookSetting := models.WebHook{
		WebhookUrls: webhookUrls,
	}

	// 选择前50个
	receiverIds := ""
	for i := 1; i <= 50; i++ {
		if i > 1 {
			receiverIds += ","
		}
		receiverIds += string(rune(i + '0'))
	}

	msg := Message{
		"task_receiver_id": receiverIds,
	}

	activeUrls := webHook.getActiveWebhookUrls(webHookSetting, msg)

	if len(activeUrls) == 0 {
		t.Error("expected some active urls, got 0")
	}
}

// BenchmarkWebHook_getActiveWebhookUrls 性能测试
func BenchmarkWebHook_getActiveWebhookUrls(b *testing.B) {
	webHook := &WebHook{}

	webhookUrls := make([]models.WebhookUrl, 10)
	for i := 0; i < 10; i++ {
		webhookUrls[i] = models.WebhookUrl{
			Id:   i + 1,
			Name: "Webhook",
			Url:  "https://example.com",
		}
	}

	webHookSetting := models.WebHook{
		WebhookUrls: webhookUrls,
	}

	msg := Message{
		"task_receiver_id": "1,3,5,7,9",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = webHook.getActiveWebhookUrls(webHookSetting, msg)
	}
}

// TestMessage_Type 测试Message类型
func TestMessage_Type(t *testing.T) {
	msg := Message{
		"task_id":          123,
		"name":             "test task",
		"output":           "test output",
		"status":           "success",
		"task_receiver_id": "1,2,3",
	}

	// 验证可以正确获取值
	if taskId, ok := msg["task_id"].(int); !ok || taskId != 123 {
		t.Error("failed to get task_id")
	}

	if name, ok := msg["name"].(string); !ok || name != "test task" {
		t.Error("failed to get name")
	}

	if receiverId, ok := msg["task_receiver_id"].(string); !ok || receiverId != "1,2,3" {
		t.Error("failed to get task_receiver_id")
	}
}
