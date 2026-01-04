<script setup lang="ts">
import { ref, computed } from 'vue'
import { api, type ProjectConfig, type GenerationResult } from '../api'

const step = ref(1)
const loading = ref(false)
const result = ref<GenerationResult | null>(null)
const error = ref('')

const config = ref<ProjectConfig>({
  name: '',
  module_name: '',
  framework_version: '',
  framework_replace: '',
})

const defaultModuleName = computed(() => {
  if (!config.value.name) return ''
  return `github.com/soliton-go/${config.value.name}`
})

const effectiveModuleName = computed(() => {
  return config.value.module_name || defaultModuleName.value
})

async function preview() {
  error.value = ''
  loading.value = true
  try {
    result.value = await api.previewInitProject({
      ...config.value,
      module_name: effectiveModuleName.value,
    })
    step.value = 2
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
    result.value = await api.initProject({
      ...config.value,
      module_name: effectiveModuleName.value,
    })
    step.value = 3
  } catch (e: any) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

function reset() {
  step.value = 1
  result.value = null
  config.value = { name: '', module_name: '', framework_version: '', framework_replace: '' }
}
</script>

<template>
  <div class="wizard">
    <h1>🚀 初始化新项目 Init Project</h1>

    <!-- Step indicator -->
    <div class="steps">
      <div class="step" :class="{ active: step >= 1, done: step > 1 }">1. 配置 Configure</div>
      <div class="step" :class="{ active: step >= 2, done: step > 2 }">2. 预览 Preview</div>
      <div class="step" :class="{ active: step >= 3 }">3. 完成 Done</div>
    </div>

    <!-- Step 1: Configuration -->
    <div class="step-content" v-if="step === 1">
      <!-- Help Tips -->
      <details class="help-tips">
        <summary>💡 配置说明 Configuration Help</summary>
        <div class="tips-content">
          <p><strong>项目名称：</strong>将作为项目目录名，建议使用小写字母和连字符，如 <code>my-project</code></p>
          <p><strong>模块名称：</strong>Go 模块的导入路径，通常格式为 <code>github.com/username/project</code></p>
          <p><strong>框架替换：</strong>仅在本地开发 Soliton 框架时使用，指向框架的本地路径</p>
        </div>
      </details>

      <div class="form-group">
        <label>
          项目名称 Project Name *
          <span class="tooltip" title="项目目录名称，将创建此名称的文件夹">ⓘ</span>
        </label>
        <input
          v-model="config.name"
          placeholder="my-awesome-project"
          @keyup.enter="preview"
        />
        <span class="hint">新项目的目录名称</span>
      </div>

      <div class="form-group">
        <label>
          模块名称 Module Name
          <span class="tooltip" title="Go 模块路径，用于 import 语句">ⓘ</span>
        </label>
        <input
          v-model="config.module_name"
          :placeholder="defaultModuleName || 'github.com/yourname/my-project'"
        />
        <span class="hint">Go 模块路径（默认：github.com/soliton-go/{{ config.name || 'project' }}）</span>
      </div>

      <div class="form-group">
        <label>
          框架替换路径 Framework Replace（可选）
          <span class="tooltip" title="本地开发时使用，指向 soliton-go/framework 的路径">ⓘ</span>
        </label>
        <input
          v-model="config.framework_replace"
          placeholder="../framework 或 /path/to/framework"
        />
        <span class="hint">用于开发的 soliton-go/framework 本地路径</span>
      </div>

      <div class="error" v-if="error">{{ error }}</div>

      <div class="actions">
        <button
          class="btn primary"
          :disabled="!config.name || loading"
          @click="preview"
        >
          {{ loading ? '加载中...' : '预览 Preview →' }}
        </button>
      </div>
    </div>

    <!-- Step 2: Preview -->
    <div class="step-content" v-if="step === 2">
      <div class="preview-info">
        <strong>项目 Project:</strong> {{ config.name }}<br>
        <strong>模块 Module:</strong> {{ effectiveModuleName }}
      </div>

      <div class="file-list">
        <h3>将要创建的文件 Files to Create:</h3>
        <div class="file" v-for="file in result?.files" :key="file.path">
          <span class="file-status" :class="file.status">{{ file.status === 'new' ? 'NEW' : file.status === 'skip' ? 'SKIP' : file.status }}</span>
          <span class="file-path">{{ file.path }}</span>
        </div>
      </div>

      <div class="error" v-if="error">{{ error }}</div>

      <div class="actions">
        <button class="btn" @click="step = 1">← 返回</button>
        <button class="btn primary" :disabled="loading" @click="generate">
          {{ loading ? '创建中...' : '创建项目 Create' }}
        </button>
      </div>
    </div>

    <!-- Step 3: Complete -->
    <div class="step-content" v-if="step === 3">
      <div class="success-icon">✅</div>
      <h2>项目创建成功！Project Created Successfully!</h2>
      <p class="success-message">{{ result?.message }}</p>

      <div class="next-steps">
        <h3>下一步 Next Steps:</h3>
        <pre>cd {{ config.name }}
GOWORK=off go mod tidy
soliton-gen domain User --fields "username,email" --wire
GOWORK=off go run ./cmd/main.go</pre>
      </div>

      <div class="actions">
        <button class="btn primary" @click="reset">创建另一个项目</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wizard {
  max-width: 700px;
  margin: 0 auto;
}

h1 {
  margin-bottom: 24px;
}

.steps {
  display: flex;
  gap: 8px;
  margin-bottom: 32px;
}

.step {
  flex: 1;
  padding: 12px;
  background: var(--bg-card);
  border-radius: 8px;
  text-align: center;
  color: var(--text-muted);
  border: 1px solid var(--border);
}

.step.active {
  border-color: var(--primary);
  color: var(--text);
}

.step.done {
  background: rgba(99, 102, 241, 0.2);
}

.step-content {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid var(--border);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-group input {
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

.actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text);
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:hover {
  background: var(--border);
}

.btn.primary {
  background: var(--primary);
  border-color: var(--primary);
}

.btn.primary:hover {
  background: var(--primary-dark);
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
  margin-top: 16px;
}

.preview-info {
  background: var(--bg-input);
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  line-height: 1.8;
}

.file-list h3 {
  margin-bottom: 12px;
}

.file {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border);
}

.file:last-child {
  border-bottom: none;
}

.file-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  text-transform: uppercase;
}

.file-status.new {
  background: rgba(34, 197, 94, 0.2);
  color: var(--success);
}

.file-status.skip {
  background: rgba(148, 163, 184, 0.2);
  color: var(--text-muted);
}

.file-path {
  font-family: monospace;
  font-size: 0.9rem;
}

.success-icon {
  font-size: 64px;
  text-align: center;
  margin-bottom: 16px;
}

.success-message {
  text-align: center;
  color: var(--text-muted);
  margin-bottom: 24px;
}

.next-steps {
  background: var(--bg-input);
  padding: 16px;
  border-radius: 8px;
}

.next-steps h3 {
  margin-bottom: 12px;
}

.next-steps pre {
  background: var(--bg-dark);
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  font-size: 0.9rem;
}
</style>
