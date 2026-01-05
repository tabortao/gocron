package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// WebhookPayload webhook接收的数据结构
type WebhookPayload struct {
	TaskID   int    `json:"task_id"`
	TaskName string `json:"task_name"`
	Status   string `json:"status"`
	Output   string `json:"output"`
	Remark   string `json:"remark"`
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/health", handleHealth)

	fmt.Println("🚀 Webhook测试服务启动")
	fmt.Println("📡 监听地址: http://localhost:8080/webhook")
	fmt.Println("💚 健康检查: http://localhost:8080/health")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// 记录请求时间
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	fmt.Printf("\n=== [%s] 收到Webhook请求 ===\n", timestamp)
	fmt.Printf("方法: %s\n", r.Method)
	fmt.Printf("路径: %s\n", r.URL.Path)
	
	// 打印请求头
	fmt.Println("请求头:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}
	
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("❌ 读取请求体失败: %v\n", err)
		http.Error(w, "读取请求体失败", http.StatusBadRequest)
		return
	}
	
	fmt.Printf("请求体: %s\n", string(body))
	
	// 尝试解析JSON
	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("⚠️  JSON解析失败: %v\n", err)
		fmt.Println("将作为纯文本处理")
	} else {
		fmt.Println("✅ JSON解析成功:")
		fmt.Printf("  任务ID: %d\n", payload.TaskID)
		fmt.Printf("  任务名称: %s\n", payload.TaskName)
		fmt.Printf("  状态: %s\n", payload.Status)
		fmt.Printf("  输出: %s\n", payload.Output)
		fmt.Printf("  备注: %s\n", payload.Remark)
	}
	
	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"success":   true,
		"message":   "webhook接收成功",
		"timestamp": timestamp,
		"received":  len(body) > 0,
	}
	
	json.NewEncoder(w).Encode(response)
	fmt.Println("✅ 响应已发送")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"status": "ok",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	}
	
	json.NewEncoder(w).Encode(response)
}