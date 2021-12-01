import Vue from 'vue'
import Vuex from 'vuex'
import { Hosts, IHosts } from '@/store/modules/hosts'
Vue.use(Vuex)

export interface IRootState {
  hosts: IHosts
}

const store = new Vuex.Store<IRootState>({
  modules: {
    hosts: Hosts
  }
})
export default store
