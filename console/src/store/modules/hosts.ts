import { Action, Module, Mutation, VuexModule } from 'vuex-module-decorators'
import store from '@/store'
import Service from '@/service'

export interface IHosts {
  List: string[]
  Tag: string[]
}

@Module({ namespaced: true, store, name: 'hosts' })
export class Hosts extends VuexModule implements IHosts  {
  public List: string[] = []
  public Tag: string[] = []

  @Mutation
  setList(list: string[]) {
    this.List = list
  }
  @Mutation
  setTag(tags: string[]) {
    this.Tag = tags
  }

  @Action
  async updateList() {
    this.setList(await Service.getSSHHosts())
    this.setTag(await Service.getTags())
  }
}
