<script setup lang="ts">
import { ref } from 'vue'
import { api, type ServiceConfig, type GenerationResult } from '../api'

const loading = ref(false)
const result = ref<GenerationResult | null>(null)
const error = ref('')
const showPreview = ref(false)

const config = ref<ServiceConfig>({
  name: '',
  methods: [''],
  force: false,
})

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
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function reset() {
  config.value = { name: '', methods: [''], force: false }
  result.value = null
  showPreview.value = false
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
            <span class="tooltip" title="应用服务名称，用于跨领域业务逻辑">ⓘ</span>
          </label>
          <input v-model="config.name" placeholder="OrderService / PaymentService" />
          <span class="hint">如果未包含 "Service" 后缀会自动添加</span>
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
            <label title="覆盖已存在的文件">
              <input type="checkbox" v-model="config.force" />
              强制覆盖 Force
            </label>
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
    </div>
  </div>
</template>

<style scoped>
.editor {
  max-width: 1000px;
  margin: 0 auto;
}

h1 {
  margin-bottom: 24px;
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
