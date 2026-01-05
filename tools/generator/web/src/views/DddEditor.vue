<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  api,
  type FieldConfig,
  type FieldType,
  type GenerationResult,
} from '../api'
import { showSuccess, showError } from '../toast'

const loadingDomains = ref(false)
const domains = ref<string[]>([])
const selectedDomain = ref('')
const fieldTypes = ref<FieldType[]>([])

const valueObjectLoading = ref(false)
const valueObjectResult = ref<GenerationResult | null>(null)
const specLoading = ref(false)
const specResult = ref<GenerationResult | null>(null)
const policyLoading = ref(false)
const policyResult = ref<GenerationResult | null>(null)
const eventLoading = ref(false)
const eventResult = ref<GenerationResult | null>(null)

const valueObject = ref({
  name: '',
  fields: [{ name: 'value', type: 'string', comment: '', enum_values: [] as string[] }],
  force: false,
})

const specification = ref({
  name: '',
  target: '',
  force: false,
})

const policy = ref({
  name: '',
  target: '',
  force: false,
})

const eventFlow = ref({
  name: '',
  topic: '',
  fields: [{ name: 'user_id', type: 'uuid', comment: '', enum_values: [] as string[] }],
  generateEvent: true,
  generateHandler: true,
  handlerTopic: '',
  eventForce: false,
  handlerForce: false,
})

const domainHint = computed(() => {
  if (!selectedDomain.value.trim()) {
    return '请选择或输入已有领域名称'
  }
  const selectedLower = selectedDomain.value.trim().toLowerCase()
  const exists = domains.value.some((d) => d.toLowerCase() === selectedLower)
  if (domains.value.length > 0 && !exists) {
    return '未在项目中找到该领域，请先生成领域模块'
  }
  return '可选领域来自已生成模块'
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
    domains.value = res.domains.map((d) => d.name)
    if (!selectedDomain.value && domains.value.length > 0) {
      const first = domains.value[0]
      if (first) {
        selectedDomain.value = first
      }
    }
  } catch (e) {
    console.error('Failed to load domains:', e)
  } finally {
    loadingDomains.value = false
  }
}

function ensureDomain(): string | null {
  const domain = selectedDomain.value.trim()
  if (!domain) {
    showError('请先选择领域')
    return null
  }
  return domain
}

function addField(target: FieldConfig[]) {
  target.push({ name: '', type: 'string', comment: '', enum_values: [] })
}

function removeField(target: FieldConfig[], index: number) {
  if (target.length > 1) {
    target.splice(index, 1)
  }
}

function updateEnumValues(field: FieldConfig, value: string) {
  field.enum_values = value.split('|').map((v) => v.trim()).filter(Boolean)
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    new: '新建',
    overwrite: '覆盖',
    skip: '跳过',
    error: '错误',
  }
  return map[status] || status
}

function mergeResults(results: GenerationResult[]): GenerationResult {
  const merged: GenerationResult = { success: true, files: [] }
  const messages: string[] = []
  const errors: string[] = []

  results.forEach((res) => {
    merged.success = merged.success && res.success
    merged.files.push(...res.files)
    if (res.message) {
      messages.push(res.message)
    }
    if (res.errors && res.errors.length > 0) {
      errors.push(...res.errors)
    }
  })

  if (messages.length > 0) {
    merged.message = messages.join(' / ')
  }
  if (errors.length > 0) {
    merged.errors = errors
  }
  return merged
}

async function previewValueObject() {
  const domain = ensureDomain()
  if (!domain) return
  if (!valueObject.value.name.trim()) {
    showError('请输入 Value Object 名称')
    return
  }

  valueObjectLoading.value = true
  try {
    const fields = valueObject.value.fields.filter((f) => f.name.trim())
    valueObjectResult.value = await api.previewValueObject({
      domain,
      name: valueObject.value.name,
      fields,
      force: valueObject.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    valueObjectLoading.value = false
  }
}

async function generateValueObject() {
  const domain = ensureDomain()
  if (!domain) return
  if (!valueObject.value.name.trim()) {
    showError('请输入 Value Object 名称')
    return
  }

  valueObjectLoading.value = true
  try {
    const fields = valueObject.value.fields.filter((f) => f.name.trim())
    valueObjectResult.value = await api.generateValueObject({
      domain,
      name: valueObject.value.name,
      fields,
      force: valueObject.value.force,
    })
    if (valueObjectResult.value.success) {
      showSuccess(valueObjectResult.value.message || 'Value Object 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    valueObjectLoading.value = false
  }
}

async function previewSpec() {
  const domain = ensureDomain()
  if (!domain) return
  if (!specification.value.name.trim()) {
    showError('请输入 Specification 名称')
    return
  }

  specLoading.value = true
  try {
    specResult.value = await api.previewSpecification({
      domain,
      name: specification.value.name,
      target: specification.value.target.trim(),
      force: specification.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    specLoading.value = false
  }
}

async function generateSpec() {
  const domain = ensureDomain()
  if (!domain) return
  if (!specification.value.name.trim()) {
    showError('请输入 Specification 名称')
    return
  }

  specLoading.value = true
  try {
    specResult.value = await api.generateSpecification({
      domain,
      name: specification.value.name,
      target: specification.value.target.trim(),
      force: specification.value.force,
    })
    if (specResult.value.success) {
      showSuccess(specResult.value.message || 'Specification 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    specLoading.value = false
  }
}

async function previewPolicy() {
  const domain = ensureDomain()
  if (!domain) return
  if (!policy.value.name.trim()) {
    showError('请输入 Policy 名称')
    return
  }

  policyLoading.value = true
  try {
    policyResult.value = await api.previewPolicy({
      domain,
      name: policy.value.name,
      target: policy.value.target.trim(),
      force: policy.value.force,
    })
  } catch (e: any) {
    showError(e.message)
  } finally {
    policyLoading.value = false
  }
}

async function generatePolicy() {
  const domain = ensureDomain()
  if (!domain) return
  if (!policy.value.name.trim()) {
    showError('请输入 Policy 名称')
    return
  }

  policyLoading.value = true
  try {
    policyResult.value = await api.generatePolicy({
      domain,
      name: policy.value.name,
      target: policy.value.target.trim(),
      force: policy.value.force,
    })
    if (policyResult.value.success) {
      showSuccess(policyResult.value.message || 'Policy 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    policyLoading.value = false
  }
}

async function previewEventFlow() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateEvent && !eventFlow.value.generateHandler) {
    showError('请至少选择一个生成项')
    return
  }

  eventLoading.value = true
  try {
    const tasks: GenerationResult[] = []
    if (eventFlow.value.generateEvent) {
      const fields = eventFlow.value.fields.filter((f) => f.name.trim())
      const res = await api.previewEvent({
        domain,
        name: eventFlow.value.name,
        fields,
        topic: eventFlow.value.topic.trim(),
        force: eventFlow.value.eventForce,
      })
      tasks.push(res)
    }

    if (eventFlow.value.generateHandler) {
      const res = await api.previewEventHandler({
        domain,
        event_name: eventFlow.value.name,
        topic: (eventFlow.value.handlerTopic || eventFlow.value.topic).trim(),
        force: eventFlow.value.handlerForce,
      })
      tasks.push(res)
    }

    eventResult.value = mergeResults(tasks)
  } catch (e: any) {
    showError(e.message)
  } finally {
    eventLoading.value = false
  }
}

async function generateEventFlow() {
  const domain = ensureDomain()
  if (!domain) return
  if (!eventFlow.value.name.trim()) {
    showError('请输入 Event 名称')
    return
  }
  if (!eventFlow.value.generateEvent && !eventFlow.value.generateHandler) {
    showError('请至少选择一个生成项')
    return
  }

  eventLoading.value = true
  try {
    const tasks: GenerationResult[] = []
    if (eventFlow.value.generateEvent) {
      const fields = eventFlow.value.fields.filter((f) => f.name.trim())
      const res = await api.generateEvent({
        domain,
        name: eventFlow.value.name,
        fields,
        topic: eventFlow.value.topic.trim(),
        force: eventFlow.value.eventForce,
      })
      tasks.push(res)
    }

    if (eventFlow.value.generateHandler) {
      const res = await api.generateEventHandler({
        domain,
        event_name: eventFlow.value.name,
        topic: (eventFlow.value.handlerTopic || eventFlow.value.topic).trim(),
        force: eventFlow.value.handlerForce,
      })
      tasks.push(res)
    }

    eventResult.value = mergeResults(tasks)
    if (eventResult.value.success) {
      showSuccess(eventResult.value.message || 'Event 生成成功')
    }
  } catch (e: any) {
    showError(e.message)
  } finally {
    eventLoading.value = false
  }
}
</script>

<template>
  <div class="editor">
    <h1>🧩 领域增强 DDD</h1>
    <p class="subtitle">以中文为主，专业术语保留英文：Value Object / Specification / Policy / Event & Handler</p>

    <details class="help-tips">
      <summary>📖 使用说明 Usage Guide</summary>
      <div class="tips-content">
        <p><strong>适用范围：</strong>以下功能依赖已有领域模块，请先在「领域模块」中生成。</p>
        <p><strong>Value Object：</strong>领域值对象，用于表达不可变的业务概念。</p>
        <p><strong>Specification：</strong>领域规格，用于描述是否满足某个业务条件。</p>
        <p><strong>Policy：</strong>领域策略，用于约束或判断业务行为。</p>
        <p><strong>命名建议：</strong>使用 PascalCase，如 <code>EmailAddress</code>、<code>ActiveUserSpec</code></p>
        <p><strong>Event 与 Handler：</strong>可单独生成，也可在同一操作中组合生成。</p>
        <p><strong>Topic：</strong>留空时自动按领域生成，例如 <code>user.activated</code></p>
      </div>
    </details>

    <div class="domain-select">
      <label>
        选择领域 Domain *
        <span class="tooltip" data-tooltip="需要已存在的领域模块">ⓘ</span>
      </label>
      <div class="domain-input">
        <input v-model="selectedDomain" list="domain-list" placeholder="user / order / product" />
        <datalist id="domain-list">
          <option v-for="d in domains" :key="d" :value="d">{{ d }}</option>
        </datalist>
        <button class="btn" @click="loadDomains" :disabled="loadingDomains">
          {{ loadingDomains ? '刷新中...' : '刷新领域' }}
        </button>
      </div>
      <span class="hint">{{ domainHint }}</span>
    </div>

    <div class="grid">
      <!-- Value Object -->
      <section class="card">
        <div class="card-header">
          <h2>值对象 Value Object</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于建模不可变概念（如金额、邮箱、地址）。</p>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="valueObject.name" placeholder="EmailAddress" />
        </div>

        <div class="fields-section">
          <div class="section-header">
            <h3>字段 Fields</h3>
            <button class="btn-add" @click="addField(valueObject.fields)">+ 添加字段</button>
          </div>
          <div class="field-row" v-for="(field, index) in valueObject.fields" :key="index">
            <input v-model="field.name" placeholder="value" class="field-name" />
            <select v-model="field.type" class="field-type">
              <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                {{ t.type }} - {{ t.description }}
              </option>
            </select>
            <input v-if="field.type === 'enum'" :value="field.enum_values?.join('|')"
              @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
              placeholder="basic|premium" class="field-enum" />
            <input v-model="field.comment" placeholder="字段备注" class="field-comment" />
            <button class="btn-remove" @click="removeField(valueObject.fields, index)"
              :disabled="valueObject.fields.length === 1">×</button>
          </div>
        </div>

        <label class="checkbox">
          <input type="checkbox" v-model="valueObject.force" />
          强制覆盖 Force
        </label>

        <div class="actions">
          <button class="btn" @click="previewValueObject" :disabled="valueObjectLoading">
            {{ valueObjectLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateValueObject" :disabled="valueObjectLoading">
            {{ valueObjectLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>

        <div class="result" v-if="valueObjectResult">
          <div class="result-header">
            <span>{{ valueObjectResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in valueObjectResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Specification -->
      <section class="card">
        <div class="card-header">
          <h2>规格 Specification</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于封装业务规则判断（可复用）。</p>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="specification.name" placeholder="ActiveUserSpec" />
        </div>
        <div class="form-group">
          <label>目标类型 Target（可选）</label>
          <input v-model="specification.target" :placeholder="selectedDomain || 'User'" />
          <span class="hint">留空表示 any</span>
        </div>
        <label class="checkbox">
          <input type="checkbox" v-model="specification.force" />
          强制覆盖 Force
        </label>
        <div class="actions">
          <button class="btn" @click="previewSpec" :disabled="specLoading">
            {{ specLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateSpec" :disabled="specLoading">
            {{ specLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>
        <div class="result" v-if="specResult">
          <div class="result-header">
            <span>{{ specResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in specResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Policy -->
      <section class="card">
        <div class="card-header">
          <h2>策略 Policy</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">用于表达业务策略或决策逻辑。</p>
        <div class="form-group">
          <label>名称 Name *</label>
          <input v-model="policy.name" placeholder="PasswordPolicy" />
        </div>
        <div class="form-group">
          <label>目标类型 Target（可选）</label>
          <input v-model="policy.target" :placeholder="selectedDomain || 'User'" />
          <span class="hint">留空表示 any</span>
        </div>
        <label class="checkbox">
          <input type="checkbox" v-model="policy.force" />
          强制覆盖 Force
        </label>
        <div class="actions">
          <button class="btn" @click="previewPolicy" :disabled="policyLoading">
            {{ policyLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generatePolicy" :disabled="policyLoading">
            {{ policyLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>
        <div class="result" v-if="policyResult">
          <div class="result-header">
            <span>{{ policyResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in policyResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Event + Handler -->
      <section class="card full">
        <div class="card-header">
          <h2>事件与处理器 Event & Handler</h2>
          <span class="tag">Domain</span>
        </div>
        <p class="card-desc">支持事件与处理器组合生成，自动注入模块与 EventBus。</p>
        <div class="form-group">
          <label>事件名称 Event *</label>
          <input v-model="eventFlow.name" placeholder="UserActivated" />
        </div>
        <div class="form-group">
          <label>事件 Topic（可选）</label>
          <input v-model="eventFlow.topic" placeholder="user.activated" />
          <span class="hint">留空自动生成 topic</span>
        </div>

        <div class="fields-section" v-if="eventFlow.generateEvent">
          <div class="section-header">
            <h3>事件字段 Event Fields</h3>
            <button class="btn-add" @click="addField(eventFlow.fields)">+ 添加字段</button>
          </div>
          <div class="field-row" v-for="(field, index) in eventFlow.fields" :key="index">
            <input v-model="field.name" placeholder="user_id" class="field-name" />
            <select v-model="field.type" class="field-type">
              <option v-for="t in fieldTypes" :key="t.type" :value="t.type">
                {{ t.type }} - {{ t.description }}
              </option>
            </select>
            <input v-if="field.type === 'enum'" :value="field.enum_values?.join('|')"
              @input="updateEnumValues(field, ($event.target as HTMLInputElement).value)"
              placeholder="created|deleted" class="field-enum" />
            <input v-model="field.comment" placeholder="字段备注" class="field-comment" />
            <button class="btn-remove" @click="removeField(eventFlow.fields, index)"
              :disabled="eventFlow.fields.length === 1">×</button>
          </div>
        </div>

        <div class="options-grid">
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.generateEvent" />
            生成 Event
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.generateHandler" />
            生成 Handler
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.eventForce" />
            Event 强制覆盖
          </label>
          <label class="checkbox">
            <input type="checkbox" v-model="eventFlow.handlerForce" />
            Handler 强制覆盖
          </label>
        </div>

        <div class="form-group" v-if="eventFlow.generateHandler">
          <label>处理器 Topic（可选）</label>
          <input v-model="eventFlow.handlerTopic" placeholder="默认使用事件 Topic" />
          <span class="hint">留空时复用事件 Topic</span>
        </div>

        <div class="actions">
          <button class="btn" @click="previewEventFlow" :disabled="eventLoading">
            {{ eventLoading ? '预览中...' : '预览 Preview' }}
          </button>
          <button class="btn primary" @click="generateEventFlow" :disabled="eventLoading">
            {{ eventLoading ? '生成中...' : '生成 Generate' }}
          </button>
        </div>

        <div class="result" v-if="eventResult">
          <div class="result-header">
            <span>{{ eventResult.success ? '✅ 生成结果' : '❌ 生成失败' }}</span>
          </div>
          <div class="file-list">
            <div class="file" v-for="file in eventResult.files" :key="file.path">
              <span class="file-status" :class="file.status">{{ getStatusText(file.status) }}</span>
              <span class="file-path">{{ file.path.split('/').pop() }}</span>
            </div>
          </div>
          <div class="message" v-if="eventResult.message">{{ eventResult.message }}</div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.editor {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px 40px;
}

.subtitle {
  color: var(--text-muted);
  margin-bottom: 16px;
}

.domain-select {
  margin: 16px 0 24px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  padding: 16px;
  border-radius: 12px;
}

.domain-input {
  display: flex;
  gap: 12px;
  align-items: center;
}

.domain-input input {
  flex: 1;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 18px;
}

.card.full {
  grid-column: 1 / -1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.tag {
  font-size: 0.75rem;
  padding: 4px 8px;
  background: rgba(99, 102, 241, 0.15);
  border-radius: 8px;
  color: var(--text);
}

.card-desc {
  color: var(--text-muted);
  font-size: 0.9rem;
  margin-bottom: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.form-group label {
  font-weight: 500;
}

.form-group input,
.domain-input input,
select {
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 10px 12px;
  border-radius: 8px;
}

.hint {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.help-tips {
  background: rgba(99, 102, 241, 0.1);
  border: 1px solid var(--primary);
  border-radius: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
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

.fields-section {
  margin-bottom: 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.field-row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1fr 1fr auto;
  gap: 8px;
  margin-bottom: 8px;
}

.field-name,
.field-type,
.field-enum,
.field-comment {
  background: var(--bg-input);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 8px;
  border-radius: 8px;
}

.btn-add,
.btn,
.btn-remove,
.btn.primary {
  border: none;
  cursor: pointer;
  border-radius: 8px;
}

.btn {
  background: var(--bg-input);
  color: var(--text);
  padding: 8px 12px;
}

.btn.primary {
  background: var(--primary);
  color: white;
  padding: 8px 12px;
}

.btn-add {
  background: var(--primary);
  color: white;
  padding: 6px 10px;
}

.btn-remove {
  background: var(--error);
  color: white;
  padding: 0 10px;
}

.checkbox {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  color: var(--text);
}

.actions {
  display: flex;
  gap: 10px;
  margin-top: 12px;
}

.result {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 12px;
  margin-top: 16px;
}

.result-header {
  font-weight: 600;
  margin-bottom: 8px;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.file {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.file-status {
  padding: 2px 6px;
  border-radius: 6px;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.file-status.new {
  background: var(--success);
  color: white;
}

.file-status.overwrite {
  background: var(--warning);
  color: white;
}

.file-status.skip {
  background: var(--bg-input);
  color: var(--text);
}

.file-status.error {
  background: var(--error);
  color: white;
}

.options-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
  margin-bottom: 12px;
}

@media (max-width: 960px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .field-row {
    grid-template-columns: 1fr;
  }

  .domain-input {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
