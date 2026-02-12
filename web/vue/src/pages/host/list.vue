<template>
  <el-main>
    <div class="page-header">
      <div class="page-title">{{ t('host.list') }}</div>
      <div class="toolbar">
        <el-button type="success" v-if="isAdmin" @click="showAgentInstall" icon="Download">{{
          t('host.autoRegister')
        }}</el-button>
        <el-button type="primary" v-if="isAdmin" @click="toEdit(null)">{{
          t('common.add')
        }}</el-button>
        <el-button type="info" @click="refresh" icon="Refresh">{{ t('common.refresh') }}</el-button>
      </div>
    </div>

    <el-card class="card-section filter-card" shadow="never">
      <el-form :inline="true" size="small">
        <el-form-item label="ID">
          <el-input v-model.trim="searchParams.id" style="width: 200px" clearable></el-input>
        </el-form-item>
        <el-form-item :label="t('host.name')">
          <el-input v-model.trim="searchParams.name" style="width: 200px" clearable></el-input>
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
        :total="hostTotal"
        v-model:current-page="searchParams.page"
        v-model:page-size="searchParams.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
      >
      </el-pagination>
      <el-table :data="hosts" tooltip-effect="dark" border style="width: 100%">
        <el-table-column prop="id" label="ID"> </el-table-column>
        <el-table-column prop="alias" :label="t('host.alias')"> </el-table-column>
        <el-table-column prop="name" :label="t('host.name')"> </el-table-column>
        <el-table-column prop="port" :label="t('host.port')"> </el-table-column>
        <el-table-column :label="t('task.viewLog')">
          <template #default="scope">
            <el-button type="success" @click="toTasks(scope.row)">{{ t('task.list') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="t('host.remark')"> </el-table-column>
        <el-table-column
          :label="t('common.operation')"
          :width="locale === 'zh-CN' ? 260 : 300"
          v-if="this.isAdmin"
        >
          <template #default="scope">
            <el-button type="primary" size="small" @click="toEdit(scope.row)">{{
              t('common.edit')
            }}</el-button>
            <el-button type="info" size="small" @click="ping(scope.row)">{{
              t('system.testSend')
            }}</el-button>
            <el-button type="danger" size="small" @click="remove(scope.row)">{{
              t('common.delete')
            }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="agentDialogVisible" :title="t('host.agentInstall')" width="750px">
      <div v-if="installCommand">
        <el-alert
          :title="t('host.installTip')"
          type="info"
          :closable="false"
          style="margin-bottom: 20px"
          show-icon
        />

        <el-tabs v-model="activeTab" type="card">
          <el-tab-pane label="Linux / macOS" name="linux">
            <div style="padding: 15px; background: #f5f7fa; border-radius: 4px">
              <div style="margin-bottom: 10px; color: #606266; font-size: 14px">
                <el-icon style="vertical-align: middle"><Monitor /></el-icon>
                {{ t('host.bashCommand') }}
              </div>
              <el-input
                v-model="installCommand"
                type="textarea"
                :rows="3"
                readonly
                style="font-family: monospace; font-size: 13px"
              />
              <div style="margin-top: 10px; text-align: right">
                <el-button type="primary" @click="copyCommand('linux')" icon="DocumentCopy"
                  >Copy</el-button
                >
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="Windows" name="windows">
            <div style="padding: 15px">
              <el-alert type="warning" :closable="false" style="margin-bottom: 15px">
                <template #title>
                  <strong>{{ t('host.windowsManualInstall') }}</strong>
                </template>
                {{ t('host.windowsManualInstallTip') }}
              </el-alert>

              <el-steps direction="vertical" :active="3">
                <el-step
                  :title="t('host.windowsStep1')"
                  :description="t('host.windowsStep1Desc')"
                />
                <el-step
                  :title="t('host.windowsStep2')"
                  :description="t('host.windowsStep2Desc')"
                />
                <el-step
                  :title="t('host.windowsStep3')"
                  :description="t('host.windowsStep3Desc')"
                />
              </el-steps>
            </div>
          </el-tab-pane>
        </el-tabs>

        <el-divider />

        <div style="padding: 10px 0">
          <el-descriptions :column="1" border>
            <el-descriptions-item :label="t('host.tokenExpires')">
              <el-tag type="warning" effect="plain">{{ expiresAt }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('host.tokenUsage')">
              <span style="color: #67c23a">{{ t('host.tokenReusable') }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <div v-else style="text-align: center; padding: 20px">
        <el-icon class="is-loading" :size="30"><Loading /></el-icon>
        <p>{{ t('common.loading') }}</p>
      </div>
    </el-dialog>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'
import { ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import hostService from '../../api/host'
import agentService from '../../api/agent'
import { useUserStore } from '../../stores/user'
import { copyText } from '../../utils/clipboard'

export default {
  name: 'host-list',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale }
  },
  data() {
    const userStore = useUserStore()
    return {
      hosts: [],
      hostTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: '',
        name: '',
        alias: ''
      },
      isAdmin: userStore.isAdmin,
      agentDialogVisible: false,
      installCommand: '',
      expiresAt: '',
      activeTab: 'linux',
      cachedToken: null,
      cachedTokenExpires: null
    }
  },
  components: {
    Loading
  },
  created() {
    this.search()
  },
  watch: {
    $route(to, from) {
      if (
        to.path === '/host' &&
        (from.path === '/host/create' || from.path.startsWith('/host/edit/'))
      ) {
        this.search()
      }
    }
  },
  methods: {
    changePage(page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize(pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search(callback = null) {
      hostService.list(this.searchParams, data => {
        this.hosts = data.data
        this.hostTotal = data.total
        if (callback) {
          callback()
        }
      })
    },
    remove(item) {
      ElMessageBox.confirm(this.t('message.confirmDeleteNode'), this.t('common.tip'), {
        confirmButtonText: this.t('common.confirm'),
        cancelButtonText: this.t('common.cancel'),
        type: 'warning',
        center: true
      })
        .then(() => {
          hostService.remove(item.id, () => this.refresh())
        })
        .catch(() => {})
    },
    ping(item) {
      if (!item.id || item.id <= 0) {
        this.$message.error(this.t('message.dataNotFound'))
        return
      }
      hostService.ping(item.id, () => {
        this.$message.success(this.t('message.connectionSuccess'))
      })
    },
    toEdit(item) {
      let path = ''
      if (item === null) {
        path = '/host/create'
      } else {
        path = `/host/edit/${item.id}`
      }
      this.$router.push(path)
    },
    refresh() {
      this.search(() => {
        this.$message.success(this.t('message.refreshSuccess'))
      })
    },
    toTasks(item) {
      this.$router.push({
        path: '/task',
        query: {
          host_id: item.id
        }
      })
    },
    showAgentInstall() {
      this.agentDialogVisible = true

      // 检查是否有缓存的token且未过期
      const now = new Date()
      if (this.cachedToken && this.cachedTokenExpires && now < this.cachedTokenExpires) {
        // 使用缓存的token
        this.installCommand = this.cachedToken.install_cmd
        this.expiresAt = this.cachedTokenExpires.toLocaleString()
        return
      }

      // 生成新token
      this.installCommand = ''
      this.expiresAt = ''
      agentService.generateToken(data => {
        this.installCommand = data.install_cmd
        const expiresDate = new Date(data.expires_at)
        this.expiresAt = expiresDate.toLocaleString()

        // 缓存token信息
        this.cachedToken = data
        this.cachedTokenExpires = expiresDate
      })
    },
    copyCommand(type) {
      const cmd = type === 'windows' ? this.installCommandWindows : this.installCommand
      copyText(cmd)
        .then(() => {
          this.$message.success(this.t('message.copySuccess'))
        })
        .catch(() => {
          this.$message.error(this.t('message.copyFailed'))
        })
    }
  }
}
</script>
