<script setup lang="ts">
import { RouterView, RouterLink, useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  { path: '/', name: '首页', icon: '🏠' },
  { path: '/init', name: '初始化项目', icon: '🚀' },
  { path: '/domain', name: '领域模块', icon: '📦' },
  { path: '/service', name: '应用服务', icon: '⚙️' },
]
</script>

<template>
  <div class="app">
    <aside class="sidebar">
      <div class="logo">
        <svg class="logo-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64" fill="none">
          <defs>
            <linearGradient id="bgGrad" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" style="stop-color:#6366f1" />
              <stop offset="100%" style="stop-color:#8b5cf6" />
            </linearGradient>
            <linearGradient id="waveGrad" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" style="stop-color:#ffffff;stop-opacity:0.9" />
              <stop offset="100%" style="stop-color:#c7d2fe;stop-opacity:0.9" />
            </linearGradient>
          </defs>
          <circle cx="32" cy="32" r="30" fill="url(#bgGrad)" />
          <path d="M12 32 Q22 22, 32 32 T52 32" stroke="url(#waveGrad)" stroke-width="4" fill="none"
            stroke-linecap="round" />
          <path d="M12 40 Q22 30, 32 40 T52 40" stroke="url(#waveGrad)" stroke-width="3" fill="none"
            stroke-linecap="round" opacity="0.7" />
          <path d="M12 24 Q22 34, 32 24 T52 24" stroke="url(#waveGrad)" stroke-width="3" fill="none"
            stroke-linecap="round" opacity="0.7" />
          <path d="M18 20 L26 32 L18 44" stroke="white" stroke-width="3" fill="none" stroke-linecap="round"
            stroke-linejoin="round" opacity="0.5" />
          <path d="M46 20 L38 32 L46 44" stroke="white" stroke-width="3" fill="none" stroke-linecap="round"
            stroke-linejoin="round" opacity="0.5" />
        </svg>
        <span class="logo-text">Soliton-Gen</span>
      </div>
      <nav class="nav">
        <RouterLink v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item"
          :class="{ active: route.path === item.path }">
          <span class="nav-icon">{{ item.icon }}</span>
          <span class="nav-text">{{ item.name }}</span>
        </RouterLink>
      </nav>
    </aside>
    <main class="main">
      <RouterView />
    </main>
  </div>
</template>

<style>
:root {
  --primary: #6366f1;
  --primary-dark: #4f46e5;
  --bg-dark: #0f172a;
  --bg-card: #1e293b;
  --bg-input: #334155;
  --text: #f1f5f9;
  --text-muted: #94a3b8;
  --border: #475569;
  --success: #22c55e;
  --warning: #f59e0b;
  --error: #ef4444;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: var(--bg-dark);
  color: var(--text);
  line-height: 1.6;
}

.app {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 240px;
  background: var(--bg-card);
  border-right: 1px solid var(--border);
  padding: 20px 0;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px 20px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 20px;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--primary), #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 12px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  color: var(--text-muted);
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item:hover {
  background: var(--bg-input);
  color: var(--text);
}

.nav-item.active {
  background: var(--primary);
  color: white;
}

.nav-icon {
  font-size: 18px;
}

.main {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}

/* Tooltip styles - using data-tooltip attribute */
[data-tooltip] {
  position: relative;
  cursor: help;
}

[data-tooltip]::before {
  content: attr(data-tooltip);
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-8px);
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.95);
  color: white;
  font-size: 0.875rem;
  line-height: 1.4;
  white-space: nowrap;
  border-radius: 6px;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s, transform 0.2s;
  z-index: 1000;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

[data-tooltip]::after {
  content: '';
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(-2px);
  border: 5px solid transparent;
  border-top-color: rgba(0, 0, 0, 0.95);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s, transform 0.2s;
  z-index: 1000;
}

[data-tooltip]:hover::before,
[data-tooltip]:hover::after {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

/* For long tooltips, allow wrapping */
[data-tooltip].tooltip-wrap::before {
  white-space: normal;
  max-width: 250px;
  text-align: center;
}
</style>
