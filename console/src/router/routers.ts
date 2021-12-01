import { RouteConfig } from 'vue-router'
import Editor from '@/views/Editor.vue'

export const backgroundRouters: RouteConfig[] = [
]

export const displayRouters: RouteConfig[] = [
  {
    path: '/',
    name: '',
    meta: {},
    component: () => import('@/views/Dashboard.vue')
  },
  {
    path: '/edit/:id',
    name: '配置修改',
    component: Editor
  },
  {
    path: '/new',
    name: '配置新增',
    component: Editor
  }
]
