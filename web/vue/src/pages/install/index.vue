<template>
  <el-main class="install-page">
    <div class="install-header">
      <div class="language-switcher">
        <LanguageSwitcher />
      </div>
    </div>
    <el-form
      ref="form"
      :model="form"
      :rules="formRules"
      :label-width="isMobile ? '100px' : '150px'"
      class="install-form"
    >
      <h3 class="section-title">{{ t('install.dbConfig') }}</h3>
      <el-form-item :label="t('install.dbType')" prop="db_type">
        <el-select v-model.trim="form.db_type" @change="update_port" class="w-full">
          <el-option
            v-for="item in dbList"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <template v-if="form.db_type !== 'sqlite'">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('install.dbHost')" prop="db_host">
              <el-input v-model="form.db_host"></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('install.dbPort')" prop="db_port">
              <el-input v-model.number="form.db_port"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('install.dbUser')" prop="db_username">
              <el-input v-model="form.db_username"></el-input>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item :label="t('install.dbPassword')" prop="db_password">
              <el-input v-model="form.db_password" type="password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item
            :label="form.db_type === 'sqlite' ? t('install.dbFilePath') : t('install.dbName')"
            prop="db_name"
          >
            <el-input
              v-model="form.db_name"
              :placeholder="
                form.db_type === 'sqlite'
                  ? t('install.dbFilePathPlaceholder')
                  : t('install.dbNamePlaceholder')
              "
            ></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('install.dbTablePrefix')" prop="db_table_prefix">
            <el-input v-model="form.db_table_prefix"></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <h3 class="section-title">{{ t('install.adminConfig') }}</h3>
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('install.adminUsername')" prop="admin_username">
            <el-input v-model="form.admin_username"></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('install.adminEmail')" prop="admin_email">
            <el-input v-model="form.admin_email"></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('install.adminPassword')" prop="admin_password">
            <el-input
              v-model="form.admin_password"
              type="password"
              :placeholder="t('install.passwordPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
        <el-col :xs="24" :sm="12">
          <el-form-item :label="t('install.confirmPassword')" prop="confirm_admin_password">
            <el-input
              v-model="form.confirm_admin_password"
              type="password"
              :placeholder="t('install.passwordPlaceholder')"
            ></el-input>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item class="submit-item">
        <el-button type="primary" @click="submit()" size="large" class="submit-button">{{
          t('install.install')
        }}</el-button>
      </el-form-item>
    </el-form>

    <!-- 语言选择对话框 -->
    <el-dialog
      v-model="showLanguageDialog"
      :title="currentDialogTitle"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      center
      class="language-dialog"
    >
      <div class="language-selection">
        <p class="language-prompt">{{ currentDialogPrompt }}</p>
        <div class="language-options">
          <el-button
            v-for="lang in availableLanguages"
            :key="lang.value"
            :type="selectedLanguage === lang.value ? 'primary' : 'default'"
            size="large"
            class="language-button"
            @click="selectLanguage(lang.value)"
          >
            <span class="language-icon">{{ lang.icon }}</span>
            <span class="language-label">{{ lang.label }}</span>
          </el-button>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="confirmLanguage" :disabled="!selectedLanguage">
          {{ currentConfirmText }}
        </el-button>
      </template>
    </el-dialog>
  </el-main>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import installService from '../../api/install'
import LanguageSwitcher from '../../components/common/LanguageSwitcher.vue'

export default {
  name: 'index',
  components: { LanguageSwitcher },
  setup() {
    const { t, locale } = useI18n()
    const isMobile = ref(false)

    const setLocale = lang => {
      locale.value = lang
    }

    const checkMobile = () => {
      isMobile.value = window.innerWidth <= 768
    }

    onMounted(() => {
      checkMobile()
      window.addEventListener('resize', checkMobile)
    })

    onUnmounted(() => {
      window.removeEventListener('resize', checkMobile)
    })

    return {
      t,
      setLocale,
      isMobile
    }
  },
  data() {
    return {
      showLanguageDialog: false,
      selectedLanguage: '',
      availableLanguages: [
        {
          value: 'zh-CN',
          label: '简体中文',
          icon: '🇨🇳'
        },
        {
          value: 'en-US',
          label: 'English',
          icon: '🇺🇸'
        }
      ],
      form: {
        db_type: 'mysql',
        db_host: '127.0.0.1',
        db_port: 3306,
        db_username: '',
        db_password: '',
        db_name: '',
        db_table_prefix: '',
        admin_username: '',
        admin_password: '',
        confirm_admin_password: '',
        admin_email: ''
      },
      formRules: {},
      dbList: [
        {
          value: 'mysql',
          label: 'MySQL'
        },
        {
          value: 'postgres',
          label: 'PostgreSql'
        },
        {
          value: 'sqlite',
          label: 'SQLite'
        }
      ],
      default_ports: {
        mysql: 3306,
        postgres: 5432,
        sqlite: 0
      }
    }
  },
  computed: {
    currentDialogTitle() {
      return this.selectedLanguage === 'en-US' ? 'Select Language' : '选择语言'
    },
    currentDialogPrompt() {
      return this.selectedLanguage === 'en-US'
        ? 'Please select your preferred language'
        : '请选择您的首选语言'
    },
    currentConfirmText() {
      return this.selectedLanguage === 'en-US' ? 'Confirm' : '确认'
    }
  },
  created() {
    this.checkAndShowLanguageDialog()
    this.initFormRules()
  },
  mounted() {
    console.log('Install page mounted')
    console.log('Saved locale:', localStorage.getItem('locale'))
    console.log('Show dialog:', this.showLanguageDialog)
  },
  methods: {
    checkAndShowLanguageDialog() {
      const savedLocale = localStorage.getItem('locale')
      console.log('Checking language dialog, savedLocale:', savedLocale)

      console.log('Showing language selection dialog')
      this.showLanguageDialog = true
      this.selectedLanguage = savedLocale || 'en-US'
    },
    selectLanguage(lang) {
      this.selectedLanguage = lang
    },
    confirmLanguage() {
      if (this.selectedLanguage) {
        this.setLocale(this.selectedLanguage)
        localStorage.setItem('locale', this.selectedLanguage)
        this.showLanguageDialog = false
      }
    },
    initFormRules() {
      this.formRules = {
        db_type: [{ required: true, message: this.t('install.selectDb'), trigger: 'blur' }],
        db_name: [{ required: true, message: this.t('install.enterDbName'), trigger: 'blur' }],
        admin_username: [
          { required: true, message: this.t('install.enterAdminUsername'), trigger: 'blur' }
        ],
        admin_email: [
          {
            type: 'email',
            required: true,
            message: this.t('install.enterAdminEmail'),
            trigger: 'blur'
          }
        ],
        admin_password: [
          { required: true, message: this.t('install.enterAdminPassword'), trigger: 'blur' },
          { min: 8, message: this.t('install.passwordMinLength'), trigger: 'blur' }
        ],
        confirm_admin_password: [
          { required: true, message: this.t('install.confirmAdminPassword'), trigger: 'blur' },
          { min: 8, message: this.t('install.passwordMinLength'), trigger: 'blur' }
        ]
      }
    },
    update_port(dbType) {
      this.form['db_port'] = this.default_ports[dbType]
      if (dbType === 'sqlite') {
        this.form['db_host'] = ''
        this.form['db_username'] = ''
        this.form['db_password'] = ''
        this.form['db_name'] = './data/gocron.db'
      } else {
        this.form['db_host'] = '127.0.0.1'
        this.form['db_name'] = ''
      }
    },
    submit() {
      if (this.form.db_type !== 'sqlite') {
        if (!this.form.db_host) {
          this.$message.error(this.t('install.enterDbHost'))
          return
        }
        if (!this.form.db_port) {
          this.$message.error(this.t('install.enterDbPort'))
          return
        }
        if (!this.form.db_username) {
          this.$message.error(this.t('install.enterDbUser'))
          return
        }
        if (!this.form.db_password) {
          this.$message.error(this.t('install.enterDbPassword'))
          return
        }
      }

      this.$refs['form'].validate(valid => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save() {
      installService.store(this.form, () => {
        this.$router.push('/')
      })
    }
  }
}
</script>

<style scoped>
.install-page {
  padding: 20px;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
}

.install-header {
  position: relative;
  width: 100%;
  margin-bottom: 20px;
}

.language-switcher {
  position: absolute;
  top: 0;
  right: 20px;
}

.install-form {
  max-width: 700px;
  margin: 0 auto;
  background: #fff;
  padding: 32px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 24px 0 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.section-title:first-of-type {
  margin-top: 0;
}

.w-full {
  width: 100%;
}

.submit-item {
  margin-top: 32px;
  margin-bottom: 0;
}

.submit-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
}

.language-selection {
  padding: 20px 0;
}

.language-prompt {
  text-align: center;
  font-size: 14px;
  color: #606266;
  margin-bottom: 30px;
  line-height: 1.6;
}

.language-options {
  display: flex;
  flex-direction: column;
  gap: 15px;
  align-items: center;
}

.language-button {
  width: 280px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  font-size: 16px;
  transition: all 0.3s;
}

.language-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.language-icon {
  font-size: 24px;
}

.language-label {
  font-weight: 500;
}

@media screen and (max-width: 768px) {
  .install-page {
    padding: 16px;
  }

  .install-form {
    padding: 20px;
    border-radius: 12px;
  }

  .section-title {
    font-size: 16px;
    margin: 20px 0 12px;
  }

  .install-form :deep(.el-form-item__label) {
    font-size: 14px;
  }

  .submit-button {
    height: 44px;
    font-size: 15px;
  }

  .language-switcher {
    right: 16px;
  }
}

@media screen and (max-width: 480px) {
  .install-page {
    padding: 12px;
  }

  .install-form {
    padding: 16px;
    border-radius: 8px;
  }

  .section-title {
    font-size: 15px;
  }

  .install-form :deep(.el-form-item) {
    margin-bottom: 16px;
  }

  .submit-button {
    height: 42px;
    font-size: 14px;
    border-radius: 10px;
  }

  .language-button {
    width: 100%;
    max-width: 280px;
    height: 54px;
    font-size: 15px;
  }

  .language-dialog :deep(.el-dialog) {
    width: 95% !important;
    margin: 10vh auto !important;
  }
}
</style>
