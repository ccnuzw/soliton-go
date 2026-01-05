import { createApp, h } from 'vue'
import Toast, { type ToastProps } from './components/Toast.vue'

export function showToast(props: ToastProps) {
    const container = document.createElement('div')
    document.body.appendChild(container)

    const app = createApp({
        render() {
            return h(Toast, {
                ...props,
                onClose: () => {
                    app.unmount()
                    document.body.removeChild(container)
                }
            })
        }
    })

    app.mount(container)
}

export function showSuccess(message: string) {
    showToast({ message, type: 'success' })
}

export function showError(message: string) {
    showToast({ message, type: 'error' })
}

export function showWarning(message: string) {
    showToast({ message, type: 'warning' })
}

export function showInfo(message: string) {
    showToast({ message, type: 'info' })
}
