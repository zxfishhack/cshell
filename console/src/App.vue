<template>
  <div id="app">
    <app-layout>
      <router-view />
    </app-layout>
  </div>
</template>

<script lang="ts">
  import { Component, Vue, Watch } from 'vue-property-decorator'
  import settings from '@/settings'
  import Service from '@/service'
  import { getModule } from 'vuex-module-decorators'
  import { Hosts } from '@/store/modules/hosts'

  @Component({
    name: 'App',
    components: {}
  })
  export default class extends Vue {
  }
</script>
