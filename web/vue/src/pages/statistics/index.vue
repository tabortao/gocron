<template>
  <el-main class="statistics-main">
    <div class="page-header">
        <h2>{{ t('statistics.title') }}</h2>
        <el-button type="primary" size="small" @click="refresh">{{ t('common.refresh') }}</el-button>
      </div>
      
      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stat-cards">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #409EFF;">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.totalTasks }}</div>
                <div class="stat-label">{{ t('statistics.totalTasks') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #67C23A;">
                <el-icon :size="24"><CircleCheck /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.todayExecutions }}</div>
                <div class="stat-label">{{ t('statistics.last7DaysExecutions') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #E6A23C;">
                <el-icon :size="24"><TrendCharts /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.successRate }}%</div>
                <div class="stat-label">{{ t('statistics.successRate') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-content">
              <div class="stat-icon" style="background: #F56C6C;">
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
                {{ Math.round((i - 1) * getMaxValue() / 5) }}
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
            <g v-for="(item, index) in stats.last7Days" :key="'x-tick-' + index">
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
              v-if="stats.last7Days.length > 0"
              :points="getChartLinePoints('success')"
              fill="none"
              stroke="#67C23A"
              stroke-width="3"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            
            <!-- 失败折线 -->
            <polyline 
              v-if="stats.last7Days.length > 0"
              :points="getChartLinePoints('failed')"
              fill="none"
              stroke="#F56C6C"
              stroke-width="3"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            
            <!-- 成功数据点 -->
            <g v-for="(item, index) in stats.last7Days" :key="'success-point-' + index">
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
            <g v-for="(item, index) in stats.last7Days" :key="'failed-point-' + index">
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
            <text x="20" y="97" text-anchor="middle" font-size="12" fill="#606266" transform="rotate(-90, 20, 97)">
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
        <el-table :data="stats.last7Days" border style="width: 100%" size="small">
          <el-table-column prop="date" :label="t('common.date')" width="180" />
          <el-table-column prop="total" :label="t('statistics.total')" width="120" />
          <el-table-column prop="success" :label="t('statistics.success')" width="120">
            <template #default="scope">
              <el-tag type="success" size="small">{{ scope.row.success }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="failed" :label="t('statistics.failed')" width="120">
            <template #default="scope">
              <el-tag type="danger" size="small">{{ scope.row.failed }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('statistics.successRate')">
            <template #default="scope">
              <el-progress 
                :percentage="calculateSuccessRate(scope.row)" 
                :color="getProgressColor(calculateSuccessRate(scope.row))"
              />
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Document, CircleCheck, CircleClose, TrendCharts } from '@element-plus/icons-vue'
import statisticsApi from '../../api/statistics'

const { t } = useI18n()

const stats = ref({
  totalTasks: 0,
  todayExecutions: 0,
  successRate: 0,
  failedCount: 0,
  last7Days: []
})

// 获取统计数据
const fetchStatistics = () => {
  statisticsApi.getOverview((data) => {
    if (data) {
      // 转换后端返回的下划线格式为驼峰格式
      const last7Days = data.last_7_days || []
      
      // 计算最近7天的总成功数和总失败数
      const total7DaysSuccess = last7Days.reduce((sum, item) => sum + item.success, 0)
      const total7DaysFailed = last7Days.reduce((sum, item) => sum + item.failed, 0)
      const total7DaysExecutions = last7Days.reduce((sum, item) => sum + item.total, 0)
      
      // 计算最近7天的成功率
      let successRate7Days = 0
      if (total7DaysExecutions > 0) {
        successRate7Days = Math.round((total7DaysSuccess / total7DaysExecutions) * 1000) / 10
      }
      
      stats.value = {
        totalTasks: data.total_tasks || 0,
        todayExecutions: total7DaysExecutions, // 改为显示7天总执行次数
        successRate: successRate7Days, // 改为7天成功率
        failedCount: total7DaysFailed, // 改为7天失败总数
        last7Days: last7Days
      }
    }
  })
}

// 计算成功率
const calculateSuccessRate = (row) => {
  if (row.total === 0) return 0
  return Math.round((row.success / row.total) * 100)
}

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage >= 90) return '#67C23A'
  if (percentage >= 70) return '#E6A23C'
  return '#F56C6C'
}

// 获取最大值（用于折线图高度计算）
const getMaxValue = () => {
  if (stats.value.last7Days.length === 0) return 1
  const allValues = stats.value.last7Days.flatMap(item => [item.success, item.failed])
  return Math.max(...allValues, 1)
}

// 计算折线图点的X坐标（图表区域：70-870）
const getChartPointX = (index) => {
  const totalDays = stats.value.last7Days.length
  if (totalDays <= 1) return 470 // 中心位置
  const chartWidth = 800 // 870 - 70
  const spacing = chartWidth / (totalDays - 1)
  return 70 + spacing * index
}

// 计算折线图点的Y坐标（图表区域：15-180，需要反转）
const getChartPointY = (value) => {
  const maxValue = getMaxValue()
  if (maxValue === 0) return 180
  const chartHeight = 165 // 180 - 15
  const ratio = value / maxValue
  return 180 - (ratio * chartHeight)
}

// 获取折线的点坐标字符串
const getChartLinePoints = (type) => {
  return stats.value.last7Days.map((item, index) => {
    const x = getChartPointX(index)
    const y = getChartPointY(type === 'success' ? item.success : item.failed)
    return `${x},${y}`
  }).join(' ')
}

// 格式化日期显示
const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 刷新数据
const refresh = () => {
  fetchStatistics()
}

onMounted(() => {
  fetchStatistics()
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
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-3px);
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
  border-radius: 8px;
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
  font-weight: bold;
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
}

.chart-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.table-card :deep(.el-card__body) {
  padding: 16px 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
}

.chart-wrapper {
  padding: 10px 0;
  margin-bottom: 12px;
}

.line-chart {
  width: 100%;
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
  border-radius: 2px;
}

.success-color {
  background: #67C23A;
}

.failed-color {
  background: #F56C6C;
}
</style>
