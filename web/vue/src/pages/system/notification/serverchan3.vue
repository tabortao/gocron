<template>
  <el-main>
    <notification-tab></notification-tab>
    <el-form ref="form" :model="form" :rules="formRules" label-width="auto" style="width: 800px">
      <h3>Server 酱³</h3>
      <el-alert
        :title="
          isZh
            ? 'POST 请求，Header: application/x-www-form-urlencoded（title + desp）'
            : 'POST, Content-Type: application/x-www-form-urlencoded (title + desp)'
        "
        type="info"
        :closable="false"
        style="margin-bottom: 15px"
      />
      <el-form-item :label="isZh ? '标题模板' : 'Title template'" prop="title_template">
        <el-input
          type="textarea"
          :rows="3"
          :placeholder="titlePlaceholder"
          v-model.trim="form.title_template"
        />
      </el-form-item>
      <el-form-item
        :label="isZh ? '内容模板（desp）' : 'Content template (desp)'"
        prop="desp_template"
      >
        <el-input
          type="textarea"
          :rows="8"
          :placeholder="despPlaceholder"
          v-model.trim="form.desp_template"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submit()">{{ t('common.save') }}</el-button>
      </el-form-item>

      <el-button type="primary" @click="createUrl">{{
        isZh ? '新增API地址' : 'Add API URL'
      }}</el-button>
      <br /><br />
      <h3>{{ isZh ? 'API地址列表' : 'API URL list' }}</h3>
      <el-tag v-for="item in urls" :key="item.id" closable @close="deleteUrl(item)">
        {{ item.name }} - {{ item.url }}
      </el-tag>
    </el-form>

    <el-dialog :title="isZh ? '新增API地址' : 'Add API URL'" v-model="dialogVisible" width="30%">
      <el-form :model="form">
        <el-form-item :label="isZh ? '名称' : 'Name'">
          <el-input v-model.trim="name" />
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model.trim="url" />
          <div class="url-tip">
            {{
              isZh
                ? '提示：URL 需要为 uid + sendkey 形式，例如 https://{uid}.push.ft07.com/send/{sendkey}.send'
                : 'Tip: URL should be uid + sendkey, e.g. https://{uid}.push.ft07.com/send/{sendkey}.send'
            }}
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveUrl">{{ t('common.confirm') }}</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'

export default {
  name: 'notification-serverchan3',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale }
  },
  data() {
    return {
      form: {
        title_template: '',
        desp_template: ''
      },
      formRules: {},
      urls: [],
      name: '',
      url: '',
      dialogVisible: false
    }
  },
  computed: {
    isZh() {
      return this.locale === 'zh-CN'
    },
    titlePlaceholder() {
      return '{{.TaskName}} - {{.StatusZh}}'
    },
    despPlaceholder() {
      return `**任务**：{{.TaskName}}（ID: {{.TaskId}}）

**状态**：{{.StatusZh}}

{{ if .Host }}**节点**：{{.Host}}

{{ end }}**摘要**：{{.ResultSummary}}

{{ if .Remark }}**备注**：{{.Remark}}{{ end }}`
    },
    computedFormRules() {
      return {
        title_template: [
          {
            required: true,
            message: this.isZh ? '请输入标题模板' : 'Please enter title template',
            trigger: 'blur'
          }
        ],
        desp_template: [
          {
            required: true,
            message: this.isZh ? '请输入内容模板' : 'Please enter content template',
            trigger: 'blur'
          }
        ]
      }
    }
  },
  watch: {
    computedFormRules: {
      handler(newVal) {
        this.formRules = newVal
      },
      immediate: true
    }
  },
  components: { notificationTab },
  created() {
    this.init()
  },
  methods: {
    createUrl() {
      this.dialogVisible = true
    },
    saveUrl() {
      if (this.name === '' || this.url === '') {
        this.$message.error(this.t('system.incompleteParameters'))
        return
      }
      notificationService.createServerchan3Url(
        {
          name: this.name,
          url: this.url
        },
        () => {
          this.dialogVisible = false
          this.init()
        }
      )
    },
    deleteUrl(item) {
      notificationService.removeServerchan3Url(item.id, () => {
        this.init()
      })
    },
    submit() {
      this.$refs['form'].validate(valid => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save() {
      notificationService.updateServerchan3(this.form, () => {
        this.$message.success(this.t('message.updateSuccess'))
        this.init()
      })
    },
    init() {
      this.name = ''
      this.url = ''
      notificationService.serverchan3(data => {
        this.form.title_template = data.title_template || ''
        this.form.desp_template = data.desp_template || ''
        this.urls = data.urls || []
      })
    }
  }
}
</script>

<style scoped>
.el-tag + .el-tag {
  margin-left: 10px;
}
.url-tip {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.4;
  color: var(--el-text-color-secondary);
}
</style>
