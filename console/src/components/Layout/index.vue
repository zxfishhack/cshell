<template>
  <el-container class="container">
    <el-aside>
      <el-aside>
        <el-menu :router="true">
          <el-menu-item v-for="item in Hosts" :key="item" :index="`/edit/${item}`">
            <template slot="title">
              {{item}}
            </template>
          </el-menu-item>
          <el-menu-item :index="`/new`">新增</el-menu-item>
        </el-menu>
      </el-aside>
    </el-aside>
    <el-main>
      <slot></slot>
    </el-main>
  </el-container>
</template>

<script lang="ts">
import { Vue, Watch, Prop, Component } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import settings from '@/settings'
import { displayRouters } from '@/router/routers'

@Component({})
export default class AppLayout extends Vue {
  get Hosts() {
    return this.$store.state.hosts.List
  }
}
</script>

<style lang="scss">
.container {
  height: 100%;
}
.current-user {
  height: 60px;
  line-height: 60px;
  position: absolute;
  right: 10px;
}
</style>
