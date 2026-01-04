<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { api, type ServiceConfig, type GenerationResult, type ServiceListItem, type ServiceDetectionResult } from '../api'
import { showSuccess } from '../toast'

const loading = ref(false)
const result = ref<GenerationResult | null>(null)
const error = ref('')
const showPreview = ref(false)
const services = ref<ServiceListItem[]>([])
const activeTab = ref<'new' | 'existing'>('new')
const loadingServices = ref(false)
const editingService = ref<string | null>(null)
const searchQuery = ref('')
const detectionResult = ref<ServiceDetectionResult | null>(null)
const showDetection = ref(false)
const manualMode = ref(false)

const config = ref<ServiceConfig>({
  name: '',
  methods: [''],
  force: false,  // 默认不勾选
})

const filteredServices = computed(() => {
  if (!services.value || !searchQuery.value.trim()) {
    return services.value || []
  }
  const query = searchQuery.value.toLowerCase()
  return services.value.filter(service => 
    service.name.toLowerCase().includes(query) ||
    service.methods.some(method => method.toLowerCase().includes(query))
  )
})

onMounted(async () => {
  await loadServices()
})

async function loadServices() {
  loadingServices.value = true
  try {
    const res = await api.listServices()
    services.value = res.services
  } catch (e) {
    console.error('Failed to load services:', e)
  } finally {
    loadingServices.value = false
  }
}

async function deleteService(serviceName: string, event: Event) {
  event.stopPropagation() // 防止触发卡片点击
  
  if (!confirm(`确定要删除应用服务 "${serviceName}" 吗？\n\n这将删除服务文件，此操作不可恢复！`)) {
    return
  }

  try {
    await api.deleteService(serviceName)
    await loadServices() // 刷新列表
  } catch (e: any) {
    alert(`删除失败: ${e.message}`)
  }
}

async function detectServiceType() {
  if (!config.value.name || config.value.name.trim() === '') {
    showDetection.value = false
    return
  }
  
  try {
    const result = await api.detectServiceType(config.value.name)
    detectionResult.value = result
    showDetection.value = true
  } catch (e) {
    console.error('检测失败:', e)
    showDetection.value = false
  }
}

// Watch service name changes for auto-detection
let detectTimeout: number | null = null
watch(() => config.value.name, () => {
  if (detectTimeout) clearTimeout(detectTimeout)
  detectTimeout = setTimeout(detectServiceType, 500) as unknown as number
})

async function loadService(serviceName: string) {
  loading.value = true
  error.value = ''
  try {
    const detail = await api.getServiceDetail(serviceName)
    
    config.value = {
      name: detail.name,
      methods: detail.methods.map(m => m.name),
      force: true, // Auto-enable force when editing
    }

    editingService.value = serviceName
    activeTab.value = 'new' // Switch to editor tab
  } catch (e: any) {
    error.value = `加载失败: ${e.message}`
  } finally {
    loading.value = false
  }
}

function addMethod() {
  config.value.methods.push('')
}

function removeMethod(index: number) {
  if (config.value.methods.length > 1) {
    config.value.methods.splice(index, 1)
  }
}

async function preview() {
  error.value = ''
  loading.value = true
  try {
    const validMethods = config.value.methods.filter(m => m.trim())
    result.value = await api.previewService({
      ...config.value,
      methods: validMethods,
    })
    showPreview.value = true
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function generate() {
  error.value = ''
  loading.value = true
  try {
    const validMethods = config.value.methods.filter(m => m.trim())
    result.value = await api.generateService({
      ...config.value,
      methods: validMethods,
    })
    showPreview.value = true
    await loadServices()
    
    // 显示成功提示
    if (result.value.success) {
      showSuccess(result.value.message || '生成成功！')
    }
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function reset() {
  config.value = {
    name: '',
    methods: [''],
    force: false,
  }
  result.value = null
  showPreview.value = false
  editingService.value = null
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    'new': '新建',
    'overwrite': '覆盖',
    'skip': '跳过',
    'error': '错误'
  }
  return map[status] || status
}
</script>

<template>
  <div class="editor">
    <h1>⚙️ 生成应用服务 Service</h1>

    <!-- Tabs -->
    <div class="tabs">
      <button 
        class="tab" 
        :class="{ active: activeTab === 'new' }"
        @click="activeTab = 'new'"
      >
        ✨ 新建服务
      </button>
      <button 
        class="tab" 
        :class="{ active: activeTab === 'existing' }"
        @click="activeTab = 'existing'"
      >
        📋 已生成服务 ({{ services?.length || 0 }})
      </button>
    </div>

    <!-- Existing Services List -->
    <div v-if="activeTab === 'existing'" class="services-list">
      <!-- Search Box -->
      <div class="search-box">
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="🔍 搜索应用服务或方法..."
          class="search-input"
        />
        <span v-if="searchQuery" class="search-clear" @click="searchQuery = ''">✕</span>
      </div>

      <div v-if="loadingServices" class="loading">加载中...</div>
      <div v-else-if="filteredServices.length === 0 && !searchQuery" class="empty">
        <p>暂无已生成的应用服务</p>
        <p class="hint">点击"新建服务"开始创建</p>
      </div>
      <div v-else-if="filteredServices.length === 0 && searchQuery" class="empty">
        <p>未找到匹配的服务</p>
        <p class="hint">尝试其他关键词</p>
      </div>
      <div v-else class="service-grid">
        <div 
          v-for="service in filteredServices" 
          :key="service.name"
          class="service-card"
          @click="loadService(service.name)"
        >
          <div class="service-header">
            <h3>{{ service.name }}</h3>
            <div class="header-actions">
              <span class="badge">{{ service.methods?.length || 0 }} 方法</span>
              <button 
                class="btn-delete" 
                @click="deleteService(service.name, $event)"
                title="删除服务"
              >
                🗑️
              </button>
            </div>
          </div>
          <div class="service-methods">
            <span v-for="(method, idx) in service.methods" :key="idx" class="method-tag">
              {{ method }}
            </span>
          </div>
          <div class="service-action">
            点击编辑 →
          </div>
        </div>
      </div>
    </div>

    <!-- Editor (New/Edit) -->
    <div v-if="activeTab === 'new'">
      <!-- Editing indicator -->
      <div v-if="editingService" class="editing-banner">
        ✏️ 正在编辑: <strong>{{ editingService }}</strong>
        <button class="btn-small" @click="reset">取消编辑</button>
      </div>

      <div class="layout">
      <!-- Left: Form -->
      <div class="form-panel">
        <!-- Usage Guide -->
        <details class="help-tips">
          <summary>📖 使用指南 Usage Guide</summary>
          <div class="tips-content">
            <p><strong>服务名称：</strong>使用 PascalCase 格式，如 <code>OrderService</code>、<code>PaymentService</code></p>
            <p><strong>方法定义：</strong>每行一个方法名，如 <code>CreateOrder</code>、<code>ProcessPayment</code></p>
            <p><strong>默认方法：</strong>如果不填写方法，将自动生成 Create、Get、List 三个基础方法</p>
            <p><strong>提示：</strong>Service 用于编排跨领域的业务逻辑，可以调用多个 Repository</p>
          </div>
        </details>

        <div class="form-group">
          <label>
            服务名称 Service Name *
            <span class="tooltip" data-tooltip="应用服务名称，用于跨领域业务逻辑">ⓘ</span>
          </label>
          <input v-model="config.name" placeholder="OrderService / PaymentService" />
          <span class="hint">如果未包含 "Service" 后缀会自动添加</span>
        </div>

        <!-- Detection Result -->
        <div v-if="showDetection && detectionResult" class="detection-result">
          <div class="detection-header">
            <div class="detection-icon">
              {{ detectionResult.domain_exists ? '✅' : 'ℹ️' }}
            </div>
            <div class="detection-content">
              <p class="detection-message">{{ detectionResult.message }}</p>
              <div class="detection-details">
                <span class="detail-item">
                  <strong>类型:</strong> {{ detectionResult.service_type === 'domain_service' ? '领域服务' : '跨领域服务' }}
                </span>
                <span class="detail-item">
                  <strong>目标:</strong> {{ detectionResult.target_dir }}
                </span>
                <span v-if="detectionResult.should_reuse_dto" class="detail-item highlight">
                  ✅ 复用现有 DTO
                </span>
              </div>
            </div>
            <button 
              class="btn-toggle-manual" 
              @click="manualMode = !manualMode"
              type="button"
            >
              {{ manualMode ? '🔄 自动' : '⚙️ 手动' }}
            </button>
          </div>
        </div>

        <div class="methods-section">
          <div class="section-header">
            <h3>方法 Methods</h3>
            <button class="btn-add" @click="addMethod">+ 添加方法</button>
          </div>

          <div class="method-row" v-for="(_method, index) in config.methods" :key="index">
            <input
              v-model="config.methods[index]"
              placeholder="CreateOrder / ProcessPayment / CancelOrder"
              class="method-name"
            />
            <button class="btn-remove" @click="removeMethod(index)" :disabled="config.methods.length === 1">×</button>
          </div>

          <p class="hint">留空将生成默认方法：Create、Get、List</p>
        </div>

        <div class="options">
          <div class="form-group inline">
            <label data-tooltip="覆盖已存在的文件">
              <input type="checkbox" v-model="config.force" />
              强制覆盖 Force
            </label>
            <div v-if="config.force" class="force-warning">
              ⚠️ <strong>警告：</strong>强制覆盖将<strong>永久删除</strong>所有手动修改的代码！<br>
              只在首次生成后立即修改时使用。一旦开始写业务逻辑，请勿勾选此选项。
            </div>
          </div>
        </div>

        <div class="error" v-if="error">{{ error }}</div>

        <div class="actions">
          <button class="btn" @click="preview" :disabled="!config.name || loading">
            {{ loading ? '加载中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generate" :disabled="!config.name || loading">
            {{ loading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>
      </div>

      <!-- Right: Preview -->
      <div class="preview-panel" v-if="showPreview && result">
        <div class="preview-header">
          <h3>{{ result.success ? '✅ 已生成文件' : '❌ 错误' }}</h3>
          <button class="btn-close" @click="showPreview = false">×</button>
        </div>

        <div class="file-list">
          <div class="file" v-for="file in result.files" :key="file.path">
            <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
            <span class="file-path">{{ file.path.split('/').pop() }}</span>
          </div>
        </div>

        <div class="message" v-if="result.message">{{ result.message }}</div>

        <div class="next-steps">
          <h4>下一步 Next Steps:</h4>
          <ol>
            <li>在服务结构体中注入所需的 Repository</li>
            <li>在每个方法中实现业务逻辑</li>
            <li>在 main.go 中注册服务</li>
          </ol>
        </div>

        <button class="btn primary" @click="reset" style="width: 100%; margin-top: 16px;">
          生成另一个
        </button>
      </div>
      </div> <!-- end layout -->
    </div> <!-- end activeTab === 'new' -->
  </div> <!-- end editor -->
</template>

<style scoped>
.editor {
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 20px;
}

h1 {
  margin-bottom: 24px;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 2px solid var(--border);
}

.tab {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: var(--text-muted);
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s;
  margin-bottom: -2px;
}

.tab:hover {
  color: var(--text);
}

.tab.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

.services-list {
  min-height: 400px;
}

.search-box {
  position: relative;
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  padding: 12px 40px 12px 16px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text);
  font-size: 1rem;
  transition: border-color 0.2s;
}

.search-input:focus {
  outline: none;
  border-color: var(--primary);
}

.search-clear {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  cursor: pointer;
  font-size: 18px;
  padding: 4px 8px;
  transition: color 0.2s;
}

.search-clear:hover {
  color: var(--error);
}

.loading, .empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
}

.empty .hint {
  margin-top: 8px;
  font-size: 0.9rem;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.service-card {
  background: var(--bg-card);
  border: 2px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.service-card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.service-header h3 {
  margin: 0;
  font-size: 1.2rem;
  flex: 1;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-delete {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
  opacity: 0.6;
}

.btn-delete:hover {
  opacity: 1;
  background: rgba(239, 68, 68, 0.1);
}

.badge {
  padding: 4px 8px;
  background: rgba(99, 102, 241, 0.2);
  color: var(--primary);
  border-radius: 4px;
  font-size: 0.75rem;
}

.service-methods {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
  max-height: 120px;
  overflow-y: auto;
  padding: 2px;
}

.service-methods::-webkit-scrollbar {
  width: 6px;
}

.service-methods::-webkit-scrollbar-track {
  background: var(--bg-input);
  border-radius: 3px;
}

.service-methods::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 3px;
}

.service-methods::-webkit-scrollbar-thumb:hover {
  background: var(--primary);
}

.method-tag {
  padding: 2px 8px;
  background: var(--bg-input);
  border-radius: 4px;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.more {
  padding: 2px 8px;
  color: var(--primary);
  font-size: 0.8rem;
}

.service-action {
  color: var(--primary);
  font-size: 0.9rem;
  text-align: right;
}

.editing-banner {
  background: rgba(245, 158, 11, 0.2);
  border: 1px solid var(--warning);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detection-result {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, rgba(139, 92, 246, 0.05) 100%);
  border: 2px solid rgba(99, 102, 241, 0.3);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.detection-header {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.detection-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.detection-content {
  flex: 1;
}

.detection-message {
  margin: 0 0 12px 0;
  font-size: 1rem;
  color: var(--text);
  font-weight: 500;
}

.detection-details {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.detail-item {
  padding: 4px 12px;
  background: var(--bg-input);
  border-radius: 6px;
  font-size: 0.9rem;
  color: var(--text-muted);
}

.detail-item strong {
  color: var(--text);
}

.detail-item.highlight {
  background: rgba(34, 197, 94, 0.2);
  color: var(--success, #22c55e);
  font-weight: 500;
}

.btn-toggle-manual {
  padding: 8px 16px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text);
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-toggle-manual:hover {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.btn-small {
  padding: 6px 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
  cursor: pointer;
  font-size: 0.9rem;
}

.layout {
  display: flex;
  gap: 24px;
}

.form-panel {
  flex: 1;
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid var(--border);
}

.preview-panel {
  width: 350px;
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid var(--border);
  height: fit-content;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.btn-close {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 24px;
  cursor: pointer;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-group input:not([type="checkbox"]) {
  width: 100%;
  padding: 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text);
  font-size: 1rem;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary);
}

.form-group.inline {
  display: inline-block;
}

.form-group.inline label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
  cursor: pointer;
}

.hint {
  display: block;
  margin-top: 6px;
  color: var(--text-muted);
  font-size: 0.85rem;
}

.help-tips {
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid var(--primary);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 24px;
}

.help-tips summary {
  cursor: pointer;
  font-weight: 500;
  color: var(--primary);
  user-select: none;
  list-style: none;
}

.help-tips summary::-webkit-details-marker {
  display: none;
}

.help-tips summary::before {
  content: '▶';
  display: inline-block;
  margin-right: 6px;
  font-size: 0.8em;
  transition: transform 0.2s;
}

.help-tips[open] summary::before {
  transform: rotate(90deg);
}

.tips-content {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.tips-content p {
  margin-bottom: 8px;
  color: var(--text-muted);
  line-height: 1.6;
}

.tips-content code {
  background: var(--bg-dark);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
}

.tooltip {
  display: inline-block;
  margin-left: 6px;
  color: var(--primary);
  cursor: help;
  font-size: 0.9em;
}

.tooltip:hover {
  color: var(--primary-dark);
}

.methods-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.btn-add {
  padding: 8px 16px;
  background: var(--primary);
  border: none;
  border-radius: 6px;
  color: white;
  cursor: pointer;
}

.method-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.method-name {
  flex: 1;
  padding: 10px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
}

.btn-remove {
  width: 40px;
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid var(--error);
  border-radius: 6px;
  color: var(--error);
  cursor: pointer;
}

.btn-remove:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.options {
  margin-bottom: 20px;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text);
  font-size: 1rem;
  cursor: pointer;
}

.btn.primary {
  background: var(--primary);
  border-color: var(--primary);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error {
  background: rgba(239, 68, 68, 0.2);
  color: var(--error);
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 16px;
}

.file-list {
  max-height: 200px;
  overflow-y: auto;
}

.file {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border);
}

.file:last-child {
  border-bottom: none;
}

.file-status {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.7rem;
  text-transform: uppercase;
}

.file-status.new {
  background: rgba(34, 197, 94, 0.2);
  color: var(--success);
}

.file-status.overwrite {
  background: rgba(245, 158, 11, 0.2);
  color: var(--warning);
}

.file-status.skip {
  background: rgba(148, 163, 184, 0.2);
  color: var(--text-muted);
}

.file-path {
  font-family: monospace;
  font-size: 0.85rem;
}

.message {
  margin-top: 16px;
  padding: 12px;
  background: rgba(34, 197, 94, 0.2);
  border-radius: 8px;
  color: var(--success);
}

.next-steps {
  margin-top: 16px;
  padding: 12px;
  background: var(--bg-input);
  border-radius: 8px;
}

.next-steps h4 {
  margin-bottom: 8px;
}

.next-steps ol {
  margin-left: 20px;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.next-steps li {
  margin-bottom: 4px;
}
</style>
