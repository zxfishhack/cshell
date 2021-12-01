import Vue from 'vue'
import VueRouter from 'vue-router'
import { backgroundRouters, displayRouters } from '@/router/routers'

Vue.use(VueRouter)

const routers = [
  ...backgroundRouters,
  ...displayRouters
]

const router = new VueRouter({
  mode: 'history',
  routes: routers
})

export default router
