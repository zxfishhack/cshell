import { Action, Module, Mutation, VuexModule } from 'vuex-module-decorators'
import store from '@/store'
import Service from '@/service'

export interface IHosts {
  List: string[]
}

@Module({ namespaced: true, store, name: 'hosts' })
export class Hosts extends VuexModule implements IHosts  {
  public List: string[] = []
}
