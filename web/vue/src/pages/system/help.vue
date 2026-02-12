<template>
  <el-main class="help-page">
    <div class="help-hero">
      <div class="help-hero__title">{{ t('system.help') }}</div>
      <div class="help-hero__subtitle">
        {{
          isZh
            ? '常见配置与排错入口，建议先从「快速入口」和「飞书机器人通知」开始。'
            : 'Common setup and troubleshooting. Start with Quick links and Feishu bot notifications.'
        }}
      </div>
    </div>

    <el-row :gutter="16" class="help-grid">
      <el-col :xs="24" :md="12">
        <el-card class="help-card" shadow="never">
          <template #header>
            <div class="help-card__header">
              <span>{{ isZh ? '快速入口' : 'Quick links' }}</span>
            </div>
          </template>
          <div class="help-links">
            <div class="help-link-row">
              <div class="help-link-row__label">{{ isZh ? 'NAS 教程' : 'NAS guide' }}</div>
              <el-link
                :href="nasDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                docs/nas_user_guide.md
              </el-link>
            </div>
            <div class="help-link-row">
              <div class="help-link-row__label">{{ isZh ? 'Windows 教程' : 'Windows guide' }}</div>
              <el-link
                :href="windowsDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                docs/windows_user_guide.md
              </el-link>
            </div>
            <div class="help-link-row">
              <div class="help-link-row__label">{{ isZh ? '开发环境指南' : 'Dev guide' }}</div>
              <el-link
                :href="devDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                docs/development_guide.md
              </el-link>
            </div>
            <div class="help-link-row">
              <div class="help-link-row__label">{{ isZh ? '飞书机器人通知' : 'Feishu bot' }}</div>
              <el-link
                :href="feishuDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                docs/飞书机器人Webhook通知教程.md
              </el-link>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="help-card help-card--accent" shadow="never">
          <template #header>
            <div class="help-card__header">
              <span>{{
                isZh ? '飞书机器人通知（推荐模板）' : 'Feishu Bot Notification (Recommended)'
              }}</span>
            </div>
          </template>
          <div class="help-feishu">
            <div class="help-feishu__desc">
              {{
                isZh
                  ? '系统管理 → 通知设置 → WebHook 通知：把模板改为交互卡片，消息更清爽；并优先使用 ResultSummary（避免把节点返回的 JSON 整段发到群里）。'
                  : 'System → Notification → WebHook: use interactive card template and prefer ResultSummary to avoid sending raw JSON to group.'
              }}
            </div>
            <div class="help-feishu__codeTitle">
              {{ isZh ? '交互卡片模板（可复制）' : 'Interactive card template' }}
            </div>
            <el-input
              class="help-code"
              type="textarea"
              :rows="10"
              :model-value="feishuInteractiveTemplate"
              readonly
            />
            <div class="help-feishu__actions">
              <el-link
                :href="feishuDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                {{ isZh ? '查看完整教程' : 'Open full guide' }}
              </el-link>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="help-grid">
      <el-col :xs="24" :md="12">
        <el-card class="help-card" shadow="never">
          <template #header>
            <div class="help-card__header">
              <span>{{ isZh ? '任务节点怎么填' : 'How to fill Task Node' }}</span>
            </div>
          </template>
          <el-descriptions :column="1" border class="help-desc">
            <el-descriptions-item :label="isZh ? '别名' : 'Alias'">
              {{
                isZh
                  ? '随意填写，建议写清机器用途，例如：NAS-1、Windows-DEV。'
                  : 'Any name you like, e.g. NAS-1, Windows-DEV.'
              }}
            </el-descriptions-item>
            <el-descriptions-item :label="isZh ? '主机名' : 'Host Name'">
              {{
                isZh
                  ? '填写运行 gocron-node 的机器地址（IP/域名）。如果 gocron 主控在 Docker/容器内运行，不要填 127.0.0.1，请填宿主机/局域网 IP。'
                  : 'Use the reachable IP/hostname of the machine running gocron-node. If gocron runs in Docker/container, do not use 127.0.0.1; use host/LAN IP.'
              }}
            </el-descriptions-item>
            <el-descriptions-item :label="isZh ? '端口' : 'Port'">
              {{
                isZh
                  ? '默认 5921（与 gocron-node 启动参数一致）。'
                  : 'Default 5921 (must match gocron-node listen port).'
              }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="help-card" shadow="never">
          <template #header>
            <div class="help-card__header">
              <span>{{ isZh ? 'Windows 本地开发常见问题' : 'Windows local dev FAQ' }}</span>
            </div>
          </template>
          <div class="help-text">
            <div v-if="isZh">
              <div class="help-text__item">
                1）如果主控在 Docker 里运行，节点主机名填 127.0.0.1 会失败（因为是容器自身）。
              </div>
              <div class="help-text__item">
                2）请确保 Windows 防火墙放行 5921/tcp，且 gocron-node 以 0.0.0.0:5921 监听。
              </div>
            </div>
            <div v-else>
              <div class="help-text__item">
                1) If gocron runs in Docker, using 127.0.0.1 for node host will fail (it points to
                container).
              </div>
              <div class="help-text__item">
                2) Allow inbound 5921/tcp on Windows firewall, and run gocron-node listening on
                0.0.0.0:5921.
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="help-card" shadow="never">
      <template #header>
        <div class="help-card__header">
          <span>{{ isZh ? 'Cron 示例' : 'Cron examples' }}</span>
        </div>
      </template>
      <el-descriptions :column="1" border class="help-desc">
        <el-descriptions-item :label="isZh ? '每 5 分钟' : 'Every 5 minutes'">
          <span class="help-mono">0 */5 * * * *</span>
        </el-descriptions-item>
        <el-descriptions-item :label="isZh ? '每天 02:30' : 'Daily 02:30'">
          <span class="help-mono">0 30 2 * * *</span>
        </el-descriptions-item>
        <el-descriptions-item :label="isZh ? '工作日 09:00' : 'Weekdays 09:00'">
          <span class="help-mono">0 0 9 * * 1-5</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'

export default {
  name: 'system-help',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale }
  },
  computed: {
    isZh() {
      return this.locale === 'zh-CN'
    },
    nasDocUrl() {
      return 'https://github.com/tabortao/gocron/blob/master/docs/nas_user_guide.md'
    },
    windowsDocUrl() {
      return 'https://github.com/tabortao/gocron/blob/master/docs/windows_user_guide.md'
    },
    devDocUrl() {
      return 'https://github.com/tabortao/gocron/blob/master/docs/development_guide.md'
    },
    feishuDocUrl() {
      return 'https://github.com/tabortao/gocron/blob/master/docs/%E9%A3%9E%E4%B9%A6%E6%9C%BA%E5%99%A8%E4%BA%BAWebhook%E9%80%9A%E7%9F%A5%E6%95%99%E7%A8%8B.md'
    },
    feishuInteractiveTemplate() {
      return `{
  "msg_type": "interactive",
  "card": {
    "config": { "wide_screen_mode": true },
    "header": {
      "template": "{{ if .IsSuccess }}green{{ else }}red{{ end }}",
      "title": { "tag": "plain_text", "content": "gocron 任务通知" }
    },
    "elements": [
      { "tag": "div", "text": { "tag": "lark_md", "content": "**任务**：{{.TaskName}}（ID: {{.TaskId}}）\\n**状态**：{{.StatusZh}}" } },
      { "tag": "div", "text": { "tag": "lark_md", "content": "**摘要**：{{.ResultSummary}}" } },
      { "tag": "hr" },
      { "tag": "div", "text": { "tag": "lark_md", "content": "**备注**：{{.Remark}}" } }
    ]
  }
}`
    }
  }
}
</script>

<style scoped>
.help-page {
  padding: 18px 18px 24px;
  background: linear-gradient(180deg, rgba(248, 250, 252, 0.9), rgba(255, 255, 255, 0.95));
  min-height: calc(100vh - 60px);
}

.help-hero {
  padding: 16px 16px 14px;
  border-radius: 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  background:
    radial-gradient(900px 240px at 10% 20%, rgba(37, 99, 235, 0.12), transparent 60%),
    radial-gradient(700px 220px at 90% 20%, rgba(249, 115, 22, 0.1), transparent 55%),
    rgba(255, 255, 255, 0.9);
  margin-bottom: 16px;
}

.help-hero__title {
  font-size: 18px;
  font-weight: 600;
  color: #0f172a;
  letter-spacing: -0.2px;
}

.help-hero__subtitle {
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.6;
  color: #475569;
}

.help-grid {
  margin-bottom: 16px;
}

.help-card {
  border-radius: 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  background: rgba(255, 255, 255, 0.98);
}

.help-card--accent {
  border-color: rgba(191, 219, 254, 1);
  background:
    radial-gradient(800px 220px at 50% 0%, rgba(37, 99, 235, 0.1), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card__header {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #0f172a;
  font-weight: 600;
}

.help-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.help-link-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
}

.help-link-row__label {
  color: #334155;
  font-size: 13px;
}

.help-feishu__desc {
  color: #475569;
  font-size: 13px;
  line-height: 1.7;
  margin-bottom: 10px;
}

.help-feishu__codeTitle {
  font-size: 12px;
  color: #334155;
  margin: 8px 0 8px;
}

.help-code :deep(textarea) {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
  font-size: 12px;
  line-height: 1.55;
}

.help-feishu__actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.help-text {
  color: #334155;
  line-height: 1.8;
}

.help-text__item {
  padding: 10px 12px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
}

.help-text__item + .help-text__item {
  margin-top: 10px;
}

.help-mono {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
}
</style>
