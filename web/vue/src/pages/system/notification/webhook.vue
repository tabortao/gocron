<template>
  <el-main>
    <notification-tab></notification-tab>
      <el-form ref="form" :model="form" :rules="formRules" label-width="auto" style="width: 800px;">
        <h3>{{ t('system.webhook') }}</h3>
        <el-alert
          :title="t('system.webhookTip')"
          type="info"
          :closable="false"
          style="margin-bottom: 15px;">
        </el-alert>
        <el-form-item :label="t('system.template')" prop="template">
          <el-input
            type="textarea"
            :rows="8"
            :placeholder="webhookPlaceholder"
            v-model.trim="form.template">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit()">{{ t('common.save') }}</el-button>
        </el-form-item>
        <el-button type="primary" @click="createUrl">{{ t('system.addWebhookUrl') }}</el-button> <br><br>
        <h3>{{ t('system.webhookUrls') }}</h3>
        <el-tag
          v-for="item in webhookUrls"
          :key="item.id"
          closable
          @close="deleteUrl(item)">
          {{item.name}} - {{item.url}}
        </el-tag>
      </el-form>
      <el-dialog
        :title="t('system.addWebhookUrl')"
        v-model="dialogVisible"
        width="30%">
        <el-form :model="form">
          <el-form-item :label="t('system.webhookName')" >
            <el-input v-model.trim="name"></el-input>
          </el-form-item>
          <el-form-item label="URL" >
            <el-input v-model.trim="url"></el-input>
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
  name: 'notification-webhook',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale }
  },
  data () {
    return {
      form: {
        template: ''
      },
      formRules: {},
      webhookUrls: [],
      name: '',
      url: '',
      dialogVisible: false
    }
  },
  computed: {
    webhookPlaceholder() {
      return `{"task_id": "{{.TaskId}}", "task_name": "{{.TaskName}}", "status": "{{.Status}}", "result": "{{.Result}}", "remark": "{{.Remark}}"}`
    },
    computedFormRules() {
      return {
        template: [
          {required: true, message: this.t('system.pleaseEnterTemplate'), trigger: 'blur'}
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
  components: {notificationTab},
  created () {
    this.init()
  },
  methods: {
    createUrl () {
      this.dialogVisible = true
    },
    saveUrl () {
      if (this.name === '' || this.url === '') {
        this.$message.error(this.t('system.incompleteParameters'))
        return
      }
      notificationService.createWebhookUrl({
        name: this.name,
        url: this.url
      }, () => {
        this.dialogVisible = false
        this.init()
      })
    },
    deleteUrl (item) {
      notificationService.removeWebhookUrl(item.id, () => {
        this.init()
      })
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      notificationService.updateWebHook(this.form, () => {
        this.$message.success(this.t('message.updateSuccess'))
        this.init()
      })
    },
    init () {
      this.name = ''
      this.url = ''
      notificationService.webhook((data) => {
        this.form.template = data.template || ''
        this.webhookUrls = data.webhook_urls || []
      })
    }
  }
}
</script>

<style scoped>
  .el-tag + .el-tag {
    margin-left: 10px;
  }
</style>
