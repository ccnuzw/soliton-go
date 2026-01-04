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
    ],
})

export default router
