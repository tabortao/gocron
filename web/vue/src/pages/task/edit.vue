<template>
  <el-main>
    <el-form ref="form" :model="form" :rules="formRules" label-width="auto">
      <el-input v-model="form.id" type="hidden"></el-input>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="t('task.name')" prop="name">
            <el-input v-model.trim="form.name"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('task.tag')">
            <el-input v-model.trim="form.tag" :placeholder="t('task.tagPlaceholder')"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="form.level === 1">
        <el-col>
          <el-alert :title="t('task.mainTaskTip')" type="info" :closable="false"> </el-alert>
          <el-alert :title="t('task.dependencyTip')" type="info" :closable="false"> </el-alert>
          <br />
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="7">
          <el-form-item :label="t('task.type')">
            <el-select v-model.trim="form.level" :disabled="form.id !== ''">
              <el-option
                v-for="item in levelList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="7" v-if="form.level === 1">
          <el-form-item :label="t('task.dependency')">
            <el-select v-model.trim="form.dependency_status">
              <el-option
                v-for="item in dependencyStatusList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="10">
          <el-form-item :label="t('task.childTaskId')" v-if="form.level === 1">
            <el-input
              v-model.trim="form.dependency_task_id"
              :placeholder="t('task.childTaskIdPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="form.level === 1">
        <el-col :span="12">
          <el-form-item :label="t('task.cronExpression')" prop="spec">
            <el-input v-model.trim="form.spec" :placeholder="t('task.cronPlaceholder')">
              <template #append>
                <el-popover placement="bottom" :width="500" trigger="click">
                  <template #reference>
                    <el-button>{{ t('task.cronExample') }}</el-button>
                  </template>
                  <div>
                    <h4>{{ t('task.cronStandard') }}</h4>
                    <ul style="padding-left: 20px; margin: 10px 0">
                      <li>0 * * * * * - {{ t('message.everyMinute') }}</li>
                      <li>*/20 * * * * * - {{ t('message.every20Seconds') }}</li>
                      <li>0 30 21 * * * - {{ t('message.everyDay21_30') }}</li>
                      <li>0 0 23 * * 6 - {{ t('message.everySaturday23') }}</li>
                    </ul>
                    <h4>{{ t('task.cronShortcut') }}</h4>
                    <ul style="padding-left: 20px; margin: 10px 0">
                      <li>@reboot - {{ t('message.reboot') }}</li>
                      <li>@yearly - {{ t('message.yearly') }}</li>
                      <li>@monthly - {{ t('message.monthly') }}</li>
                      <li>@weekly - {{ t('message.weekly') }}</li>
                      <li>@daily - {{ t('message.daily') }}</li>
                      <li>@hourly - {{ t('message.hourly') }}</li>
                      <li>@every 30s - {{ t('message.every30s') }}</li>
                      <li>@every 1m20s - {{ t('message.every1m20s') }}</li>
                    </ul>
                  </div>
                </el-popover>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="8">
          <el-form-item :label="t('task.protocol')">
            <el-select v-model.trim="form.protocol" @change="handleProtocolChange">
              <el-option
                v-for="item in protocolList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8" v-if="form.protocol === 1">
          <el-form-item :label="t('task.httpMethod')">
            <el-select key="http-method" v-model.trim="form.http_method">
              <el-option
                v-for="item in httpMethods"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8" v-else>
          <el-form-item :label="t('task.taskNode')" prop="host_ids">
            <el-select
              key="shell"
              v-model="form.host_ids"
              filterable
              multiple
              :placeholder="t('task.taskNodePlaceholder')"
            >
              <el-option
                v-for="item in hosts"
                :key="item.id"
                :label="item.alias + ' - ' + item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="16">
          <el-form-item :label="t('task.command')" prop="command">
            <el-input
              type="textarea"
              :rows="5"
              :placeholder="commandPlaceholder"
              v-model="form.command"
              @blur="validateCommand"
            >
            </el-input>
            <div
              v-if="commandWarning"
              class="command-warning"
              style="color: #e6a23c; font-size: 12px; margin-top: 4px"
            >
              ⚠️ {{ commandWarning }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col>
          <el-alert :title="t('task.timeoutTip')" type="info" :closable="false"> </el-alert>
          <el-alert :title="t('task.singleInstanceTip')" type="info" :closable="false"> </el-alert>
          <br />
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="t('task.timeout')" prop="timeout">
            <el-input v-model.number.trim="form.timeout"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="t('task.singleInstance')">
            <el-select v-model.trim="form.multi">
              <el-option
                v-for="item in runStatusList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="12">
          <el-form-item :label="t('task.retryTimes')" prop="retry_times">
            <el-input
              v-model.number.trim="form.retry_times"
              :placeholder="t('task.retryTimesPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="t('task.retryInterval')" prop="retry_interval">
            <el-input
              v-model.number.trim="form.retry_interval"
              :placeholder="t('task.retryIntervalPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="8">
          <el-form-item :label="t('task.notification')">
            <el-select v-model.trim="form.notify_status">
              <el-option
                v-for="item in notifyStatusList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8" v-if="form.notify_status !== 0">
          <el-form-item :label="t('task.notifyType')">
            <el-select v-model="form.notify_type" multiple collapse-tags collapse-tags-tooltip>
              <el-option
                v-for="item in notifyTypes"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8" v-if="form.notify_status !== 0 && form.notify_type.includes(0)">
          <el-form-item :label="t('task.notifyReceiver')">
            <el-select
              key="notify-mail"
              v-model="selectedMailNotifyIds"
              filterable
              multiple
              :placeholder="t('task.notifyReceiverPlaceholder')"
            >
              <el-option
                v-for="item in mailUsers"
                :key="item.id"
                :label="item.username"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="8" v-if="form.notify_status !== 0 && form.notify_type.includes(1)">
          <el-form-item :label="t('task.notifyChannel')">
            <el-select
              key="notify-slack"
              v-model="selectedSlackNotifyIds"
              filterable
              multiple
              :placeholder="t('task.notifyReceiverPlaceholder')"
            >
              <el-option
                v-for="item in slackChannels"
                :key="item.id"
                :label="item.name"
                selected="true"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="8" v-if="form.notify_status !== 0 && form.notify_type.includes(2)">
          <el-form-item :label="t('task.notifyWebhookReceiver')">
            <el-select
              key="notify-webhook"
              v-model="selectedWebhookNotifyIds"
              filterable
              multiple
              :placeholder="t('task.notifyReceiverPlaceholder')"
            >
              <el-option :label="t('task.notifyAllWebhook')" :value="-2"></el-option>
              <el-option
                v-for="item in webhookUrls"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="8" v-if="form.notify_status !== 0 && form.notify_type.includes(3)">
          <el-form-item :label="t('task.notifyServerChan3Receiver')">
            <el-select
              key="notify-serverchan3"
              v-model="selectedServerChan3NotifyIds"
              filterable
              multiple
              :placeholder="t('task.notifyReceiverPlaceholder')"
            >
              <el-option :label="t('task.notifyAllServerChan3')" :value="-3"></el-option>
              <el-option
                v-for="item in serverChan3Urls"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="8" v-if="form.notify_status !== 0 && form.notify_type.includes(4)">
          <el-form-item :label="t('task.notifyBarkReceiver')">
            <el-select
              key="notify-bark"
              v-model="selectedBarkNotifyIds"
              filterable
              multiple
              :placeholder="t('task.notifyReceiverPlaceholder')"
            >
              <el-option :label="t('task.notifyAllBark')" :value="-4"></el-option>
              <el-option
                v-for="item in barkUrls"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="form.notify_status === 3">
        <el-col :span="12">
          <el-form-item :label="t('task.notifyKeyword')" prop="notify_keyword">
            <el-input
              v-model.trim="form.notify_keyword"
              :placeholder="t('task.notifyKeywordPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="16">
          <el-form-item :label="t('task.remark')">
            <el-input type="textarea" :rows="3" v-model="form.remark"> </el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item>
        <el-button type="primary" @click="submit">{{ t('common.save') }}</el-button>
        <el-button @click="cancel">{{ t('common.cancel') }}</el-button>
      </el-form-item>
    </el-form>
  </el-main>
</template>

<script>
import { useI18n } from 'vue-i18n'
import taskService from '../../api/task'
import notificationService from '../../api/notification'
import { validateCronSpec, getCronExamples } from '../../utils/cronValidator'

const createDefaultForm = () => ({
  id: '',
  name: '',
  tag: '',
  level: 1,
  dependency_status: 1,
  dependency_task_id: '',
  spec: '',
  protocol: 2,
  http_method: 1,
  command: '',
  host_id: '',
  host_ids: [],
  timeout: 3600,
  multi: 0,
  notify_status: 0,
  notify_type: [],
  notify_receiver_id: '',
  notify_keyword: '',
  retry_times: 0,
  retry_interval: 0,
  remark: ''
})

export default {
  name: 'task-edit',
  setup() {
    const { t, locale } = useI18n()
    return { t, locale }
  },
  data() {
    return {
      form: createDefaultForm(),
      formRules: {},
      httpMethods: [
        {
          value: 1,
          label: 'get'
        },
        {
          value: 2,
          label: 'post'
        }
      ],
      protocolList: [
        {
          value: 1,
          label: 'http'
        },
        {
          value: 2,
          label: 'shell'
        }
      ],
      levelList: [],
      dependencyStatusList: [],
      runStatusList: [],
      notifyStatusList: [],
      notifyTypes: [],
      hosts: [],
      mailUsers: [],
      slackChannels: [],
      webhookUrls: [],
      serverChan3Urls: [],
      barkUrls: [],
      selectedMailNotifyIds: [],
      selectedSlackNotifyIds: [],
      selectedWebhookNotifyIds: [],
      selectedServerChan3NotifyIds: [],
      selectedBarkNotifyIds: [],
      notifyReceiverInitialized: false
    }
  },
  computed: {
    commandPlaceholder() {
      if (this.form.protocol === 1) {
        return this.t('message.pleaseEnterUrl')
      }
      return this.t('message.pleaseEnterShellCommand')
    },
    commandWarning() {
      if (!this.form.command) return ''
      if (this.form.command.includes('&quot;')) {
        return (
          this.t('message.htmlEntityDetected') || 'HTML 实体编码已检测到，将自动转换为正确的引号'
        )
      }
      return ''
    }
  },
  watch: {
    $route() {
      this.initializeForm()
    },
    'form.notify_status'() {
      this.updateNotifyKeywordRule()
      if (this.form.notify_status === 0) {
        this.form.notify_type = []
        this.form.notify_receiver_id = ''
        this.selectedMailNotifyIds = []
        this.selectedSlackNotifyIds = []
        this.selectedWebhookNotifyIds = []
        this.selectedServerChan3NotifyIds = []
        this.selectedBarkNotifyIds = []
        this.notifyReceiverInitialized = false
      }
    },
    'form.level'() {
      this.updateSpecRule()
    }
  },
  created() {
    this.initFormRules()
    this.initSelectOptions()
    this.loadNotificationOptions()
    this.initializeForm()
  },
  methods: {
    initFormRules() {
      this.formRules = {
        name: [{ required: true, message: this.t('message.pleaseEnterTaskName'), trigger: 'blur' }],
        spec: [
          { required: true, message: this.t('message.pleaseEnterCronExpression'), trigger: 'blur' },
          {
            validator: (rule, value, callback) => this.validateCronSpecField(rule, value, callback),
            trigger: 'blur'
          },
          {
            validator: (rule, value, callback) => this.validateCronSpecField(rule, value, callback),
            trigger: 'change'
          }
        ],
        command: [
          { required: true, message: this.t('message.pleaseEnterCommand'), trigger: 'blur' }
        ],
        timeout: [
          {
            type: 'number',
            required: true,
            message: this.t('message.pleaseEnterValidTimeout'),
            trigger: 'blur'
          }
        ],
        retry_times: [
          {
            type: 'number',
            required: true,
            message: this.t('message.pleaseEnterValidRetryTimes'),
            trigger: 'blur'
          }
        ],
        retry_interval: [
          {
            type: 'number',
            required: true,
            message: this.t('message.pleaseEnterValidRetryInterval'),
            trigger: 'blur'
          }
        ],
        notify_keyword: [
          { required: true, message: this.t('message.pleaseEnterNotifyKeyword'), trigger: 'blur' }
        ],
        host_ids: [
          { required: true, message: this.t('message.selectTaskNode'), trigger: 'blur' },
          {
            validator: (rule, value, callback) => this.validateHostIds(rule, value, callback),
            trigger: 'change'
          }
        ]
      }
    },
    initSelectOptions() {
      this.levelList = [
        { value: 1, label: this.t('task.mainTask') },
        { value: 2, label: this.t('task.childTask') }
      ]
      this.dependencyStatusList = [
        { value: 1, label: this.t('task.strongDependency') },
        { value: 2, label: this.t('task.weakDependency') }
      ]
      this.runStatusList = [
        { value: 0, label: this.t('common.yes') },
        { value: 1, label: this.t('common.no') }
      ]
      this.notifyStatusList = [
        { value: 0, label: this.t('task.notifyDisabled') },
        { value: 1, label: this.t('task.notifyOnFailure') },
        { value: 2, label: this.t('task.notifyAlways') },
        { value: 3, label: this.t('task.notifyKeywordMatch') }
      ]
      this.notifyTypes = [
        { value: 0, label: this.t('task.notifyEmail') },
        { value: 1, label: this.t('task.notifySlack') },
        { value: 2, label: this.t('task.notifyWebhook') },
        { value: 3, label: this.t('task.notifyServerChan3') },
        { value: 4, label: this.t('task.notifyBark') }
      ]
    },
    updateNotifyKeywordRule() {
      const keywordRules = this.formRules.notify_keyword
      const needKeyword = this.form.notify_status === 3
      if (!keywordRules || !keywordRules.length) {
        return
      }
      keywordRules[0].required = needKeyword
      if (!needKeyword) {
        this.form.notify_keyword = ''
        if (this.$refs.form) {
          this.$refs.form.clearValidate('notify_keyword')
        }
      }
      // 移除主动验证，只在用户交互时才验证
    },
    updateSpecRule() {
      const specRules = this.formRules.spec
      if (!specRules || !specRules.length) {
        return
      }
      const needSpec = this.form.level === 1
      specRules[0].required = needSpec
      if (!needSpec && this.$refs.form) {
        this.$refs.form.clearValidate('spec')
      }
      // 移除主动验证，只在用户交互时才验证
    },
    validateHostIds(rule, value, callback) {
      if (Number(this.form.protocol) === 2 && (!value || value.length === 0)) {
        callback(new Error(this.t('message.selectTaskNode')))
        return
      }
      callback()
    },
    handleProtocolChange(value, skipValidation = false) {
      const protocolValue = Number(value)
      if (Number.isNaN(protocolValue)) {
        return
      }
      this.form.protocol = protocolValue
      if (protocolValue === 2) {
        if (!skipValidation) {
          this.$nextTick(() => {
            this.$refs.form && this.$refs.form.validateField('host_ids')
          })
        }
        return
      }
      this.form.host_ids = []
      this.form.host_id = ''
      if (this.$refs.form) {
        this.$refs.form.clearValidate('host_ids')
      }
    },
    validateCronSpecField(rule, value, callback) {
      if (this.form.level !== 1) {
        callback()
        return
      }
      const result = validateCronSpec(value)
      if (!result.valid) {
        callback(new Error(result.message))
        return
      }
      callback()
    },
    validateCommand() {
      if (this.form.command && this.form.command.includes('&quot;')) {
        // 自动修复 HTML 实体编码
        this.form.command = this.form.command
          .replace(/&quot;/g, '"')
          .replace(/&apos;/g, "'")
          .replace(/&lt;/g, '<')
          .replace(/&gt;/g, '>')
          .replace(/&amp;/g, '&')
      }
    },
    resetForm() {
      if (this.$refs.form) {
        this.$refs.form.clearValidate()
      }
      const defaults = createDefaultForm()
      Object.assign(this.form, defaults)
      this.selectedMailNotifyIds = []
      this.selectedSlackNotifyIds = []
      this.selectedWebhookNotifyIds = []
      this.selectedServerChan3NotifyIds = []
      this.selectedBarkNotifyIds = []
      this.notifyReceiverInitialized = false
      this.handleProtocolChange(this.form.protocol, true)
      this.updateNotifyKeywordRule()
      this.updateSpecRule()
    },
    parseNotifyTypeToArray(value) {
      const n = Number(value)
      if (!Number.isFinite(n)) {
        return []
      }
      if (n === 0) {
        return [0]
      }
      const isPowerOfTwo = n > 0 && (n & (n - 1)) === 0
      if (n > 0 && n <= 4 && !isPowerOfTwo) {
        return [n]
      }
      const types = []
      if (n & 1) types.push(0)
      if (n & 2) types.push(1)
      if (n & 4) types.push(2)
      if (n & 8) types.push(3)
      if (n & 16) types.push(4)
      return types
    },
    notifyTypeArrayToMask(types) {
      if (!Array.isArray(types) || types.length === 0) {
        return 0
      }
      let mask = 0
      types.forEach(v => {
        const n = Number(v)
        if (Number.isInteger(n) && n >= 0 && n <= 4) {
          mask |= 1 << n
        }
      })
      return mask
    },
    normalizeAllReceiverSelection() {
      if (this.selectedWebhookNotifyIds.includes(-2)) {
        this.selectedWebhookNotifyIds = [-2]
      }
      if (this.selectedServerChan3NotifyIds.includes(-3)) {
        this.selectedServerChan3NotifyIds = [-3]
      }
      if (this.selectedBarkNotifyIds.includes(-4)) {
        this.selectedBarkNotifyIds = [-4]
      }
    },
    buildNotifyReceiverId(types) {
      const tokens = new Set()
      if (types.includes(0)) {
        this.selectedMailNotifyIds.forEach(v => {
          const id = Number(v)
          if (Number.isFinite(id)) tokens.add(`m:${id}`)
        })
      }
      if (types.includes(1)) {
        this.selectedSlackNotifyIds.forEach(v => {
          const id = Number(v)
          if (Number.isFinite(id)) tokens.add(`s:${id}`)
        })
      }
      if (types.includes(2)) {
        if (this.selectedWebhookNotifyIds.includes(-2)) {
          tokens.add('w:*')
        } else {
          this.selectedWebhookNotifyIds.forEach(v => {
            const id = Number(v)
            if (Number.isFinite(id)) tokens.add(`w:${id}`)
          })
        }
      }
      if (types.includes(3)) {
        if (this.selectedServerChan3NotifyIds.includes(-3)) {
          tokens.add('c:*')
        } else {
          this.selectedServerChan3NotifyIds.forEach(v => {
            const id = Number(v)
            if (Number.isFinite(id)) tokens.add(`c:${id}`)
          })
        }
      }
      if (types.includes(4)) {
        if (this.selectedBarkNotifyIds.includes(-4)) {
          tokens.add('b:*')
        } else {
          this.selectedBarkNotifyIds.forEach(v => {
            const id = Number(v)
            if (Number.isFinite(id)) tokens.add(`b:${id}`)
          })
        }
      }
      return Array.from(tokens).join(',')
    },
    tryInitNotifyReceiverSelections() {
      if (this.notifyReceiverInitialized) {
        return
      }
      if (this.form.notify_status <= 0) {
        return
      }
      const raw = String(this.form.notify_receiver_id || '')
        .split(',')
        .map(v => v.trim())
        .filter(Boolean)
      if (raw.length === 0) {
        return
      }

      const mailSet = new Set((this.mailUsers || []).map(v => v.id))
      const slackSet = new Set((this.slackChannels || []).map(v => v.id))
      const webhookSet = new Set((this.webhookUrls || []).map(v => v.id))
      const serverChan3Set = new Set((this.serverChan3Urls || []).map(v => v.id))
      const barkSet = new Set((this.barkUrls || []).map(v => v.id))

      this.selectedMailNotifyIds = []
      this.selectedSlackNotifyIds = []
      this.selectedWebhookNotifyIds = []
      this.selectedServerChan3NotifyIds = []
      this.selectedBarkNotifyIds = []

      let hasTyped = false
      const legacyIds = []
      raw.forEach(token => {
        if (token.includes(':')) {
          const [k, v] = token.split(':')
          const key = String(k || '').trim()
          const val = String(v || '').trim()
          if (!key || !val) return
          hasTyped = true
          if (key === 'm') {
            const id = Number(val)
            if (Number.isFinite(id)) this.selectedMailNotifyIds.push(id)
          } else if (key === 's') {
            const id = Number(val)
            if (Number.isFinite(id)) this.selectedSlackNotifyIds.push(id)
          } else if (key === 'w') {
            if (val === '*' || val === 'all') {
              this.selectedWebhookNotifyIds = [-2]
              return
            }
            const id = Number(val)
            if (Number.isFinite(id)) this.selectedWebhookNotifyIds.push(id)
          } else if (key === 'c') {
            if (val === '*' || val === 'all') {
              this.selectedServerChan3NotifyIds = [-3]
              return
            }
            const id = Number(val)
            if (Number.isFinite(id)) this.selectedServerChan3NotifyIds.push(id)
          } else if (key === 'b') {
            if (val === '*' || val === 'all') {
              this.selectedBarkNotifyIds = [-4]
              return
            }
            const id = Number(val)
            if (Number.isFinite(id)) this.selectedBarkNotifyIds.push(id)
          }
          return
        }
        const id = Number(token)
        if (Number.isFinite(id)) legacyIds.push(id)
      })

      if (!hasTyped && legacyIds.length > 0) {
        const notifyTypes = Array.isArray(this.form.notify_type) ? this.form.notify_type : []
        if (notifyTypes.length === 1) {
          const t = Number(notifyTypes[0])
          if (t === 0) this.selectedMailNotifyIds = legacyIds
          if (t === 1) this.selectedSlackNotifyIds = legacyIds
          if (t === 2) this.selectedWebhookNotifyIds = legacyIds.includes(-2) ? [-2] : legacyIds
          if (t === 3) this.selectedServerChan3NotifyIds = legacyIds.includes(-3) ? [-3] : legacyIds
          if (t === 4) this.selectedBarkNotifyIds = legacyIds.includes(-4) ? [-4] : legacyIds
        } else {
          legacyIds.forEach(id => {
            if (id === -2) {
              this.selectedWebhookNotifyIds = [-2]
              return
            }
            if (id === -3) {
              this.selectedServerChan3NotifyIds = [-3]
              return
            }
            if (id === -4) {
              this.selectedBarkNotifyIds = [-4]
              return
            }
            if (mailSet.has(id)) {
              this.selectedMailNotifyIds.push(id)
              return
            }
            if (slackSet.has(id)) {
              this.selectedSlackNotifyIds.push(id)
              return
            }
            if (webhookSet.has(id)) {
              this.selectedWebhookNotifyIds.push(id)
              return
            }
            if (serverChan3Set.has(id)) {
              this.selectedServerChan3NotifyIds.push(id)
              return
            }
            if (barkSet.has(id)) {
              this.selectedBarkNotifyIds.push(id)
            }
          })
        }
      }

      this.notifyReceiverInitialized = true
    },
    initializeForm() {
      this.resetForm()
      const id = this.$route.params.id
      if (id) {
        taskService.detail(id, (taskData, hosts) => {
          this.hosts = hosts || []
          if (!taskData) {
            this.$message.error(this.t('message.dataNotFound'))
            this.cancel()
            return
          }
          this.populateForm(taskData)
        })
        return
      }
      taskService.detail(null, (...args) => {
        const hosts = args.length > 1 ? args[1] : args[0]
        this.hosts = hosts || []
        this.handleProtocolChange(this.form.protocol, true)
        this.updateSpecRule()
      })
    },
    populateForm(taskData) {
      Object.assign(this.form, {
        id: taskData.id,
        name: taskData.name,
        tag: taskData.tag,
        level: taskData.level,
        dependency_status: taskData.dependency_status || 1,
        dependency_task_id: taskData.dependency_task_id || '',
        spec: taskData.spec,
        protocol: taskData.protocol,
        http_method: taskData.http_method || 1,
        command: taskData.command,
        timeout: taskData.timeout,
        multi: taskData.multi,
        notify_keyword: taskData.notify_keyword,
        notify_status: taskData.notify_status,
        notify_type:
          taskData.notify_status === 0 ? [] : this.parseNotifyTypeToArray(taskData.notify_type),
        notify_receiver_id: taskData.notify_receiver_id,
        retry_times: taskData.retry_times,
        retry_interval: taskData.retry_interval,
        remark: taskData.remark || ''
      })
      const taskHosts = taskData.hosts || []
      this.form.host_ids = Number(this.form.protocol) === 2 ? taskHosts.map(v => v.host_id) : []
      this.handleProtocolChange(this.form.protocol, true)
      this.updateNotifyKeywordRule()
      this.updateSpecRule()

      this.selectedMailNotifyIds = []
      this.selectedSlackNotifyIds = []
      this.selectedWebhookNotifyIds = []
      this.selectedServerChan3NotifyIds = []
      this.selectedBarkNotifyIds = []
      this.notifyReceiverInitialized = false
      this.tryInitNotifyReceiverSelections()
    },
    loadNotificationOptions() {
      notificationService.mail(data => {
        this.mailUsers = data.mail_users || []
        this.tryInitNotifyReceiverSelections()
      })
      notificationService.slack(data => {
        this.slackChannels = data.channels || []
        this.tryInitNotifyReceiverSelections()
      })
      notificationService.webhook(data => {
        this.webhookUrls = data.webhook_urls || []
        this.tryInitNotifyReceiverSelections()
      })
      notificationService.serverchan3(data => {
        this.serverChan3Urls = data.urls || []
        this.tryInitNotifyReceiverSelections()
      })
      notificationService.bark(data => {
        this.barkUrls = data.urls || []
        this.tryInitNotifyReceiverSelections()
      })
    },
    submit() {
      this.$refs.form.validate(valid => {
        if (!valid) {
          return false
        }
        if (this.form.notify_status > 0) {
          const notifyTypes = Array.isArray(this.form.notify_type) ? this.form.notify_type : []
          if (notifyTypes.length === 0) {
            this.$message.error(this.t('message.selectNotifyType'))
            return false
          }
          if (notifyTypes.includes(0) && this.selectedMailNotifyIds.length === 0) {
            this.$message.error(this.t('message.selectMailReceiver'))
            return false
          }
          if (notifyTypes.includes(1) && this.selectedSlackNotifyIds.length === 0) {
            this.$message.error(this.t('message.selectSlackChannel'))
            return false
          }
          if (notifyTypes.includes(2) && this.selectedWebhookNotifyIds.length === 0) {
            this.$message.error(this.t('message.selectWebhookUrl'))
            return false
          }
          if (notifyTypes.includes(3) && this.selectedServerChan3NotifyIds.length === 0) {
            this.$message.error(this.t('message.selectServerChan3Url'))
            return false
          }
          if (notifyTypes.includes(4) && this.selectedBarkNotifyIds.length === 0) {
            this.$message.error(this.t('message.selectBarkUrl'))
            return false
          }
        }

        this.save()
      })
    },
    save() {
      this.normalizeAllReceiverSelection()
      const payload = { ...this.form }
      if (payload.command) {
        payload.command = payload.command
          .replace(/&quot;/g, '"')
          .replace(/&apos;/g, "'")
          .replace(/&lt;/g, '<')
          .replace(/&gt;/g, '>')
          .replace(/&amp;/g, '&')
      }

      if (Number(payload.protocol) === 2) {
        payload.host_id = (payload.host_ids || []).join(',')
      } else {
        payload.host_id = ''
        payload.host_ids = []
      }

      const notifyTypes = Array.isArray(this.form.notify_type) ? this.form.notify_type : []
      if (payload.notify_status > 0) {
        payload.notify_type = this.notifyTypeArrayToMask(notifyTypes)
        payload.notify_receiver_id = this.buildNotifyReceiverId(notifyTypes)
      } else {
        payload.notify_type = 0
        payload.notify_receiver_id = ''
      }

      taskService.update(payload, () => {
        this.$router.push('/task')
      })
    },
    cancel() {
      this.$router.push('/task')
    }
  }
}
</script>

<style scoped>
:deep(.el-form-item__error) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
