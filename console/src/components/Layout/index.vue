<template>
  <el-container class="container">
    <el-aside class="pretty-scrollbar">
      <el-menu :router="true" :unique-opened="true">
        <el-menu-item :index="`/new`">新增</el-menu-item>
        <el-submenu :index="`/all`">
          <template slot="title">所有 ({{Hosts?.length}})</template>
          <el-menu-item v-for="item in Hosts" :key="item" :index="`/edit/${item}`">{{item}}</el-menu-item>
        </el-submenu>
        <el-submenu v-for="tag in Tags" :key="`key-${tag}`" :index="`/tag/${tag}`">
          <template slot="title">{{tag}} ({{serverByTag[tag]?.length}})</template>
          <el-menu-item v-for="item in serverByTag[tag]" :key="`${tag}-${item}`" :index="`/edit/${item}`">{{item}}</el-menu-item>
        </el-submenu>
      </el-menu>
    </el-aside>
    <el-main>
      <slot></slot>
    </el-main>
  </el-container>
</template>

<script lang="ts">
import { Vue, Component, Watch } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import { Hosts } from '@/store/modules/hosts'
import { Debounce } from 'vue-debounce-decorator'
import Service from '@/service'

@Component({})
export default class AppLayout extends Vue {
  serverByTag: {[index:string]: string[]} = {}
  get Tags() : string[] {
    return this.$store.state.hosts.Tag
  }
  get Hosts() : string[] {
    return this.$store.state.hosts.List
  }

  @Watch('Tags')
  @Watch('Hosts')
  @Debounce(250)
  async updateList() {
    const waitList: Promise<string[]>[] = []
    let serverByTag: {[index:string]: string[]} = {}
    this.Tags.forEach(tag => {
      const v = Service.getSSHHostsByTag(tag)
      waitList.push(v)
      v.then(res => {
        serverByTag[tag] = res
      })
    })
    await Promise.all(waitList)
    this.serverByTag = serverByTag
    console.log(this.serverByTag)
  }

  async mounted() {
    await getModule(Hosts, this.$store).updateList()
  }
}
</script>

<style lang="scss">
.container {
  height: 100%;
}

</style>
