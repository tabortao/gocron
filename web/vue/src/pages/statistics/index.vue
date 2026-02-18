<template>
  <el-main class="statistics-main">
    <div class="page-header">
      <div class="page-title">{{ t('statistics.title') }}</div>
      <div class="toolbar">
        <el-button type="primary" size="small" @click="refresh">{{
          t('common.refresh')
        }}</el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="isMobile ? 12 : 16" class="stat-cards">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div
              class="stat-icon"
              style="background: linear-gradient(135deg, #409eff 0%, #337ecc 100%)"
            >
              <el-icon :size="24"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalTasks }}</div>
              <div class="stat-label">{{ t('statistics.totalTasks') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div
              class="stat-icon"
              style="background: linear-gradient(135deg, #67c23a 0%, #4e9a2d 100%)"
            >
              <el-icon :size="24"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.todayExecutions }}</div>
              <div class="stat-label">{{ t('statistics.last7DaysExecutions') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div
              class="stat-icon"
              style="background: linear-gradient(135deg, #e6a23c 0%, #c88a2e 100%)"
            >
              <el-icon :size="24"><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.successRate }}%</div>
              <div class="stat-label">{{ t('statistics.successRate') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div
              class="stat-icon"
              style="background: linear-gradient(135deg, #f56c6c 0%, #c45656 100%)"
            >
              <el-icon :size="24"><CircleClose /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.failedCount }}</div>
              <div class="stat-label">{{ t('statistics.failedCount') }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 趋势图表 -->
    <el-card shadow="hover" class="chart-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('statistics.last7DaysTrend') }}</span>
        </div>
      </template>

      <!-- 折线图可视化 -->
      <div class="chart-wrapper">
        <svg class="line-chart" viewBox="0 0 900 240" xmlns="http://www.w3.org/2000/svg">
          <!-- Y轴 -->
          <line x1="70" y1="15" x2="70" y2="180" stroke="#909399" stroke-width="2" />
          <!-- X轴 -->
          <line x1="70" y1="180" x2="870" y2="180" stroke="#909399" stroke-width="2" />

          <!-- Y轴刻度和标签 -->
          <g v-for="i in 6" :key="'y-tick-' + i">
            <line
              :x1="65"
              :y1="180 - (i - 1) * 33"
              :x2="70"
              :y2="180 - (i - 1) * 33"
              stroke="#909399"
              stroke-width="2"
            />
            <text
              :x="58"
              :y="180 - (i - 1) * 33 + 4"
              text-anchor="end"
              font-size="11"
              fill="#606266"
            >
              {{ Math.round(((i - 1) * getMaxValue()) / 5) }}
            </text>
            <!-- 网格线 -->
            <line
              :x1="70"
              :y1="180 - (i - 1) * 33"
              :x2="870"
              :y2="180 - (i - 1) * 33"
              stroke="#e4e7ed"
              stroke-width="1"
              stroke-dasharray="5,5"
            />
          </g>

          <!-- X轴刻度和标签 -->
          <g v-for="(item, index) in stats.chartData" :key="'x-tick-' + index">
            <line
              :x1="getChartPointX(index)"
              :y1="180"
              :x2="getChartPointX(index)"
              :y2="185"
              stroke="#909399"
              stroke-width="2"
            />
            <text
              :x="getChartPointX(index)"
              :y="200"
              text-anchor="middle"
              font-size="11"
              fill="#606266"
            >
              {{ formatDate(item.date) }}
            </text>
          </g>

          <!-- 成功折线 -->
          <polyline
            v-if="stats.chartData.length > 0"
            :points="getChartLinePoints('success')"
            fill="none"
            stroke="#67C23A"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          />

          <!-- 失败折线 -->
          <polyline
            v-if="stats.chartData.length > 0"
            :points="getChartLinePoints('failed')"
            fill="none"
            stroke="#F56C6C"
            stroke-width="3"
            stroke-linecap="round"
            stroke-linejoin="round"
          />

          <!-- 成功数据点 -->
          <g v-for="(item, index) in stats.chartData" :key="'success-point-' + index">
            <circle
              :cx="getChartPointX(index)"
              :cy="getChartPointY(item.success)"
              r="6"
              fill="#67C23A"
              stroke="#fff"
              stroke-width="2"
              class="data-point"
            />
            <title>{{ item.date }}: {{ t('statistics.success') }} {{ item.success }}</title>
          </g>

          <!-- 失败数据点 -->
          <g v-for="(item, index) in stats.chartData" :key="'failed-point-' + index">
            <circle
              :cx="getChartPointX(index)"
              :cy="getChartPointY(item.failed)"
              r="6"
              fill="#F56C6C"
              stroke="#fff"
              stroke-width="2"
              class="data-point"
            />
            <title>{{ item.date }}: {{ t('statistics.failed') }} {{ item.failed }}</title>
          </g>

          <!-- Y轴标签 -->
          <text
            x="20"
            y="97"
            text-anchor="middle"
            font-size="12"
            fill="#606266"
            transform="rotate(-90, 20, 97)"
          >
            {{ t('statistics.executionCount') }}
          </text>

          <!-- X轴标签 -->
          <text x="470" y="225" text-anchor="middle" font-size="12" fill="#606266">
            {{ t('statistics.date') }}
          </text>
        </svg>
      </div>

      <!-- 图例 -->
      <div class="chart-legend">
        <span class="legend-item">
          <span class="legend-color success-color"></span>
          {{ t('statistics.success') }}
        </span>
        <span class="legend-item">
          <span class="legend-color failed-color"></span>
          {{ t('statistics.failed') }}
        </span>
      </div>
    </el-card>

    <!-- 详细数据表格 -->
    <el-card shadow="hover" class="table-card">
      <template #header>
        <span>{{ t('statistics.last7DaysTrend') }} - {{ t('statistics.detailedData') }}</span>
      </template>
      <el-table :data="stats.last7Days" border style="width: 100%" size="small" class="stats-table">
        <el-table-column prop="date" :label="t('common.date')" min-width="120" />
        <el-table-column prop="total" :label="t('statistics.total')" min-width="80" />
        <el-table-column prop="success" :label="t('statistics.success')" min-width="80">
          <template #default="scope">
            <el-tag type="success" size="small">{{ scope.row.success }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="failed" :label="t('statistics.failed')" min-width="80">
          <template #default="scope">
            <el-tag type="danger" size="small">{{ scope.row.failed }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('statistics.successRate')" min-width="150">
          <template #default="scope">
            <el-progress
              :percentage="calculateSuccessRate(scope.row)"
              :color="getProgressColor(calculateSuccessRate(scope.row))"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </el-main>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Document, CircleCheck, CircleClose, TrendCharts } from '@element-plus/icons-vue'
import statisticsApi from '../../api/statistics'

const { t } = useI18n()

const isMobile = ref(false)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

const stats = ref({
  totalTasks: 0,
  todayExecutions: 0,
  successRate: 0,
  failedCount: 0,
  last7Days: [],
  chartData: []
})

const fetchStatistics = () => {
  statisticsApi.getOverview(data => {
    if (data) {
      const last7Days = data.last_7_days || []

      const total7DaysSuccess = last7Days.reduce((sum, item) => sum + item.success, 0)
      const total7DaysFailed = last7Days.reduce((sum, item) => sum + item.failed, 0)
      const total7DaysExecutions = last7Days.reduce((sum, item) => sum + item.total, 0)

      let successRate7Days = 0
      if (total7DaysExecutions > 0) {
        successRate7Days = Math.round((total7DaysSuccess / total7DaysExecutions) * 1000) / 10
      }

      stats.value = {
        totalTasks: data.total_tasks || 0,
        todayExecutions: total7DaysExecutions,
        successRate: successRate7Days,
        failedCount: total7DaysFailed,
        last7Days: last7Days,
        chartData: [...last7Days].reverse()
      }
    }
  })
}

const calculateSuccessRate = row => {
  if (row.total === 0) return 0
  return Math.round((row.success / row.total) * 100)
}

const getProgressColor = percentage => {
  if (percentage >= 90) return '#67C23A'
  if (percentage >= 70) return '#E6A23C'
  return '#F56C6C'
}

const getMaxValue = () => {
  if (stats.value.chartData.length === 0) return 1
  const allValues = stats.value.chartData.flatMap(item => [item.success, item.failed])
  return Math.max(...allValues, 1)
}

const getChartPointX = index => {
  const totalDays = stats.value.chartData.length
  if (totalDays <= 1) return 470
  const chartWidth = 800
  const spacing = chartWidth / (totalDays - 1)
  return 70 + spacing * index
}

const getChartPointY = value => {
  const maxValue = getMaxValue()
  if (maxValue === 0) return 180
  const chartHeight = 165
  const ratio = value / maxValue
  return 180 - ratio * chartHeight
}

const getChartLinePoints = type => {
  return stats.value.chartData
    .map((item, index) => {
      const x = getChartPointX(index)
      const y = getChartPointY(type === 'success' ? item.success : item.failed)
      return `${x},${y}`
    })
    .join(' ')
}

const formatDate = dateStr => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

const refresh = () => {
  fetchStatistics()
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchStatistics()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.statistics-main {
  padding: 16px 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.stat-cards {
  margin-bottom: 16px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-card :deep(.el-card__body) {
  padding: 16px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 4px;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.chart-card {
  margin-bottom: 16px;
  border-radius: 12px;
}

.chart-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.table-card {
  border-radius: 12px;
}

.table-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  font-size: 16px;
}

.chart-wrapper {
  padding: 10px 0;
  margin-bottom: 12px;
  overflow-x: auto;
}

.line-chart {
  width: 100%;
  min-width: 600px;
  height: 240px;
  display: block;
}

.data-point {
  cursor: pointer;
  transition: r 0.2s;
}

.data-point:hover {
  r: 8;
}

.chart-legend {
  display: flex;
  justify-content: center;
  gap: 30px;
  margin-top: 8px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.legend-color {
  width: 14px;
  height: 14px;
  border-radius: 4px;
}

.success-color {
  background: linear-gradient(135deg, #67c23a 0%, #4e9a2d 100%);
}

.failed-color {
  background: linear-gradient(135deg, #f56c6c 0%, #c45656 100%);
}

.stats-table :deep(.el-progress) {
  line-height: 1;
}

@media screen and (max-width: 768px) {
  .statistics-main {
    padding: 12px;
  }

  .stat-cards {
    margin-bottom: 12px;
  }

  .stat-card :deep(.el-card__body) {
    padding: 12px;
  }

  .stat-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
  }

  .stat-icon :deep(.el-icon) {
    font-size: 20px !important;
  }

  .stat-value {
    font-size: 20px;
  }

  .stat-label {
    font-size: 12px;
  }

  .chart-card :deep(.el-card__body),
  .table-card :deep(.el-card__body) {
    padding: 12px;
  }

  .card-header {
    font-size: 15px;
  }

  .chart-wrapper {
    padding: 8px 0;
  }

  .line-chart {
    min-width: 500px;
    height: 200px;
  }

  .chart-legend {
    gap: 20px;
  }

  .legend-item {
    font-size: 12px;
  }

  .stats-table :deep(.el-progress) {
    display: none;
  }

  .stats-table :deep(.el-table__body-wrapper) {
    overflow-x: auto;
  }
}

@media screen and (max-width: 480px) {
  .statistics-main {
    padding: 8px;
  }

  .page-header {
    margin-bottom: 12px;
  }

  .stat-cards {
    margin-bottom: 8px;
  }

  .stat-card :deep(.el-card__body) {
    padding: 10px;
  }

  .stat-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
  }

  .stat-icon :deep(.el-icon) {
    font-size: 18px !important;
  }

  .stat-value {
    font-size: 18px;
  }

  .stat-label {
    font-size: 11px;
  }

  .line-chart {
    min-width: 400px;
    height: 180px;
  }
}
</style>
