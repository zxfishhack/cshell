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

  @Component({
    name: 'App',
    components: {}
  })
  export default class extends Vue {
    async created() {
      const hosts = await Service.getSSHHosts()
      this.$store.state.hosts.List = hosts
    }
    mounted() {
      document.title = settings.title
    }
  }
</script>
