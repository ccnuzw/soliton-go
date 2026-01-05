<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api, type MigrationLogEntry, type MigrationResult, type ProjectLayout } from '../api'
import { showError, showSuccess } from '../toast'

type MigrationRun = MigrationResult & { id: string }

const layout = ref<ProjectLayout | null>(null)
const loading = ref(true)
const running = ref(false)
const currentRun = ref<MigrationRun | null>(null)
const logs = ref<MigrationLogEntry[]>([])
const history = ref<MigrationRun[]>([])

const autoTidy = ref(true)
const timeoutSeconds = ref(300)
const keepHistory = ref(5)
const confirmRun = ref(true)

const filterInfo = ref(true)
const filterError = ref(true)
const filterSystem = ref(true)
const filterTidy = ref(true)
const filterMigrate = ref(true)

const showConfirmModal = ref(false)

const filteredLogs = computed(() => {
  return logs.value.filter((log) => {
    if (log.level === 'error' && !filterError.value) return false
    if (log.level !== 'error' && !filterInfo.value) return false
    if (log.step === 'system' && !filterSystem.value) return false
    if (log.step === 'tidy' && !filterTidy.value) return false
    if (log.step === 'migrate' && !filterMigrate.value) return false
    return true
  })
})

const projectPath = computed(() => layout.value?.module_dir || '')
const commandHint = computed(() => {
  return '默认使用 cmd/migrate/main.go；若不存在则回退 cmd/migrate.go'
})

onMounted(async () => {
  await loadLayout()
  loadHistory()
  loading.value = false
})

async function loadLayout() {
  try {
    layout.value = await api.getLayout()
  } catch (e) {
    console.error(e)
  }
}

function loadHistory() {
  const raw = localStorage.getItem('soliton-gen:migration-history')
  if (!raw) return
  try {
    const data = JSON.parse(raw)
    if (Array.isArray(data)) {
      history.value = data
      if (history.value.length > 0) {
        const first = history.value[0]
        if (first) {
          setCurrentRun(first)
        }
      }
    }
  } catch (e) {
    console.error('Failed to parse history:', e)
  }
}

function saveHistory() {
  localStorage.setItem('soliton-gen:migration-history', JSON.stringify(history.value))
}

function setCurrentRun(run: MigrationRun) {
  currentRun.value = run
  logs.value = run.logs || []
}

function clearLogs() {
  logs.value = []
  currentRun.value = null
}

function clearHistory() {
  history.value = []
  saveHistory()
  clearLogs()
}

function formatTime(iso: string) {
  if (!iso) return '-'
  const date = new Date(iso)
  return date.toLocaleString()
}

function formatDuration(ms: number) {
  if (!ms) return '-'
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const remain = seconds % 60
  if (minutes > 0) {
    return `${minutes}m ${remain}s`
  }
  return `${remain}s`
}

function levelLabel(level: string) {
  if (level === 'error') return 'ERROR'
  return 'INFO'
}

function stepLabel(step: string) {
  if (step === 'tidy') return 'TIDY'
  if (step === 'migrate') return 'MIGRATE'
  return 'SYSTEM'
}

async function runMigration() {
  if (!layout.value?.found || !projectPath.value) {
    showError('未检测到项目，请先初始化或切换项目')
    return
  }
  if (confirmRun.value) {
    showConfirmModal.value = true
    return
  }
  await executeMigration()
}

async function executeMigration() {
  showConfirmModal.value = false
  running.value = true

  try {
    const res = await api.runMigration(projectPath.value, autoTidy.value, timeoutSeconds.value)
    const run: MigrationRun = { ...res, id: `${Date.now()}` }
    currentRun.value = run
    logs.value = res.logs || []
    history.value = [run, ...history.value].slice(0, keepHistory.value)
    saveHistory()

    if (res.success) {
      showSuccess(res.message || '迁移完成')
    } else {
      showError(res.message || '迁移失败')
    }
  } catch (e: any) {
    showError(e.message || '迁移失败')
  } finally {
    running.value = false
  }
}

async function copyLogs() {
  const text = logs.value.map((log) => `[${log.time}] [${stepLabel(log.step)}] [${levelLabel(log.level)}] ${log.message}`).join('\n')
  try {
    await navigator.clipboard.writeText(text)
    showSuccess('日志已复制到剪贴板')
  } catch (e: any) {
    showError(`复制失败: ${e.message}`)
  }
}

function downloadLogs() {
  const text = logs.value.map((log) => `[${log.time}] [${stepLabel(log.step)}] [${levelLabel(log.level)}] ${log.message}`).join('\n')
  const blob = new Blob([text], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `migration-log-${Date.now()}.txt`
  a.click()
  URL.revokeObjectURL(url)
}
</script>

<template>
  <div class="page">
    <header class="header">
      <div>
        <h1>🛠️ 迁移中心 Migration Center</h1>
        <p class="subtitle">完整迁移流程 + 详细日志，适合部署或日常维护</p>
      </div>
      <button class="btn primary" :disabled="running" @click="runMigration">
        {{ running ? '迁移执行中...' : '开始迁移' }}
      </button>
    </header>

    <section class="panel">
      <h2>项目信息</h2>
      <div v-if="loading" class="hint">加载项目中...</div>
      <div v-else-if="!layout?.found" class="hint error">未检测到项目，请先初始化或切换项目</div>
      <div v-else class="project-info">
        <div>
          <strong>模块路径：</strong>{{ layout?.module_path }}
        </div>
        <div>
          <strong>项目路径：</strong>{{ projectPath }}
        </div>
        <div class="hint">{{ commandHint }}</div>
      </div>
    </section>

    <section class="panel grid">
      <div class="card">
        <h3>运行设置</h3>
        <label class="checkbox">
          <input type="checkbox" v-model="autoTidy" />
          迁移前执行 go mod tidy
        </label>
        <label class="checkbox">
          <input type="checkbox" v-model="confirmRun" />
          执行前二次确认
        </label>
        <div class="form-group">
          <label>超时 (秒)</label>
          <input type="number" v-model.number="timeoutSeconds" min="60" max="1800" />
          <span class="hint">建议 300-600 秒，超时时会自动中断</span>
        </div>
        <div class="form-group">
          <label>保留历史记录数</label>
          <input type="number" v-model.number="keepHistory" min="1" max="20" />
          <span class="hint">保留最近几次迁移日志，保存在浏览器</span>
        </div>
      </div>

      <div class="card">
        <h3>日志筛选</h3>
        <label class="checkbox">
          <input type="checkbox" v-model="filterInfo" />
          显示 INFO
        </label>
        <label class="checkbox">
          <input type="checkbox" v-model="filterError" />
          显示 ERROR
        </label>
        <label class="checkbox">
          <input type="checkbox" v-model="filterSystem" />
          显示 SYSTEM
        </label>
        <label class="checkbox">
          <input type="checkbox" v-model="filterTidy" />
          显示 TIDY
        </label>
        <label class="checkbox">
          <input type="checkbox" v-model="filterMigrate" />
          显示 MIGRATE
        </label>
        <div class="actions">
          <button class="btn" @click="copyLogs">复制日志</button>
          <button class="btn" @click="downloadLogs">下载日志</button>
          <button class="btn danger" @click="clearLogs">清空当前</button>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-header">
        <h2>迁移日志</h2>
        <div class="status">
          <span v-if="currentRun" :class="currentRun.success ? 'ok' : 'fail'">
            {{ currentRun.success ? '成功' : '失败' }}
          </span>
          <span v-else class="hint">暂无运行记录</span>
          <span class="meta">耗时: {{ formatDuration(currentRun?.duration_ms || 0) }}</span>
          <span class="meta">ExitCode: {{ currentRun?.exit_code ?? '-' }}</span>
        </div>
      </div>
      <div class="log-panel">
        <div v-if="filteredLogs.length === 0" class="hint">暂无日志</div>
        <div v-for="(log, idx) in filteredLogs" :key="idx" class="log-line" :class="log.level">
          <span class="log-time">{{ formatTime(log.time) }}</span>
          <span class="log-step">{{ stepLabel(log.step) }}</span>
          <span class="log-level">{{ levelLabel(log.level) }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-header">
        <h2>历史记录</h2>
        <button class="btn danger" @click="clearHistory">清空历史</button>
      </div>
      <div v-if="history.length === 0" class="hint">暂无历史记录</div>
      <div v-else class="history-list">
        <button v-for="run in history" :key="run.id" class="history-item" @click="setCurrentRun(run)">
          <div>
            <strong>{{ formatTime(run.started_at) }}</strong>
            <span :class="run.success ? 'ok' : 'fail'">{{ run.success ? '成功' : '失败' }}</span>
          </div>
          <div class="hint">耗时: {{ formatDuration(run.duration_ms) }} · {{ run.command }}</div>
        </button>
      </div>
    </section>

    <div v-if="showConfirmModal" class="modal-overlay" @click="showConfirmModal = false">
      <div class="modal" @click.stop>
        <h3>确认执行迁移</h3>
        <p>项目路径：<code>{{ projectPath }}</code></p>
        <p>执行命令：<code>go run ./cmd/migrate</code></p>
        <p>前置 tidy：<strong>{{ autoTidy ? '是' : '否' }}</strong></p>
        <p>超时设置：<strong>{{ timeoutSeconds }} 秒</strong></p>
        <div class="modal-actions">
          <button class="btn" @click="showConfirmModal = false">取消</button>
          <button class="btn primary" @click="executeMigration">确认执行</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px 40px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
}

.subtitle {
  color: var(--text-muted);
}

.panel {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
}

.panel.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.card {
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 16px;
}

.project-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.hint {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.hint.error {
  color: var(--error);
}

.form-group {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

input[type='number'] {
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 8px 10px;
  border-radius: 8px;
}

.checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 12px;
}

.btn {
  border: none;
  background: var(--bg-input);
  color: var(--text);
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
}

.btn.primary {
  background: var(--primary);
  color: white;
}

.btn.danger {
  background: var(--error);
  color: white;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.status {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status .ok {
  color: var(--success);
  font-weight: 600;
}

.status .fail {
  color: var(--error);
  font-weight: 600;
}

.status .meta {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.log-panel {
  background: var(--bg-dark);
  border-radius: 10px;
  padding: 12px;
  max-height: 420px;
  overflow-y: auto;
}

.log-line {
  display: grid;
  grid-template-columns: 170px 90px 80px 1fr;
  gap: 8px;
  padding: 4px 0;
  border-bottom: 1px dashed rgba(148, 163, 184, 0.2);
  font-size: 0.85rem;
}

.log-line.error {
  color: #fca5a5;
}

.log-time {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
  color: var(--text-muted);
}

.log-step,
.log-level {
  font-weight: 600;
}

.history-list {
  display: grid;
  gap: 10px;
}

.history-item {
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 12px;
  text-align: left;
  cursor: pointer;
  color: var(--text);
}

.history-item .ok {
  margin-left: 8px;
  color: var(--success);
}

.history-item .fail {
  margin-left: 8px;
  color: var(--error);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--bg-card);
  border: 1px solid var(--border);
  padding: 20px;
  border-radius: 12px;
  width: min(480px, 90%);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 16px;
}

code {
  background: var(--bg-dark);
  padding: 2px 6px;
  border-radius: 6px;
}

@media (max-width: 960px) {
  .panel.grid {
    grid-template-columns: 1fr;
  }

  .log-line {
    grid-template-columns: 1fr;
  }
}
</style>
