<template>
  <div class="app-container">
    <slot></slot>
  </div>
</template>

<script lang="ts">
    import { Component, Vue } from 'vue-property-decorator'

    @Component({
        name: 'common',
        components: {}
    })
    export default class CommonLayout extends Vue {
      get mainHeight() {
        return document.body.scrollHeight
      }
      created() {
      }

      mounted() {
      }
    }
</script>

<style scoped lang="scss">

</style>
