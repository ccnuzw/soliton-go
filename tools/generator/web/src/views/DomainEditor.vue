<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { api, type DomainConfig, type GenerationResult, type DomainListItem, type FieldConfig, type FieldType, type FieldDetail } from '../api'
import { showSuccess, showError } from '../toast'

const loading = ref(false)
const result = ref<GenerationResult | null>(null)
const error = ref('')
const showPreview = ref(false)
const fieldTypes = ref<FieldType[]>([])
const domains = ref<DomainListItem[]>([])
const activeTab = ref<'new' | 'existing'>('new')
const loadingDomains = ref(false)
const editingDomain = ref<string | null>(null)
const searchQuery = ref('')
const showDeleteConfirm = ref(false)
const deleteConfirmName = ref<string | null>(null)
const bulkFieldsInput = ref('')
const bulkImportError = ref('')

const config = ref<DomainConfig>({
  name: '',
  remark: '',
  fields: [{ name: '', type: 'string', enum_values: [] }],
  table_name: '',
  route_base: '',
  soft_delete: false,
  wire: true,
  force: false,
})

const filteredDomains = computed(() => {
  if (!domains.value || !searchQuery.value.trim()) {
    return domains.value || []
  }
  const query = searchQuery.value.toLowerCase()
  return domains.value.filter(domain =>
    domain.name.toLowerCase().includes(query) ||
    (domain.remark || '').toLowerCase().includes(query) ||
    domain.fields.some(field => field.toLowerCase().includes(query))
  )
})

onMounted(async () => {
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

async function deleteDomain(domainName: string, event: Event) {
  event.stopPropagation() // 防止触发卡片点击

  // 使用自定义确认对话框
  deleteConfirmName.value = domainName
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (!deleteConfirmName.value) return

  const domainName = deleteConfirmName.value
  showDeleteConfirm.value = false

  try {
    await api.deleteDomain(domainName)
    await loadDomains()
    showSuccess(`领域模块 "${domainName}" 删除成功`)
  } catch (e: any) {
    showError(`删除失败: ${e.message}`)
  } finally {
    deleteConfirmName.value = null
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
      let fieldType = mapGoTypeToFieldType(f.type, f.is_enum)

      return {
        name: f.snake_name,
        type: fieldType,
        comment: f.comment || '',
        enum_values: f.enum_values || [],
      }
    })

    config.value = {
      name: detail.name,
      remark: detail.remark || '',
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

function mapGoTypeToFieldType(goType: string, isEnum?: boolean): string {
  // If backend already identified it as enum, return enum
  if (isEnum) return 'enum'

  // Remove pointer
  goType = goType.replace('*', '')

  if (goType === 'string') return 'string'
  if (goType === 'int') return 'int'
  if (goType === 'int64') return 'int64'
  if (goType === 'float64') return 'float64'
  if (goType === 'bool') return 'bool'
  if (goType === 'time.Time') return 'time'
  if (goType.includes('Time')) return 'time?'
  if (goType === 'datatypes.JSON') return 'json'
  if (goType === '[]byte') return 'bytes'

  // Default to enum for custom types
  return 'enum'
}

function addField() {
  config.value.fields.push({ name: '', type: 'string', comment: '', enum_values: [] })
}

function removeField(index: number) {
  if (config.value.fields.length > 1) {
    config.value.fields.splice(index, 1)
  }
}

function moveFieldUp(index: number) {
  if (index > 0) {
    const fields = config.value.fields
    const item = fields.splice(index, 1)[0]!
    fields.splice(index - 1, 0, item)
  }
}

function moveFieldDown(index: number) {
  const fields = config.value.fields
  if (index < fields.length - 1) {
    const item = fields.splice(index, 1)[0]!
    fields.splice(index + 1, 0, item)
  }
}

function updateEnumValues(field: FieldConfig, value: string) {
  field.enum_values = value.split('|').map(v => v.trim()).filter(Boolean)
}

function parseBulkFields(raw: string): FieldConfig[] {
  const entries: FieldConfig[] = []
  const tokens = raw.split(/\r?\n/).flatMap(line => line.split(','))

  for (const token of tokens) {
    let value = token.trim()
    if (!value) continue

    value = value.replace(/^[-*]\s+/, '').replace(/^\d+[\.\)]\s+/, '')
    if (!value || value.startsWith('#') || value.startsWith('//')) continue

    const parts = value.split(':')
    const name = parts[0]?.trim()
    if (!name) continue

    const rawType = parts.length > 1 ? parts[1] : ''
    let type = rawType ? rawType.trim() : 'string'
    let comment = ''
    if (parts.length > 2) {
      comment = parts.slice(2).join(':').trim()
    }
    if (!type) type = 'string'

    let enumValues: string[] = []
    const enumMatch = type.match(/^enum(?:\(([^)]*)\))?$/)
    if (enumMatch) {
      type = 'enum'
      if (enumMatch[1]) {
        enumValues = enumMatch[1].split('|').map(v => v.trim()).filter(Boolean)
      }
    }

    entries.push({
      name,
      type,
      comment,
      enum_values: enumValues,
    })
  }

  return entries
}

function applyBulkImport() {
  bulkImportError.value = ''
  const parsed = parseBulkFields(bulkFieldsInput.value)
  if (parsed.length === 0) {
    bulkImportError.value = '未识别到有效字段，请检查格式'
    return
  }

  const existingNames = new Set(
    config.value.fields.map(field => field.name.trim()).filter(Boolean),
  )
  const deduped = parsed.filter(field => {
    if (existingNames.has(field.name)) return false
    existingNames.add(field.name)
    return true
  })

  if (deduped.length === 0) {
    bulkImportError.value = '字段已存在，未导入新字段'
    return
  }

  const firstField = config.value.fields[0]
  if (config.value.fields.length === 1 && firstField && !firstField.name.trim()) {
    config.value.fields = deduped
  } else {
    config.value.fields.push(...deduped)
  }

  bulkFieldsInput.value = ''
  showSuccess(`已导入 ${deduped.length} 个字段`)
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
  tidying.value = false
  tidyOutput.value = ''
  tidyError.value = ''

  try {
    const validFields = config.value.fields.filter(f => f.name && f.type)
    result.value = await api.generateDomain({
      ...config.value,
      fields: validFields,
    })
    showPreview.value = true
    // Reload domains list
    await loadDomains()

    // 显示成功提示
    if (result.value.success) {
      showSuccess(result.value.message || '生成成功！')

      // 自动运行 go mod tidy 下载依赖
      await runGoModTidy()
    }
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

// go mod tidy 状态
const tidying = ref(false)
const tidyOutput = ref('')
const tidyError = ref('')

async function runGoModTidy() {
  tidying.value = true
  tidyOutput.value = ''
  tidyError.value = ''

  try {
    // 获取当前项目路径
    const layoutRes = await fetch('/api/layout')
    const layoutData = await layoutRes.json()
    const projectPath = layoutData.module_dir || '.'

    console.log('Running go mod tidy for:', projectPath)

    const response = await fetch('/api/projects/tidy', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ project_path: projectPath }),
    })

    const tidyResult = await response.json()
    console.log('Tidy result:', tidyResult)

    if (tidyResult.success) {
      tidyOutput.value = tidyResult.message || '依赖下载成功'
    } else {
      tidyError.value = tidyResult.error || '依赖下载失败'
      if (tidyResult.output) {
        tidyOutput.value = tidyResult.output
      }
    }
  } catch (e: any) {
    console.error('Tidy error:', e)
    tidyError.value = `依赖下载失败: ${e.message}`
  } finally {
    tidying.value = false
  }
}

function reset() {
  config.value = {
    name: '',
    remark: '',
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
      <button class="tab" :class="{ active: activeTab === 'new' }" @click="activeTab = 'new'">
        ✨ 新建模块
      </button>
      <button class="tab" :class="{ active: activeTab === 'existing' }" @click="activeTab = 'existing'">
        📋 已生成模块 ({{ domains?.length || 0 }})
      </button>
    </div>

    <!-- Existing Domains List -->
    <div v-if="activeTab === 'existing'" class="domains-list">
      <!-- Search Box -->
      <div class="search-box">
        <input v-model="searchQuery" type="text" placeholder="🔍 搜索领域模块或字段..." class="search-input" />
        <span v-if="searchQuery" class="search-clear" @click="searchQuery = ''">✕</span>
      </div>

      <div v-if="loadingDomains" class="loading">加载中...</div>
      <div v-else-if="filteredDomains.length === 0 && !searchQuery" class="empty">
        <p>暂无已生成的领域模块</p>
        <p class="hint">点击"新建模块"开始创建</p>
      </div>
      <div v-else-if="filteredDomains.length === 0 && searchQuery" class="empty">
        <p>未找到匹配的模块</p>
        <p class="hint">尝试其他关键词</p>
      </div>
      <div v-else class="domain-grid">
        <div v-for="domain in filteredDomains" :key="domain.name" class="domain-card" @click="loadDomain(domain.name)">
          <div class="domain-header">
            <div class="domain-title">
              <h3>{{ domain.name }}</h3>
              <span v-if="domain.remark" class="domain-remark">{{ domain.remark }}</span>
            </div>
            <div class="header-actions">
              <span class="badge">{{ domain.fields?.length || 0 }} 字段</span>
              <button class="btn-delete" @click="deleteDomain(domain.name, $event)" title="删除模块">
                🗑️
              </button>
            </div>
          </div>
          <div class="domain-fields">
            <span v-for="(field, idx) in domain.fields" :key="idx" class="field-tag">
              {{ field }}
            </span>
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
              <p><strong>领域名称：</strong>使用 PascalCase 格式，如 <code>User</code>、<code>Order</code>、<code>Product</code></p>
              <p><strong>字段类型：</strong></p>
              <ul>
                <li><code>string</code> - 字符串 (varchar 255)</li>
                <li><code>text</code> - 长文本 (无长度限制)</li>
                <li><code>int</code> - 32位整数</li>
                <li><code>int64</code> - 64位整数 (适合金额、大数值)</li>
                <li><code>float64</code> - 浮点数 (适合评分、重量等)</li>
                <li><code>decimal</code> - 精确小数 (适合金额，精度 10,2)</li>
                <li><code>bool</code> - 布尔值 (true/false)</li>
                <li><code>time</code> - 必填时间戳</li>
                <li><code>time?</code> - 可选时间戳 (可为空)</li>
                <li><code>date</code> - 日期 (无时间部分)</li>
                <li><code>date?</code> - 可选日期</li>
                <li><code>uuid</code> - UUID 字符串 (带索引，适合外键)</li>
                <li><code>json</code> - JSON 对象 (需 gorm.io/datatypes)</li>
                <li><code>jsonb</code> - JSONB (PostgreSQL 专用)</li>
                <li><code>bytes</code> - 二进制数据</li>
                <li><code>enum</code> - 枚举类型，需填写枚举值（用 <code>|</code> 分隔，如 <code>active|inactive|banned</code>）</li>
              </ul>
              <p><strong>字段备注：</strong>可选，填写后会作为代码注释生成在字段定义行末</p>
              <p><strong>注意事项：</strong></p>
              <ul>
                <li><code>ID</code>、<code>CreatedAt</code>、<code>UpdatedAt</code> 字段自动生成，无需手动添加</li>
                <li>启用"软删除"会自动添加 <code>DeletedAt</code> 字段</li>
                <li>使用 ↑↓ 按钮可调整字段顺序</li>
                <li>勾选"强制覆盖"会完全替换现有代码，请谨慎使用</li>
                <li>勾选"自动注入"会将模块注入到 main.go</li>
                <li>生成后会自动运行 <code>go mod tidy</code> 下载依赖</li>
                <li>Value Object / Specification / Policy / Event 相关功能请前往「领域增强」</li>
              </ul>
            </div>
          </details>

          <div class="form-group">
            <label>
              领域名称 Domain Name *
              <span class="tooltip" data-tooltip="实体名称，将生成对应的 Go 结构体">ⓘ</span>
            </label>
            <input v-model="config.name" placeholder="User / Order / Product" />
          </div>
          <div class="form-group">
            <label>领域备注 Remark（可选）</label>
            <input v-model="config.remark" placeholder="用于说明该领域用途" />
            <span class="hint">会显示在已生成模块卡片上</span>
          </div>

          <div class="fields-section">
            <div class="section-header">
              <h3>字段 Fields</h3>
              <button class="btn-add" @click="addField">+ 添加字段</button>
            </div>

            <div class="field-row" v-for="(field, index) in config.fields" :key="index">
              <input v-model="field.name" placeholder="username / email / status" class="field-name" />
              <select v-model="field.type" class="field-type">
                <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                  {{ t.type }} - {{ t.description }}
                </option>
              </select>
              <input v-if="field.type === 'enum'" :value="field.enum_values?.join('|')"
                @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
                placeholder="active|inactive|banned" class="field-enum" data-tooltip="枚举值用 | 分隔，如：active|inactive" />
              <input v-model="field.comment" placeholder="字段备注" class="field-comment" />
              <div class="field-actions">
                <button class="btn-move" @click="moveFieldUp(index)" :disabled="index === 0" title="上移">↑</button>
                <button class="btn-move" @click="moveFieldDown(index)" :disabled="index === config.fields.length - 1"
                  title="下移">↓</button>
                <button class="btn-remove" @click="removeField(index)" :disabled="config.fields.length === 1">×</button>
              </div>
            </div>
          </div>

          <details class="bulk-import">
            <summary>批量导入字段 Bulk Import</summary>
            <p class="bulk-hint">
              支持格式：<code>name:type:comment</code>、<code>name:type</code>、<code>name::comment</code>。
              允许多行或逗号分隔；枚举可写为 <code>status:enum(active|inactive)</code>。
            </p>
            <textarea v-model="bulkFieldsInput" class="bulk-textarea"
              placeholder="username:string:用户名&#10;status:enum(active|inactive):状态&#10;created_at:time?"></textarea>
            <div class="bulk-actions">
              <button class="btn" @click="applyBulkImport" :disabled="!bulkFieldsInput.trim()">导入字段</button>
              <button class="btn" @click="bulkFieldsInput = ''" :disabled="!bulkFieldsInput.trim()">清空</button>
            </div>
            <div class="error" v-if="bulkImportError">{{ bulkImportError }}</div>
          </details>

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

          <!-- Force Warning -->
          <div v-if="config.force" class="force-warning">
            <div class="warning-icon">⚠️</div>
            <div class="warning-content">
              <strong>警告：强制覆盖将永久删除所有手动修改的代码！</strong>
              <p>只在首次生成后立即修改字段时使用。一旦开始写业务逻辑，请勿勾选此选项。</p>
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

  <!-- 删除确认对话框 -->
  <div v-if="showDeleteConfirm" class="modal-overlay" @click="showDeleteConfirm = false">
    <div class="modal-dialog" @click.stop>
      <div class="modal-icon">⚠️</div>
      <h3>确认删除</h3>
      <p class="modal-message">
        确定要删除领域模块 <strong>"{{ deleteConfirmName }}"</strong> 吗？
      </p>
      <p class="modal-warning">
        这将删除整个目录及其所有文件，此操作不可恢复！
      </p>
      <div class="modal-actions">
        <button class="btn" @click="showDeleteConfirm = false">取消</button>
        <button class="btn danger" @click="confirmDelete">确认删除</button>
      </div>
    </div>
  </div>
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

.loading,
.empty {
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

.domain-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.domain-header h3 {
  margin: 0;
  font-size: 1.2rem;
  flex: 1;
}

.domain-remark {
  color: var(--text-muted);
  font-size: 0.85rem;
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

.domain-fields {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
  max-height: 120px;
  overflow-y: auto;
  padding: 2px;
}

.domain-fields::-webkit-scrollbar {
  width: 6px;
}

.domain-fields::-webkit-scrollbar-track {
  background: var(--bg-input);
  border-radius: 3px;
}

.domain-fields::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 3px;
}

.domain-fields::-webkit-scrollbar-thumb:hover {
  background: var(--primary);
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

.bulk-import {
  margin-bottom: 20px;
  padding: 16px;
  background: var(--bg-input);
  border-radius: 8px;
}

.bulk-import summary {
  cursor: pointer;
  color: var(--text-muted);
  margin-bottom: 12px;
}

.bulk-hint {
  margin: 8px 0 12px;
  color: var(--text-muted);
  font-size: 0.9rem;
  line-height: 1.6;
}

.bulk-textarea {
  width: 100%;
  min-height: 120px;
  padding: 12px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text);
  font-size: 0.95rem;
  resize: vertical;
}

.bulk-textarea:focus {
  outline: none;
  border-color: var(--primary);
}

.bulk-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 12px;
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
  width: 100%;
}

.field-name {
  flex: 1.5;
  min-width: 160px;
  padding: 10px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
}

.field-type {
  width: 200px;
  flex-shrink: 0;
  padding: 10px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text);
  cursor: pointer;
}

.field-type option {
  background: var(--bg-card);
  color: var(--text);
}

.field-enum {
  min-width: 150px;
  flex: 1.5;
}

.field-comment {
  flex: 1;
  min-width: 80px;
  color: var(--text-muted);
  font-size: 0.9em;
}

.field-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.btn-move {
  width: 28px;
  height: 28px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text-muted);
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-move:hover:not(:disabled) {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

.btn-move:disabled {
  opacity: 0.3;
  cursor: not-allowed;
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

.force-warning {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-top: 16px;
  padding: 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  animation: fadeIn 0.2s ease;
}

.force-warning .warning-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.force-warning .warning-content {
  flex: 1;
}

.force-warning .warning-content strong {
  color: #ef4444;
  display: block;
  margin-bottom: 4px;
}

.force-warning .warning-content p {
  margin: 0;
  font-size: 0.9rem;
  color: var(--text-muted);
  line-height: 1.4;
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
  max-height: 500px;
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

/* 模态对话框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.modal-dialog {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 32px;
  max-width: 420px;
  width: 90%;
  text-align: center;
  animation: slideIn 0.2s ease;
}

@keyframes slideIn {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }

  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.modal-dialog h3 {
  margin: 0 0 16px;
  font-size: 1.25rem;
  color: var(--text);
}

.modal-message {
  margin: 0 0 12px;
  color: var(--text);
  line-height: 1.5;
}

.modal-message strong {
  color: var(--primary);
}

.modal-warning {
  margin: 0 0 24px;
  padding: 12px;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 8px;
  color: #ef4444;
  font-size: 0.9rem;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.btn.danger {
  background: #ef4444;
  color: white;
}

.btn.danger:hover {
  background: #dc2626;
}
</style>
