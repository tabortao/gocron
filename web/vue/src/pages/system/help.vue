<template>
  <el-main>
    <h3>{{ t('system.help') }}</h3>

    <el-alert
      v-if="isZh"
      type="info"
      :closable="false"
      show-icon
      title="快速定位"
      style="margin-bottom: 15px"
    >
      <div>
        NAS 教程：
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/nas_user_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/nas_user_guide.md
        </a>
      </div>
      <div>
        Windows 教程：
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/windows_user_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/windows_user_guide.md
        </a>
      </div>
      <div>
        开发环境指南：
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/development_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/development_guide.md
        </a>
      </div>
    </el-alert>

    <el-alert
      v-else
      type="info"
      :closable="false"
      show-icon
      title="Quick links"
      style="margin-bottom: 15px"
    >
      <div>
        NAS guide:
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/nas_user_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/nas_user_guide.md
        </a>
      </div>
      <div>
        Windows guide:
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/windows_user_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/windows_user_guide.md
        </a>
      </div>
      <div>
        Dev guide:
        <a
          href="https://github.com/tabortao/gocron/blob/main/docs/development_guide.md"
          target="_blank"
          rel="noopener noreferrer"
        >
          docs/development_guide.md
        </a>
      </div>
    </el-alert>

    <el-card style="margin-bottom: 15px">
      <template #header>
        <span>{{ isZh ? '任务节点怎么填' : 'How to fill Task Node' }}</span>
      </template>
      <el-descriptions :column="1" border>
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

    <el-card style="margin-bottom: 15px">
      <template #header>
        <span>{{ isZh ? 'Windows 本地开发常见问题' : 'Windows local dev FAQ' }}</span>
      </template>
      <div v-if="isZh" style="line-height: 1.8">
        <div>1）如果主控在 Docker 里运行，节点主机名填 127.0.0.1 会失败（因为是容器自身）。</div>
        <div>2）请确保 Windows 防火墙放行 5921/tcp，且 gocron-node 以 0.0.0.0:5921 监听。</div>
      </div>
      <div v-else style="line-height: 1.8">
        <div>
          1) If gocron runs in Docker, using 127.0.0.1 for node host will fail (it points to
          container).
        </div>
        <div>
          2) Allow inbound 5921/tcp on Windows firewall, and run gocron-node listening on
          0.0.0.0:5921.
        </div>
      </div>
    </el-card>

    <el-card>
      <template #header>
        <span>{{ isZh ? 'Cron 示例' : 'Cron examples' }}</span>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="isZh ? '每 5 分钟' : 'Every 5 minutes'">
          <span style="font-family: monospace">0 */5 * * * *</span>
        </el-descriptions-item>
        <el-descriptions-item :label="isZh ? '每天 02:30' : 'Daily 02:30'">
          <span style="font-family: monospace">0 30 2 * * *</span>
        </el-descriptions-item>
        <el-descriptions-item :label="isZh ? '工作日 09:00' : 'Weekdays 09:00'">
          <span style="font-family: monospace">0 0 9 * * 1-5</span>
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
    }
  }
}
</script>
