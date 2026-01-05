import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'dashboard',
            component: () => import('./views/Dashboard.vue'),
        },
        {
            path: '/init',
            name: 'init',
            component: () => import('./views/InitWizard.vue'),
        },
        {
            path: '/domain',
            name: 'domain',
            component: () => import('./views/DomainEditor.vue'),
        },
        {
            path: '/service',
            name: 'service',
            component: () => import('./views/ServiceEditor.vue'),
        },
        {
            path: '/migration',
            name: 'migration',
            component: () => import('./views/MigrationCenter.vue'),
        },
        {
            path: '/ddd',
            name: 'ddd',
            component: () => import('./views/DddEditor.vue'),
        },
    ],
})

export default router
