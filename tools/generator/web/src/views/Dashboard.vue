<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api, type ProjectLayout } from '../api'

const layout = ref<ProjectLayout | null>(null)
const loading = ref(true)
const showGuide = ref(false)

onMounted(async () => {
  try {
    layout.value = await api.getLayout()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
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
            <li>如果是新项目，点击 <strong>初始化项目</strong> 创建项目骨架</li>
            <li>在项目目录中运行 <code>soliton-gen serve</code> 启动 Web GUI</li>
            <li>使用 <strong>生成领域</strong> 创建业务实体和相关代码</li>
            <li>使用 <strong>生成服务</strong> 创建应用服务层</li>
          </ol>
        </div>
        <div class="guide-section">
          <h4>💡 提示</h4>
          <ul>
            <li>所有生成操作都支持 <strong>预览</strong>，可以在实际创建前查看将生成的文件</li>
            <li>勾选 <strong>自动注入到 main.go</strong> 可以自动完成模块注册</li>
            <li>使用 <strong>强制覆盖</strong> 选项可以更新已存在的文件</li>
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
  max-width: 900px;
  margin: 0 auto;
}

.header {
  text-align: center;
  margin-bottom: 24px;
}

.header h1 {
  font-size: 2.5rem;
  background: linear-gradient(135deg, var(--primary), #8b5cf6);
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
  margin-bottom: 32px;
  border: 1px solid var(--border);
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
  grid-template-columns: repeat(3, 1fr);
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
