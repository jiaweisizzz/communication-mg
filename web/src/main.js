// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import resource from 'vue-resource'
// import VueSocketio from 'vue-socket.io'
// import socketio from 'socket.io-client'
// Vue.use(VueSocketio, 'http://127.0.0.1:8080/ws')
Vue.use(resource)

Vue.config.productionTip = false

window.eventHub = new Vue()
/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
