import Vue from 'vue'
import App from '@/App.vue'
import store from '@/store'
import router from '@/router'
import ElementUI from 'element-ui'
import SvgIcon from 'vue-svgicon'

import '@/assets/icon/components'
import 'normalize.css'
import '@/assets/styles/index.scss'
import '@/components/MilkVue'
import 'element-ui/lib/theme-chalk/index.css'
import '@/assets/styles/theme/default.scss'

Vue.config.productionTip = false

Vue.use(SvgIcon, {
  tagName: 'svg-icon',
  defaultWidth: '1em',
  defaultHeight: '1em'
})

Vue.use(ElementUI, {
  size: 'small'
})

new Vue({
  store,
  router,
  render: h => h(App)
}).$mount('#app')
