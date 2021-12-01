
import Vue from 'vue'
import VueRouter from 'vue-router'
import { Route } from 'vue-router'
import { Store } from 'vuex'
import { UtilInterface } from '@/assets/utils'
import { ServerInterface } from '@/service'


declare module 'vue/types/vue' {
  interface Vue {
    $router: VueRouter;
    $route: Route;
    $store: Store<any>;
    $utils: UtilInterface,
    service: ServerInterface
  }
}