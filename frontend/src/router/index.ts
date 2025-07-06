import { createRouter, createWebHistory} from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Logs from '../views/Logs.vue'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Logs',
        component: Logs
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router