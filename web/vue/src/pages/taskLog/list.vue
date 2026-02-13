<template>
  <el-main>
    <div class="page-header">
      <div class="page-title">{{ t('taskLog.list') }}</div>
      <div class="toolbar">
        <el-button type="danger" v-if="isAdmin" @click="clearLog">{{
          t('message.clearLog')
        }}</el-button>
        <el-button type="info" @click="refresh">{{ t('common.refresh') }}</el-button>
      </div>
    </div>

    <el-card class="card-section filter-card" shadow="never">
      <el-form :inline="true" size="small">
        <el-form-item :label="t('task.id')">
          <el-input v-model.trim="searchParams.task_id" style="width: 200px" clearable></el-input>
        </el-form-item>
        <el-form-item :label="t('task.protocol')">
          <el-select
            v-model.trim="searchParams.protocol"
            :placeholder="t('task.protocol')"
            style="width: 200px"
            clearable
          >
            <el-option :label="t('message.all')" value=""></el-option>
            <el-option
              v-for="item in protocolList"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.status')">
          <el-select v-model.trim="searchParams.status" style="width: 200px" clearable>
            <el-option :label="t('message.all')" value=""></el-option>
            <el-option
              v-for="item in statusList"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search()">{{ t('common.search') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="card-section table-card" shadow="never">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="logTotal"
        v-model:current-page="searchParams.page"
        v-model:page-size="searchParams.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
      >
      </el-pagination>
      <el-table :data="logs" border ref="table" style="width: 100%">
        <el-table-column type="expand">
          <template #default="scope">
            <el-form label-position="left">
              <el-form-item>
                {{ t('message.retryCount') }}: {{ scope.row.retry_times }} <br />
                {{ t('task.cronExpression') }}: {{ scope.row.spec }} <br />
                {{ t('task.command') }}: {{ scope.row.command }}
              </el-form-item>
            </el-form>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID"> </el-table-column>
        <el-table-column prop="task_id" :label="t('task.id')"> </el-table-column>
        <el-table-column prop="name" :label="t('task.name')" width="180"> </el-table-column>
        <el-table-column prop="protocol" :label="t('task.protocol')" :formatter="formatProtocol">
        </el-table-column>
        <el-table-column :label="t('task.taskNode')" width="150">
          <template #default="scope">
            <div v-html="scope.row.hostname"></div>
          </template>
        </el-table-column>
        <el-table-column :label="t('taskLog.duration')" width="250">
          <template #default="scope">
            {{ t('taskLog.duration') }}: {{ scope.row.total_time > 0 ? scope.row.total_time : 1
            }}{{ t('message.seconds') }}<br />
            {{ t('taskLog.startTime') }}: {{ $filters.formatTime(scope.row.start_time) }}<br />
            <span v-if="scope.row.status !== 1"
              >{{ t('taskLog.endTime') }}: {{ $filters.formatTime(scope.row.end_time) }}</span
            >
          </template>
        </el-table-column>
        <el-table-column :label="t('common.status')">
          <template #default="scope">
            <span style="color: red" v-if="scope.row.status === 0">{{ t('taskLog.failed') }}</span>
            <span style="color: green" v-else-if="scope.row.status === 1">{{
              t('message.running')
            }}</span>
            <span v-else-if="scope.row.status === 2">{{ t('taskLog.success') }}</span>
            <span style="color: #4499ee" v-else-if="scope.row.status === 3">{{
              t('message.cancelled')
            }}</span>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('taskLog.result')"
          :width="locale === availableLanguages.zhCN.value ? 120 : 140"
          v-if="isAdmin"
        >
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              v-if="scope.row.status === 1"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="success"
              size="small"
              v-if="scope.row.status === 2"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="warning"
              size="small"
              v-if="scope.row.status === 0"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="info"
              size="small"
              v-if="scope.row.status === 3"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="danger"
              size="small"
              v-if="scope.row.status === 1 && scope.row.protocol === 2"
              @click="stopTask(scope.row)"
              >{{ t('message.stopTask') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          :label="t('taskLog.result')"
          :width="locale === availableLanguages.zhCN.value ? 120 : 140"
          v-else
        >
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              v-if="scope.row.status === 1"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="success"
              size="small"
              v-if="scope.row.status === 2"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="warning"
              size="small"
              v-if="scope.row.status === 0"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
            <el-button
              type="info"
              size="small"
              v-if="scope.row.status === 3"
              @click="showTaskResult(scope.row)"
              >{{ t('taskLog.viewOutput') }}</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog :title="t('message.taskExecutionResult')" v-model="dialogVisible" width="60%">
      <div v-if="currentTaskResult.hostname">
        <strong>{{ t('taskLog.host') }}:</strong>
        <pre v-html="currentTaskResult.hostname"></pre>
      </div>
      <div>
        <strong>{{ t('task.command') }}:</strong>
        <pre>{{ currentTaskResult.command }}</pre>
      </div>
      <div>
        <strong>{{ t('taskLog.output') }}:</strong>
        <pre ref="resultPre" style="max-height: 50vh; overflow: auto">{{
          currentTaskResult.result
        }}</pre>
      </div>
    </el-dialog>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import taskLogService from '../../api/taskLog'
import { useUserStore } from '../../stores/user'
import { availableLanguages } from '@/const/index'

export default {
  name: 'task-log',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale, availableLanguages }
  },
  data() {
    const userStore = useUserStore()
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        task_id: '',
        protocol: '',
        status: ''
      },
      isAdmin: userStore.isAdmin,
      dialogVisible: false,
      currentTaskResult: {
        hostname: '',
        command: '',
        result: ''
      },
      currentLogId: 0,
      currentLogStatus: 0,
      outputRefreshTimer: null,
      outputRefreshInFlight: false,
      protocolList: [
        {
          value: '1',
          label: 'http'
        },
        {
          value: '2',
          label: 'shell'
        }
      ],
      statusList: [],
      autoRefreshTimer: null,
      autoRefreshInFlight: false
    }
  },
  computed: {
    computedStatusList() {
      return [
        { value: '3', label: this.t('taskLog.success') },
        { value: '1', label: this.t('taskLog.failed') },
        { value: '4', label: this.t('message.cancelled') }
      ]
    }
  },
  watch: {
    computedStatusList: {
      handler(newVal) {
        this.statusList = newVal
      },
      immediate: true
    },
    '$route.query.task_id': {
      handler(newTaskId) {
        if (newTaskId !== undefined) {
          this.searchParams.task_id = newTaskId
          this.searchParams.page = 1
          this.search()
        }
      }
    },
    dialogVisible: {
      handler(visible) {
        if (!visible) {
          this.stopOutputRefresh()
        }
      }
    }
  },
  created() {
    this.updateTaskIdFromRoute()
    this.search()
  },
  activated() {
    this.updateTaskIdFromRoute()
    this.search()
  },
  deactivated() {
    this.stopAutoRefresh()
    this.stopOutputRefresh()
  },
  beforeUnmount() {
    this.stopAutoRefresh()
    this.stopOutputRefresh()
  },
  methods: {
    formatProtocol(row, col) {
      if (row[col.property] === 1) {
        return 'http'
      }
      return 'shell'
    },
    changePage(page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize(pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search(callback = null) {
      taskLogService.list(this.searchParams, data => {
        this.logs = data.data
        this.logTotal = data.total
        this.ensureAutoRefresh()

        if (callback) {
          callback()
        }
      })
    },
    clearLog() {
      ElMessageBox.confirm(this.t('message.confirmClearLog'), this.t('common.tip'), {
        confirmButtonText: this.t('common.confirm'),
        cancelButtonText: this.t('common.cancel'),
        type: 'warning',
        center: true
      })
        .then(() => {
          taskLogService.clear(() => {
            this.searchParams.page = 1
            this.search()
          })
        })
        .catch(() => {})
    },
    stopTask(item) {
      taskLogService.stop(item.id, item.task_id, () => {
        this.search()
      })
    },
    showTaskResult(item) {
      this.dialogVisible = true
      this.currentLogId = item.id || 0
      this.currentLogStatus = item.status
      // 清理命令中的 HTML 实体编码
      let cleanedCommand = item.command
      if (cleanedCommand) {
        cleanedCommand = cleanedCommand
          .replace(/&quot;/g, '"')
          .replace(/&apos;/g, "'")
          .replace(/&#39;/g, "'")
          .replace(/&lt;/g, '<')
          .replace(/&gt;/g, '>')
          .replace(/&amp;/g, '&')
      }
      this.currentTaskResult.hostname = item.hostname || ''
      this.currentTaskResult.command = cleanedCommand
      this.currentTaskResult.result = item.result
      if (item.status === 1) {
        this.fetchLiveOutput()
        this.startOutputRefresh()
      } else {
        this.stopOutputRefresh()
      }
      this.$nextTick(() => {
        this.scrollOutputToBottom()
      })
    },
    refresh() {
      this.search(() => {
        this.$message.success(this.t('message.refreshSuccess'))
      })
    },
    updateTaskIdFromRoute() {
      if (this.$route.query.task_id) {
        this.searchParams.task_id = this.$route.query.task_id
        this.searchParams.page = 1
      }
    },
    fetchLiveOutput() {
      if (!this.currentLogId || this.outputRefreshInFlight) return
      this.outputRefreshInFlight = true
      taskLogService.output(this.currentLogId, data => {
        this.outputRefreshInFlight = false
        if (data && data.status !== undefined) {
          this.currentLogStatus = data.status
          if (this.currentLogStatus !== 1) {
            this.stopOutputRefresh()
          }
        }
        if (data && data.output !== undefined) {
          this.currentTaskResult.result = data.output
          this.$nextTick(() => {
            this.scrollOutputToBottom()
          })
        }
      })
    },
    startOutputRefresh() {
      if (this.outputRefreshTimer || this.currentLogStatus !== 1) return
      this.outputRefreshTimer = setInterval(() => {
        if (!this.dialogVisible || this.currentLogStatus !== 1) return
        this.fetchLiveOutput()
      }, 2000)
    },
    stopOutputRefresh() {
      if (!this.outputRefreshTimer) return
      clearInterval(this.outputRefreshTimer)
      this.outputRefreshTimer = null
      this.outputRefreshInFlight = false
    },
    scrollOutputToBottom() {
      const el = this.$refs.resultPre
      if (!el) return
      el.scrollTop = el.scrollHeight
    },
    ensureAutoRefresh() {
      const hasRunning = Array.isArray(this.logs) && this.logs.some(item => item.status === 1)
      if (hasRunning) {
        this.startAutoRefresh()
      } else {
        this.stopAutoRefresh()
      }
    },
    startAutoRefresh() {
      if (this.autoRefreshTimer) return
      this.autoRefreshTimer = setInterval(() => {
        if (this.autoRefreshInFlight) return
        this.autoRefreshInFlight = true
        this.search(() => {
          this.autoRefreshInFlight = false
        })
      }, 3000)
    },
    stopAutoRefresh() {
      if (!this.autoRefreshTimer) return
      clearInterval(this.autoRefreshTimer)
      this.autoRefreshTimer = null
      this.autoRefreshInFlight = false
    }
  }
}
</script>
<style scoped>
pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  padding: 10px;
  background-color: #4c4c4c;
  color: white;
}
</style>
