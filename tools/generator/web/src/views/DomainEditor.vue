<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type FieldConfig, type DomainConfig, type GenerationResult, type FieldType, type DomainListItem, type FieldDetail } from '../api'

const loading = ref(false)
const result = ref<GenerationResult | null>(null)
const error = ref('')
const showPreview = ref(false)
const fieldTypes = ref<FieldType[]>([])
const domains = ref<DomainListItem[]>([])
const activeTab = ref<'new' | 'existing'>('new')
const loadingDomains = ref(false)
const editingDomain = ref<string | null>(null)

const config = ref<DomainConfig>({
  name: '',
  fields: [{ name: '', type: 'string', enum_values: [] }],
  table_name: '',
  route_base: '',
  soft_delete: false,
  wire: true,
  force: false,
})

onMounted(async () =>{
  try {
    const res = await api.getFieldTypes()
    fieldTypes.value = res.types
    await loadDomains()
  } catch (e) {
    console.error(e)
  }
})

async function loadDomains() {
  loadingDomains.value = true
  try {
    const res = await api.listDomains()
    domains.value = res.domains
  } catch (e) {
    console.error('Failed to load domains:', e)
  } finally {
    loadingDomains.value = false
  }
}

async function loadDomain(domainName: string) {
  loading.value = true
  error.value = ''
  try {
    const detail = await api.getDomainDetail(domainName)
    
    // Map fields from detail to config
    const fields: FieldConfig[] = detail.fields.map((f: FieldDetail) => {
      // Map Go type to field type
      let fieldType = mapGoTypeToFieldType(f.type)
      
      return {
        name: f.snake_name,
        type: fieldType,
        enum_values: f.is_enum ? [] : undefined, // TODO: Extract enum values
      }
    })

    config.value = {
      name: detail.name,
      fields: fields.length > 0 ? fields : [{ name: '', type: 'string', enum_values: [] }],
      table_name: '',
      route_base: '',
      soft_delete: false,
      wire: false,
      force: true, // Auto-enable force when editing
    }

    editingDomain.value = domainName
    activeTab.value = 'new' // Switch to editor tab
  } catch (e: any) {
    error.value = `加载失败: ${e.message}`
  } finally {
    loading.value = false
  }
}

function mapGoTypeToFieldType(goType: string): string {
  // Remove pointer
  goType = goType.replace('*', '')
  
  if (goType === 'string') return 'string'
  if (goType === 'int') return 'int'
  if (goType === 'int64') return 'int64'
  if (goType === 'float64') return 'float64'
  if (goType === 'bool') return 'bool'
  if (goType === 'time.Time') return 'time'
  if (goType.includes('Time')) return 'time?'
  
  // Default to enum for custom types
  return 'enum'
}

function addField() {
  config.value.fields.push({ name: '', type: 'string', enum_values: [] })
}

function removeField(index: number) {
  if (config.value.fields.length > 1) {
    config.value.fields.splice(index, 1)
  }
}

function updateEnumValues(field: FieldConfig, value: string) {
  field.enum_values = value.split('|').map(v => v.trim()).filter(Boolean)
}

async function preview() {
  error.value = ''
  loading.value = true
  try {
    const validFields = config.value.fields.filter(f => f.name.trim())
    result.value = await api.previewDomain({
      ...config.value,
      fields: validFields,
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
    const validFields = config.value.fields.filter(f => f.name.trim())
    result.value = await api.generateDomain({
      ...config.value,
      fields: validFields,
    })
    showPreview.value = true
    // Reload domains list
    await loadDomains()
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function reset() {
  config.value = {
    name: '',
    fields: [{ name: '', type: 'string', enum_values: [] }],
    table_name: '',
    route_base: '',
    soft_delete: false,
    wire: true,
    force: false,
  }
  result.value = null
  showPreview.value = false
  editingDomain.value = null
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
    <h1>📦 生成领域模块 Domain</h1>

    <!-- Tabs -->
    <div class="tabs">
      <button 
        class="tab" 
        :class="{ active: activeTab === 'new' }"
        @click="activeTab = 'new'"
      >
        ✨ 新建模块
      </button>
      <button 
        class="tab" 
        :class="{ active: activeTab === 'existing' }"
        @click="activeTab = 'existing'"
      >
        📋 已生成模块 ({{ domains.length }})
      </button>
    </div>

    <!-- Existing Domains List -->
    <div v-if="activeTab === 'existing'" class="domains-list">
      <div v-if="loadingDomains" class="loading">加载中...</div>
      <div v-else-if="domains.length === 0" class="empty">
        <p>暂无已生成的领域模块</p>
        <p class="hint">点击"新建模块"开始创建</p>
      </div>
      <div v-else class="domain-grid">
        <div 
          v-for="domain in domains" 
          :key="domain.name"
          class="domain-card"
          @click="loadDomain(domain.name)"
        >
          <div class="domain-header">
            <h3>{{ domain.name }}</h3>
            <span class="badge">{{ domain.fields.length }} 字段</span>
          </div>
          <div class="domain-fields">
            <span v-for="(field, idx) in domain.fields.slice(0, 5)" :key="idx" class="field-tag">
              {{ field }}
            </span>
            <span v-if="domain.fields.length > 5" class="more">+{{ domain.fields.length - 5 }}</span>
          </div>
          <div class="domain-action">
            点击编辑 →
          </div>
        </div>
      </div>
    </div>

    <!-- Editor (New/Edit) -->
    <div v-if="activeTab === 'new'">
      <!-- Editing indicator -->
      <div v-if="editingDomain" class="editing-banner">
        ✏️ 正在编辑: <strong>{{ editingDomain }}</strong>
        <button class="btn-small" @click="reset">取消编辑</button>
      </div>

      <div class="layout">
      <!-- Left: Form -->
      <div class="form-panel">
        <!-- Usage Guide -->
        <details class="help-tips">
          <summary>📖 使用指南 Usage Guide</summary>
          <div class="tips-content">
            <p><strong>领域名称：</strong>使用 PascalCase 格式，如 <code>User</code>、<code>Order</code></p>
            <p><strong>字段类型：</strong></p>
            <ul>
              <li><code>string</code> - 字符串 (varchar 255)</li>
              <li><code>text</code> - 长文本</li>
              <li><code>int</code> / <code>int64</code> - 整数</li>
              <li><code>time</code> - 时间戳，<code>time?</code> - 可选时间</li>
              <li><code>enum</code> - 枚举类型，需填写枚举值（用 | 分隔）</li>
            </ul>
            <p><strong>提示：</strong>勾选"自动注入到 main.go"可自动完成模块注册，无需手动修改代码</p>
          </div>
        </details>

        <div class="form-group">
          <label>
            领域名称 Domain Name *
            <span class="tooltip" data-tooltip="实体名称，将生成对应的 Go 结构体">ⓘ</span>
          </label>
          <input v-model="config.name" placeholder="User / Order / Product" />
        </div>

        <div class="fields-section">
          <div class="section-header">
            <h3>字段 Fields</h3>
            <button class="btn-add" @click="addField">+ 添加字段</button>
          </div>

          <div class="field-row" v-for="(field, index) in config.fields" :key="index">
            <input
              v-model="field.name"
              placeholder="username / email / status"
              class="field-name"
            />
            <select v-model="field.type" class="field-type">
              <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                {{ t.type }} - {{ t.description }}
              </option>
            </select>
            <input
              v-if="field.type === 'enum'"
              :value="field.enum_values?.join('|')"
              @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
              placeholder="active|inactive|banned"
              class="field-enum"
              data-tooltip="枚举值用 | 分隔，如：active|inactive"
            />
            <button class="btn-remove" @click="removeField(index)" :disabled="config.fields.length === 1">×</button>
          </div>
        </div>

        <div class="options">
          <div class="form-group inline">
            <label data-tooltip="启用后将添加 DeletedAt 字段，删除时标记而非真删除">
              <input type="checkbox" v-model="config.soft_delete" />
              启用软删除 Soft Delete
            </label>
          </div>
          <div class="form-group inline">
            <label data-tooltip="自动在 main.go 中注册此模块">
              <input type="checkbox" v-model="config.wire" />
              自动注入到 main.go
            </label>
          </div>
          <div class="form-group inline">
            <label data-tooltip="覆盖已存在的文件">
              <input type="checkbox" v-model="config.force" />
              强制覆盖 Force
            </label>
          </div>
        </div>

        <details class="advanced">
          <summary>高级选项 Advanced</summary>
          <div class="form-group">
            <label>自定义表名 Table Name</label>
            <input v-model="config.table_name" placeholder="（自动：名称的复数形式）" />
          </div>
          <div class="form-group">
            <label>自定义路由前缀 Route Base</label>
            <input v-model="config.route_base" placeholder="（自动：名称的复数形式）" />
          </div>
        </details>

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

.domains-list {
  min-height: 400px;
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

.domain-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.domain-card {
  background: var(--bg-card);
  border: 2px solid var(--border);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.domain-card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.domain-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.domain-header h3 {
  margin: 0;
  font-size: 1.2rem;
}

.badge {
  padding: 4px 8px;
  background: rgba(99, 102, 241, 0.2);
  color: var(--primary);
  border-radius: 4px;
  font-size: 0.75rem;
}

.domain-fields {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
  min-height: 28px;
}

.field-tag {
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

.domain-action {
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

.form-group input:not([type="checkbox"]),
.form-group select {
  width: 100%;
  padding: 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text);
  font-size: 1rem;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary);
}

.form-group.inline {
  display: inline-block;
  margin-right: 24px;
}

.form-group.inline label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
  cursor: pointer;
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

.tips-content ul {
  margin-left: 20px;
  margin-top: 8px;
  color: var(--text-muted);
}

.tips-content li {
  margin-bottom: 4px;
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

.fields-section {
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

.field-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.field-name {
  flex: 1;
  padding: 10px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
}

.field-type {
  width: 200px;
  padding: 10px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
}

.field-enum {
  width: 150px;
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

.advanced {
  margin-bottom: 20px;
  padding: 16px;
  background: var(--bg-input);
  border-radius: 8px;
}

.advanced summary {
  cursor: pointer;
  color: var(--text-muted);
  margin-bottom: 16px;
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
  max-height: 300px;
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
</style>
