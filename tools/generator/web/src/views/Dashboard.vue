<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type ProjectLayout } from '../api'

const layout = ref<ProjectLayout | null>(null)
const loading = ref(true)
const showGuide = ref(false)
const projects = ref<any[]>([])
const showProjectSelector = ref(false)
const switching = ref(false)

onMounted(async () => {
  try {
    layout.value = await api.getLayout()
    // Load available projects
    await loadProjects()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

async function loadProjects() {
  try {
    const response = await fetch('/api/projects/list')
    const data = await response.json()
    projects.value = data.projects || []
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
}

async function switchProject(projectPath: string) {
  switching.value = true
  try {
    const response = await fetch('/api/projects/switch', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: projectPath }),
    })

    if (!response.ok) {
      throw new Error('Failed to switch project')
    }

    // Reload the page to reflect the new project
    window.location.reload()
  } catch (e: any) {
    console.error('Failed to switch project:', e)
    alert('切换项目失败: ' + e.message)
  } finally {
    switching.value = false
  }
}

// go mod tidy 状态
const tidying = ref(false)
const tidyResult = ref<{ success: boolean; message: string } | null>(null)

async function runGoModTidy() {
  if (!layout.value?.module_dir) {
    alert('未检测到项目')
    return
  }

  tidying.value = true
  tidyResult.value = null

  try {
    const response = await fetch('/api/projects/tidy', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ project_path: layout.value.module_dir }),
    })

    const data = await response.json()
    tidyResult.value = {
      success: data.success,
      message: data.success ? (data.message || '依赖更新成功') : (data.error || '依赖更新失败'),
    }

    // 3秒后清除提示
    setTimeout(() => {
      tidyResult.value = null
    }, 3000)
  } catch (e: any) {
    tidyResult.value = {
      success: false,
      message: `依赖更新失败: ${e.message}`,
    }
  } finally {
    tidying.value = false
  }
}
</script>

<template>
  <div class="dashboard">
    <header class="header">
      <h1>欢迎使用 Soliton-Gen</h1>
      <p class="subtitle">轻松生成 DDD 风格的 Go 代码</p>
    </header>

    <!-- Usage Guide -->
    <details class="guide" :open="showGuide">
      <summary>📖 使用指南 Quick Guide</summary>
      <div class="guide-content">
        <div class="guide-section">
          <h4>🚀 快速开始</h4>
          <ol>
            <li><strong>新项目：</strong>点击 <strong>初始化项目</strong> 创建完整的 DDD 项目结构</li>
            <li><strong>生成领域：</strong>定义实体名称和字段，自动生成 Entity、Repository、Handler 等</li>
            <li><strong>生成服务：</strong>创建跨领域业务逻辑的应用服务层</li>
            <li><strong>更新依赖：</strong>点击 <strong>更新依赖</strong> 卡片运行 go mod tidy</li>
          </ol>
        </div>
        <div class="guide-section">
          <h4>✨ 新功能</h4>
          <ul>
            <li>✅ 生成领域/服务后 <strong>自动运行 go mod tidy</strong> 下载依赖</li>
            <li>✅ 字段支持 <strong>上下移动</strong> 调整顺序</li>
            <li>✅ 删除领域时 <strong>自动清理</strong> 所有相关文件和注入代码</li>
            <li>✅ 使用 <strong>强制覆盖</strong> 可完全替换现有领域定义</li>
          </ul>
        </div>
        <div class="guide-section">
          <h4>💡 提示</h4>
          <ul>
            <li>所有生成操作都支持 <strong>预览</strong>，可以在实际创建前查看将生成的文件</li>
            <li>勾选 <strong>自动注入到 main.go</strong> 可以自动完成模块注册</li>
            <li>ID、CreatedAt、UpdatedAt 等系统字段会自动生成，无需手动添加</li>
          </ul>
        </div>
      </div>
    </details>

    <div class="status-card" v-if="!loading">
      <div class="status-indicator" :class="layout?.found ? 'found' : 'not-found'">
        {{ layout?.found ? '✓' : '!' }}
      </div>
      <div class="status-content">
        <h3>{{ layout?.found ? '已检测到项目 Project Detected' : '未找到项目 No Project Found' }}</h3>
        <p v-if="layout?.found" class="status-path">{{ layout.module_path }}</p>
        <p v-else class="status-hint">请在包含 go.mod 和 internal/ 目录的项目中运行</p>
      </div>
      <button v-if="projects.length > 1" class="btn-switch" @click="showProjectSelector = !showProjectSelector"
        :title="showProjectSelector ? '关闭项目选择器' : '切换项目'">
        {{ showProjectSelector ? '✕' : '⇄' }}
      </button>
    </div>

    <!-- Project Selector -->
    <div class="project-selector" v-if="showProjectSelector && projects.length > 0">
      <h3>可用项目 Available Projects</h3>
      <div class="project-list">
        <div v-for="project in projects" :key="project.path" class="project-item"
          :class="{ active: project.is_current }" @click="!project.is_current && switchProject(project.path)">
          <div class="project-info">
            <div class="project-name">
              {{ project.name }}
              <span v-if="project.is_current" class="current-badge">当前</span>
            </div>
            <div class="project-module">{{ project.module_path }}</div>
          </div>
          <div class="project-action">
            <span v-if="project.is_current">✓</span>
            <span v-else class="switch-icon">→</span>
          </div>
        </div>
      </div>
      <p class="hint" v-if="switching">正在切换项目...</p>
    </div>

    <div class="cards">
      <RouterLink to="/init" class="card">
        <div class="card-icon">🚀</div>
        <h3>初始化项目 Init Project</h3>
        <p>创建一个新的 Soliton-Go 项目，包含完整的 DDD 结构</p>
      </RouterLink>

      <RouterLink to="/domain" class="card" :class="{ disabled: !layout?.found }">
        <div class="card-icon">📦</div>
        <h3>生成领域 Domain</h3>
        <p>创建 Entity、Repository、Events、Commands、Queries 和 HTTP Handler</p>
      </RouterLink>

      <RouterLink to="/service" class="card" :class="{ disabled: !layout?.found }">
        <div class="card-icon">⚙️</div>
        <h3>生成服务 Service</h3>
        <p>创建跨领域业务逻辑的应用服务层</p>
      </RouterLink>

      <div class="card action-card" :class="{ disabled: !layout?.found }" @click="runGoModTidy">
        <div class="card-icon">📦</div>
        <h3>更新依赖 Dependencies</h3>
        <p v-if="!tidying && !tidyResult">运行 go mod tidy 下载和整理项目依赖</p>
        <p v-else-if="tidying" class="loading">⏳ 正在更新依赖...</p>
        <p v-else-if="tidyResult?.success" class="success">✅ {{ tidyResult.message }}</p>
        <p v-else class="error">❌ {{ tidyResult?.message }}</p>
      </div>
    </div>

    <div class="features">
      <h2>功能特性 Features</h2>
      <ul>
        <li>✨ 可视化字段编辑器，支持拖拽</li>
        <li>👁️ 生成前预览代码 Code Preview</li>
        <li>🔌 自动注入模块到 main.go</li>
        <li>📄 开箱即用的分页支持 Pagination</li>
        <li>🗑️ 可选的软删除 Soft Delete</li>
        <li>📝 CQRS 模式的 Commands & Queries</li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.header {
  text-align: center;
  margin-bottom: 24px;
}

.header h1 {
  font-size: 2.5rem;
  background: linear-gradient(135deg, var(--primary), #8b5cf6);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 8px;
}

.subtitle {
  color: var(--text-muted);
  font-size: 1.1rem;
}

.guide {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1));
  border: 1px solid var(--primary);
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 32px;
}

.guide summary {
  cursor: pointer;
  font-weight: 600;
  font-size: 1.1rem;
  color: var(--text);
  user-select: none;
  list-style: none;
}

.guide summary::-webkit-details-marker {
  display: none;
}

.guide summary::before {
  content: '▶';
  display: inline-block;
  margin-right: 8px;
  transition: transform 0.2s;
}

.guide[open] summary::before {
  transform: rotate(90deg);
}

.guide-content {
  margin-top: 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.guide-section h4 {
  margin-bottom: 12px;
  color: var(--primary);
}

.guide-section ol,
.guide-section ul {
  margin-left: 20px;
  color: var(--text-muted);
  line-height: 1.8;
}

.guide-section li {
  margin-bottom: 8px;
}

.guide-section code {
  background: var(--bg-dark);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
}

.status-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--bg-card);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  border: 1px solid var(--border);
  position: relative;
}

.btn-switch {
  position: absolute;
  right: 20px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid var(--border);
  background: var(--bg-input);
  color: var(--text);
  font-size: 20px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-switch:hover {
  background: var(--primary);
  border-color: var(--primary);
  color: white;
}

.project-selector {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 32px;
  border: 1px solid var(--primary);
}

.project-selector h3 {
  margin-bottom: 16px;
  color: var(--primary);
}

.project-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.project-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: var(--bg-input);
  border-radius: 8px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.project-item:hover:not(.active) {
  border-color: var(--primary);
  background: var(--bg-dark);
}

.project-item.active {
  border-color: var(--success);
  background: rgba(34, 197, 94, 0.1);
  cursor: default;
}

.project-info {
  flex: 1;
}

.project-name {
  font-weight: 600;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-badge {
  display: inline-block;
  padding: 2px 8px;
  background: var(--success);
  color: white;
  border-radius: 4px;
  font-size: 0.75rem;
}

.project-module {
  font-size: 0.85rem;
  color: var(--text-muted);
  font-family: monospace;
}

.project-action {
  font-size: 24px;
  color: var(--text-muted);
}

.project-item:hover:not(.active) .switch-icon {
  color: var(--primary);
}

.status-indicator {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.status-indicator.found {
  background: rgba(34, 197, 94, 0.2);
  color: var(--success);
}

.status-indicator.not-found {
  background: rgba(245, 158, 11, 0.2);
  color: var(--warning);
}

.status-content h3 {
  margin-bottom: 4px;
}

.status-path {
  color: var(--primary);
  font-family: monospace;
}

.status-hint {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-bottom: 40px;
}

.card {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  text-decoration: none;
  color: var(--text);
  border: 1px solid var(--border);
  transition: all 0.2s;
}

.card:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.card.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.card-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.card h3 {
  margin-bottom: 8px;
}

.card p {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.features {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid var(--border);
}

.features h2 {
  margin-bottom: 16px;
}

.features ul {
  list-style: none;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.features li {
  color: var(--text-muted);
}
</style>
