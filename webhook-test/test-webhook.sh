#!/bin/bash

echo "🧪 测试Webhook服务..."

# 测试数据
test_data='{
  "task_id": 123,
  "task_name": "测试任务",
  "status": "成功",
  "output": "任务执行完成",
  "remark": "这是一个测试webhook"
}'

echo "📤 发送测试数据到webhook服务..."
echo "数据: $test_data"

# 发送POST请求
curl -X POST \
  -H "Content-Type: application/json" \
  -d "$test_data" \
  http://localhost:8080/webhook

echo -e "\n\n✅ 测试完成"