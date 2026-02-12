<template>
  <el-card class="card-section" shadow="never">
    <el-tabs v-model="activeName">
      <el-tab-pane :label="t('system.email')" name="email"></el-tab-pane>
      <el-tab-pane label="Slack" name="slack"></el-tab-pane>
      <el-tab-pane label="Webhook" name="webhook"></el-tab-pane>
      <el-tab-pane label="Server 酱³" name="serverchan3"></el-tab-pane>
    </el-tabs>
    <el-alert type="info" :closable="false" style="margin-bottom: 15px">
      <template #title>
        <div style="font-weight: 600; margin-bottom: 8px">{{ t('system.templateVariables') }}</div>
        <div style="font-size: 13px; line-height: 1.8">
          <div>
            <code v-pre>{{.TaskId}}</code> - {{ t('system.taskIdVar') }}
          </div>
          <div>
            <code v-pre>{{.TaskName}}</code> - {{ t('system.taskNameVar') }}
          </div>
          <div>
            <code v-pre>{{.Status}}</code> - {{ t('system.statusVar') }}
          </div>
          <div>
            <code v-pre>{{.StatusZh}}</code> - {{ t('system.statusVar') }}（中文）
          </div>
          <div>
            <code v-pre>{{.IsSuccess}}</code> - {{ t('system.statusVar') }}（true/false）
          </div>
          <div>
            <code v-pre>{{.Host}}</code> - 节点信息（若输出包含 Host 行）
          </div>
          <div>
            <code v-pre>{{.Result}}</code> - {{ t('system.resultVar') }}
          </div>
          <div>
            <code v-pre>{{.ResultSummary}}</code> - 输出摘要（优先 JSON.message）
          </div>
          <div>
            <code v-pre>{{.ResultBody}}</code> - 去掉 Host 行后的主体输出
          </div>
          <div>
            <code v-pre>{{.Remark}}</code> - {{ t('task.remark') }}
          </div>
        </div>
      </template>
    </el-alert>
  </el-card>
</template>

<script>
import { useI18n } from 'vue-i18n'
export default {
  name: 'notification-tab',
  setup() {
    const { t } = useI18n()
    return { t }
  },
  data() {
    return {
      activeName: ''
    }
  },
  created() {
    const segments = this.$route.path.split('/')
    if (segments.length !== 4) {
      this.activeName = 'email'
      return
    }
    this.activeName = segments[3]
  },
  watch: {
    activeName(newVal) {
      if (newVal && this.$route.path !== `/system/notification/${newVal}`) {
        this.$router.push(`/system/notification/${newVal}`)
      }
    },
    '$route.path': {
      handler(newPath) {
        const segments = newPath.split('/')
        if (segments.length === 4 && segments[2] === 'notification') {
          this.activeName = segments[3]
        }
      },
      immediate: false
    }
  }
}
</script>
