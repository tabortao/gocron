<template>
  <el-main class="help-page">
    <div class="help-hero">
      <div class="help-hero__icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" />
          <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
          <line x1="12" y1="17" x2="12.01" y2="17" />
        </svg>
      </div>
      <div class="help-hero__content">
        <div class="help-hero__title">{{ t('system.help') }}</div>
        <div class="help-hero__subtitle">
          {{
            isZh
              ? '常见配置与排错入口，建议先从「快速入口」和「飞书机器人通知」开始。'
              : 'Common setup and troubleshooting. Start with Quick links and Feishu bot notifications.'
          }}
        </div>
      </div>
    </div>

    <el-row :gutter="16" class="help-grid">
      <el-col :xs="24" :md="12">
        <el-card class="help-card help-card--links" shadow="never">
          <template #header>
            <div class="help-card__header">
              <div class="help-card__icon help-card__icon--blue">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
                  <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
                </svg>
              </div>
              <span>{{ isZh ? '快速入口' : 'Quick links' }}</span>
            </div>
          </template>
          <div class="help-links">
            <a :href="nasDocUrl" target="_blank" rel="noopener noreferrer" class="help-link-item">
              <div class="help-link-item__icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
                  <line x1="8" y1="21" x2="16" y2="21" />
                  <line x1="12" y1="17" x2="12" y2="21" />
                </svg>
              </div>
              <div class="help-link-item__content">
                <div class="help-link-item__label">{{ isZh ? 'NAS 教程' : 'NAS guide' }}</div>
                <div class="help-link-item__path">docs/nas_user_guide.md</div>
              </div>
              <svg
                class="help-link-item__arrow"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                <polyline points="15 3 21 3 21 9" />
                <line x1="10" y1="14" x2="21" y2="3" />
              </svg>
            </a>
            <a
              :href="windowsDocUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="help-link-item"
            >
              <div class="help-link-item__icon help-link-item__icon--green">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="3" width="20" height="14" rx="2" ry="2" />
                  <path d="M8 21h8" />
                  <path d="M12 17v4" />
                  <path d="M7 8l3 3-3 3" />
                  <path d="M13 14h4" />
                </svg>
              </div>
              <div class="help-link-item__content">
                <div class="help-link-item__label">
                  {{ isZh ? 'Windows 教程' : 'Windows guide' }}
                </div>
                <div class="help-link-item__path">docs/windows_user_guide.md</div>
              </div>
              <svg
                class="help-link-item__arrow"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                <polyline points="15 3 21 3 21 9" />
                <line x1="10" y1="14" x2="21" y2="3" />
              </svg>
            </a>
            <a :href="devDocUrl" target="_blank" rel="noopener noreferrer" class="help-link-item">
              <div class="help-link-item__icon help-link-item__icon--purple">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="16 18 22 12 16 6" />
                  <polyline points="8 6 2 12 8 18" />
                </svg>
              </div>
              <div class="help-link-item__content">
                <div class="help-link-item__label">{{ isZh ? '开发环境指南' : 'Dev guide' }}</div>
                <div class="help-link-item__path">docs/development_guide.md</div>
              </div>
              <svg
                class="help-link-item__arrow"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                <polyline points="15 3 21 3 21 9" />
                <line x1="10" y1="14" x2="21" y2="3" />
              </svg>
            </a>
            <a
              :href="feishuDocUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="help-link-item"
            >
              <div class="help-link-item__icon help-link-item__icon--orange">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                </svg>
              </div>
              <div class="help-link-item__content">
                <div class="help-link-item__label">
                  {{ isZh ? '飞书机器人通知' : 'Feishu bot' }}
                </div>
                <div class="help-link-item__path">docs/飞书机器人Webhook通知教程.md</div>
              </div>
              <svg
                class="help-link-item__arrow"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                <polyline points="15 3 21 3 21 9" />
                <line x1="10" y1="14" x2="21" y2="3" />
              </svg>
            </a>
            <a
              :href="serverChan3DocUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="help-link-item"
            >
              <div class="help-link-item__icon help-link-item__icon--pink">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                  <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                </svg>
              </div>
              <div class="help-link-item__content">
                <div class="help-link-item__label">
                  {{ isZh ? 'Server 酱³ 通知' : 'Server Chan' }}
                </div>
                <div class="help-link-item__path">docs/Server酱³通知教程.md</div>
              </div>
              <svg
                class="help-link-item__arrow"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                <polyline points="15 3 21 3 21 9" />
                <line x1="10" y1="14" x2="21" y2="3" />
              </svg>
            </a>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="help-card help-card--feishu" shadow="never">
          <template #header>
            <div class="help-card__header">
              <div class="help-card__icon help-card__icon--orange">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                </svg>
              </div>
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
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="16 18 22 12 16 6" />
                <polyline points="8 6 2 12 8 18" />
              </svg>
              {{ isZh ? '交互卡片模板（可复制）' : 'Interactive card template' }}
            </div>
            <div class="help-code-wrapper">
              <el-input
                class="help-code"
                type="textarea"
                :rows="10"
                :model-value="feishuInteractiveTemplate"
                readonly
              />
              <el-button
                class="help-code__copy"
                type="primary"
                size="small"
                @click="copyTemplate('feishu')"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  style="width: 14px; height: 14px; margin-right: 4px"
                >
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                </svg>
                {{ isZh ? '复制' : 'Copy' }}
              </el-button>
            </div>
            <div class="help-feishu__actions">
              <el-link
                :href="feishuDocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                {{ isZh ? '查看完整教程' : 'Open full guide' }}
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  style="width: 14px; height: 14px; margin-left: 4px"
                >
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                  <polyline points="15 3 21 3 21 9" />
                  <line x1="10" y1="14" x2="21" y2="3" />
                </svg>
              </el-link>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="help-grid">
      <el-col :xs="24" :md="24">
        <el-card class="help-card help-card--serverchan" shadow="never">
          <template #header>
            <div class="help-card__header">
              <div class="help-card__icon help-card__icon--pink">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                  <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                </svg>
              </div>
              <span>{{ isZh ? 'Server 酱³ 通知（推荐模板）' : 'Server Chan (Recommended)' }}</span>
            </div>
          </template>
          <div class="help-feishu">
            <div class="help-feishu__desc">
              {{
                isZh
                  ? '系统管理 → 通知设置 → Server 酱³：配置 API URL 后，使用 title + desp 模板（application/x-www-form-urlencoded）。desp 支持 Markdown，适合做任务日程提醒。'
                  : 'System → Notification → Server Chan: configure API URL and use title + desp templates (application/x-www-form-urlencoded).'
              }}
            </div>
            <div class="help-templates-grid">
              <div class="help-template-item">
                <div class="help-feishu__codeTitle">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="16 18 22 12 16 6" />
                    <polyline points="8 6 2 12 8 18" />
                  </svg>
                  {{ isZh ? '标题模板（title）' : 'Title template' }}
                </div>
                <div class="help-code-wrapper">
                  <el-input
                    class="help-code"
                    type="textarea"
                    :rows="2"
                    :model-value="serverChan3TitleTemplate"
                    readonly
                  />
                  <el-button
                    class="help-code__copy"
                    type="primary"
                    size="small"
                    @click="copyTemplate('title')"
                  >
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      style="width: 14px; height: 14px; margin-right: 4px"
                    >
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                    </svg>
                    {{ isZh ? '复制' : 'Copy' }}
                  </el-button>
                </div>
              </div>
              <div class="help-template-item">
                <div class="help-feishu__codeTitle">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="16 18 22 12 16 6" />
                    <polyline points="8 6 2 12 8 18" />
                  </svg>
                  {{ isZh ? '内容模板（desp，可复制）' : 'Desp template' }}
                </div>
                <div class="help-code-wrapper">
                  <el-input
                    class="help-code"
                    type="textarea"
                    :rows="9"
                    :model-value="serverChan3DespTemplate"
                    readonly
                  />
                  <el-button
                    class="help-code__copy"
                    type="primary"
                    size="small"
                    @click="copyTemplate('desp')"
                  >
                    <svg
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      style="width: 14px; height: 14px; margin-right: 4px"
                    >
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
                    </svg>
                    {{ isZh ? '复制' : 'Copy' }}
                  </el-button>
                </div>
              </div>
            </div>
            <div class="help-feishu__actions">
              <el-link
                :href="serverChan3DocUrl"
                target="_blank"
                rel="noopener noreferrer"
                type="primary"
                :underline="false"
              >
                {{ isZh ? '查看完整教程' : 'Open full guide' }}
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  style="width: 14px; height: 14px; margin-left: 4px"
                >
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                  <polyline points="15 3 21 3 21 9" />
                  <line x1="10" y1="14" x2="21" y2="3" />
                </svg>
              </el-link>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="help-grid">
      <el-col :xs="24" :md="12">
        <el-card class="help-card help-card--node" shadow="never">
          <template #header>
            <div class="help-card__header">
              <div class="help-card__icon help-card__icon--green">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="2" y="2" width="20" height="8" rx="2" ry="2" />
                  <rect x="2" y="14" width="20" height="8" rx="2" ry="2" />
                  <line x1="6" y1="6" x2="6.01" y2="6" />
                  <line x1="6" y1="18" x2="6.01" y2="18" />
                </svg>
              </div>
              <span>{{ isZh ? '任务节点怎么填' : 'How to fill Task Node' }}</span>
            </div>
          </template>
          <div class="help-desc-list">
            <div class="help-desc-item">
              <div class="help-desc-item__label">{{ isZh ? '别名' : 'Alias' }}</div>
              <div class="help-desc-item__value">
                {{
                  isZh
                    ? '随意填写，建议写清机器用途，例如：NAS-1、Windows-DEV。'
                    : 'Any name you like, e.g. NAS-1, Windows-DEV.'
                }}
              </div>
            </div>
            <div class="help-desc-item">
              <div class="help-desc-item__label">{{ isZh ? '主机名' : 'Host Name' }}</div>
              <div class="help-desc-item__value">
                {{
                  isZh
                    ? '填写运行 gocron-node 的机器地址（IP/域名）。如果 gocron 主控在 Docker/容器内运行，不要填 127.0.0.1，请填宿主机/局域网 IP。'
                    : 'Use the reachable IP/hostname of the machine running gocron-node. If gocron runs in Docker/container, do not use 127.0.0.1; use host/LAN IP.'
                }}
              </div>
            </div>
            <div class="help-desc-item">
              <div class="help-desc-item__label">{{ isZh ? '端口' : 'Port' }}</div>
              <div class="help-desc-item__value">
                {{
                  isZh
                    ? '默认 5921（与 gocron-node 启动参数一致）。'
                    : 'Default 5921 (must match gocron-node listen port).'
                }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="help-card help-card--faq" shadow="never">
          <template #header>
            <div class="help-card__header">
              <div class="help-card__icon help-card__icon--red">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="8" x2="12" y2="12" />
                  <line x1="12" y1="16" x2="12.01" y2="16" />
                </svg>
              </div>
              <span>{{ isZh ? 'Windows 本地开发常见问题' : 'Windows local dev FAQ' }}</span>
            </div>
          </template>
          <div class="help-faq-list">
            <div v-if="isZh">
              <div class="help-faq-item">
                <div class="help-faq-item__num">1</div>
                <div class="help-faq-item__text">
                  如果主控在 Docker 里运行，节点主机名填 127.0.0.1 会失败（因为是容器自身）。
                </div>
              </div>
              <div class="help-faq-item">
                <div class="help-faq-item__num">2</div>
                <div class="help-faq-item__text">
                  请确保 Windows 防火墙放行 5921/tcp，且 gocron-node 以 0.0.0.0:5921 监听。
                </div>
              </div>
            </div>
            <div v-else>
              <div class="help-faq-item">
                <div class="help-faq-item__num">1</div>
                <div class="help-faq-item__text">
                  If gocron runs in Docker, using 127.0.0.1 for node host will fail (it points to
                  container).
                </div>
              </div>
              <div class="help-faq-item">
                <div class="help-faq-item__num">2</div>
                <div class="help-faq-item__text">
                  Allow inbound 5921/tcp on Windows firewall, and run gocron-node listening on
                  0.0.0.0:5921.
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="help-card help-card--cron" shadow="never">
      <template #header>
        <div class="help-card__header">
          <div class="help-card__icon help-card__icon--cyan">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
          </div>
          <span>{{ isZh ? 'Cron 示例' : 'Cron examples' }}</span>
        </div>
      </template>
      <div class="help-cron-list">
        <div class="help-cron-item">
          <div class="help-cron-item__label">{{ isZh ? '每 5 分钟' : 'Every 5 minutes' }}</div>
          <div class="help-cron-item__code">
            <span class="help-mono">0 */5 * * * *</span>
            <el-button size="small" text @click="copyCron('0 */5 * * * *')">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                style="width: 14px; height: 14px"
              >
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
            </el-button>
          </div>
        </div>
        <div class="help-cron-item">
          <div class="help-cron-item__label">{{ isZh ? '每天 02:30' : 'Daily 02:30' }}</div>
          <div class="help-cron-item__code">
            <span class="help-mono">0 30 2 * * *</span>
            <el-button size="small" text @click="copyCron('0 30 2 * * *')">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                style="width: 14px; height: 14px"
              >
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
            </el-button>
          </div>
        </div>
        <div class="help-cron-item">
          <div class="help-cron-item__label">{{ isZh ? '工作日 09:00' : 'Weekdays 09:00' }}</div>
          <div class="help-cron-item__code">
            <span class="help-mono">0 0 9 * * 1-5</span>
            <el-button size="small" text @click="copyCron('0 0 9 * * 1-5')">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                style="width: 14px; height: 14px"
              >
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'
import { copyText } from '../../utils/clipboard'

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
    serverChan3DocUrl() {
      return 'https://github.com/tabortao/gocron/blob/master/docs/Server%E9%85%B1%C2%B3%E9%80%9A%E7%9F%A5%E6%95%99%E7%A8%8B.md'
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
    },
    serverChan3TitleTemplate() {
      return '{{.TaskName}} - {{.StatusZh}}'
    },
    serverChan3DespTemplate() {
      return `**任务**：{{.TaskName}}（ID: {{.TaskId}}）

**状态**：{{.StatusZh}}

{{ if .Host }}**节点**：{{.Host}}

{{ end }}**摘要**：{{.ResultSummary}}

{{ if .Remark }}**备注**：{{.Remark}}{{ end }}`
    }
  },
  methods: {
    copyTemplate(type) {
      let text = ''
      if (type === 'feishu') {
        text = this.feishuInteractiveTemplate
      } else if (type === 'title') {
        text = this.serverChan3TitleTemplate
      } else if (type === 'desp') {
        text = this.serverChan3DespTemplate
      }
      copyText(text)
        .then(() => {
          this.$message.success(this.isZh ? '复制成功' : 'Copied!')
        })
        .catch(() => {
          this.$message.error(this.isZh ? '复制失败' : 'Copy failed')
        })
    },
    copyCron(cron) {
      copyText(cron)
        .then(() => {
          this.$message.success(this.isZh ? '复制成功' : 'Copied!')
        })
        .catch(() => {
          this.$message.error(this.isZh ? '复制失败' : 'Copy failed')
        })
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
  padding: 20px 20px 18px;
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 1);
  background:
    radial-gradient(900px 240px at 10% 20%, rgba(99, 102, 241, 0.15), transparent 60%),
    radial-gradient(700px 220px at 90% 20%, rgba(236, 72, 153, 0.1), transparent 55%),
    rgba(255, 255, 255, 0.95);
  margin-bottom: 18px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.help-hero__icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.help-hero__icon svg {
  width: 24px;
  height: 24px;
  color: white;
}

.help-hero__content {
  flex: 1;
}

.help-hero__title {
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: -0.3px;
}

.help-hero__subtitle {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.7;
  color: #64748b;
}

.help-grid {
  margin-bottom: 16px;
}

.help-card {
  border-radius: 16px;
  border: 1px solid rgba(226, 232, 240, 1);
  background: rgba(255, 255, 255, 0.98);
  transition: all 0.3s ease;
  height: 100%;
}

.help-card:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
  transform: translateY(-2px);
}

.help-card--links {
  background:
    radial-gradient(600px 200px at 0% 0%, rgba(59, 130, 246, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card--feishu {
  background:
    radial-gradient(600px 200px at 0% 0%, rgba(249, 115, 22, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card--serverchan {
  background:
    radial-gradient(800px 200px at 50% 0%, rgba(236, 72, 153, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card--node {
  background:
    radial-gradient(600px 200px at 0% 0%, rgba(34, 197, 94, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card--faq {
  background:
    radial-gradient(600px 200px at 0% 0%, rgba(239, 68, 68, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card--cron {
  background:
    radial-gradient(800px 200px at 50% 0%, rgba(6, 182, 212, 0.08), transparent 55%),
    rgba(255, 255, 255, 0.98);
}

.help-card__header {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #0f172a;
  font-weight: 600;
  font-size: 15px;
}

.help-card__icon {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.help-card__icon svg {
  width: 18px;
  height: 18px;
}

.help-card__icon--blue {
  background: linear-gradient(135deg, #3b82f6, #60a5fa);
  color: white;
}

.help-card__icon--orange {
  background: linear-gradient(135deg, #f97316, #fb923c);
  color: white;
}

.help-card__icon--pink {
  background: linear-gradient(135deg, #ec4899, #f472b6);
  color: white;
}

.help-card__icon--green {
  background: linear-gradient(135deg, #22c55e, #4ade80);
  color: white;
}

.help-card__icon--purple {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
  color: white;
}

.help-card__icon--red {
  background: linear-gradient(135deg, #ef4444, #f87171);
  color: white;
}

.help-card__icon--cyan {
  background: linear-gradient(135deg, #06b6d4, #22d3ee);
  color: white;
}

.help-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.help-link-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;
}

.help-link-item:hover {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.05);
  transform: translateX(4px);
}

.help-link-item__icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #3b82f6, #60a5fa);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.help-link-item__icon svg {
  width: 18px;
  height: 18px;
  color: white;
}

.help-link-item__icon--green {
  background: linear-gradient(135deg, #22c55e, #4ade80);
}

.help-link-item__icon--purple {
  background: linear-gradient(135deg, #8b5cf6, #a78bfa);
}

.help-link-item__icon--orange {
  background: linear-gradient(135deg, #f97316, #fb923c);
}

.help-link-item__icon--pink {
  background: linear-gradient(135deg, #ec4899, #f472b6);
}

.help-link-item__content {
  flex: 1;
  min-width: 0;
}

.help-link-item__label {
  color: #1e293b;
  font-size: 14px;
  font-weight: 500;
}

.help-link-item__path {
  color: #94a3b8;
  font-size: 12px;
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.help-link-item__arrow {
  width: 16px;
  height: 16px;
  color: #94a3b8;
  flex-shrink: 0;
}

.help-feishu__desc {
  color: #64748b;
  font-size: 13px;
  line-height: 1.8;
  margin-bottom: 14px;
  padding: 12px 14px;
  background: rgba(248, 250, 252, 0.8);
  border-radius: 10px;
  border-left: 3px solid #f97316;
}

.help-feishu__codeTitle {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #475569;
  margin: 12px 0 10px;
  font-weight: 500;
}

.help-feishu__codeTitle svg {
  width: 16px;
  height: 16px;
  color: #6366f1;
}

.help-code-wrapper {
  position: relative;
}

.help-code :deep(textarea) {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
  font-size: 12px;
  line-height: 1.6;
  background: #1e293b !important;
  color: #e2e8f0 !important;
  border-radius: 10px;
  padding: 12px;
}

.help-code__copy {
  position: absolute;
  top: 8px;
  right: 8px;
}

.help-templates-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.help-template-item {
  min-width: 0;
}

.help-feishu__actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 14px;
}

.help-feishu__actions .el-link {
  display: flex;
  align-items: center;
  font-weight: 500;
}

.help-desc-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.help-desc-item {
  padding: 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
}

.help-desc-item__label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
  margin-bottom: 6px;
}

.help-desc-item__value {
  font-size: 13px;
  color: #334155;
  line-height: 1.7;
}

.help-faq-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.help-faq-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
}

.help-faq-item__num {
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: linear-gradient(135deg, #ef4444, #f87171);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.help-faq-item__text {
  font-size: 13px;
  color: #334155;
  line-height: 1.7;
}

.help-cron-list {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.help-cron-item {
  padding: 14px;
  border: 1px solid rgba(226, 232, 240, 1);
  border-radius: 12px;
  background: rgba(248, 250, 252, 0.6);
  text-align: center;
}

.help-cron-item__label {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 8px;
}

.help-cron-item__code {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.help-mono {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
  font-size: 14px;
  color: #0891b2;
  background: rgba(6, 182, 212, 0.1);
  padding: 4px 10px;
  border-radius: 6px;
}

@media screen and (max-width: 1024px) {
  .help-templates-grid {
    grid-template-columns: 1fr;
  }

  .help-cron-list {
    grid-template-columns: 1fr;
  }
}

@media screen and (max-width: 768px) {
  .help-page {
    padding: 12px;
  }

  .help-hero {
    padding: 16px;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .help-hero__icon {
    width: 56px;
    height: 56px;
  }

  .help-hero__icon svg {
    width: 28px;
    height: 28px;
  }

  .help-hero__title {
    font-size: 18px;
  }

  .help-hero__subtitle {
    font-size: 13px;
  }

  .help-grid {
    margin-bottom: 12px;
  }

  .help-card {
    margin-bottom: 0;
  }

  .help-card:hover {
    transform: none;
  }

  .help-link-item {
    padding: 10px 12px;
  }

  .help-link-item__icon {
    width: 32px;
    height: 32px;
  }

  .help-link-item__icon svg {
    width: 16px;
    height: 16px;
  }

  .help-link-item__label {
    font-size: 13px;
  }

  .help-link-item__path {
    display: none;
  }

  .help-feishu__desc {
    font-size: 12px;
    padding: 10px 12px;
  }

  .help-code :deep(textarea) {
    font-size: 11px;
    padding: 10px;
  }

  .help-desc-item {
    padding: 12px;
  }

  .help-desc-item__value {
    font-size: 12px;
  }

  .help-faq-item {
    padding: 12px;
  }

  .help-faq-item__text {
    font-size: 12px;
  }

  .help-cron-item {
    padding: 12px;
  }

  .help-mono {
    font-size: 12px;
  }
}
</style>
