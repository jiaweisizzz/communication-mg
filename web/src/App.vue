<style>

#app {}

</style>

<template>

<div id="app" class="st-container show-sidebar sidebar-l2">

    <navbar :store="store"></navbar>

    <sidebar :store="store"></sidebar>
    <div class="split-vertical">
        <div class="split-vertical-body">
            <div class="split-vertical-cell">
                <div id="content">
                    <div class="container-fluid">
                        <router-view></router-view>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <foot></foot>

    <message :store="store"></message>

    <permission :store="store"></permission>

    <restart :store="store"></restart>

</div>

</template>

<script>

import Navbar from './components/framework/Navbar.vue'
import Sidebar from './components/framework/Sidebar.vue'
import Foot from './components/framework/Footer.vue'
import Message from './components/widget/Message'
import Permission from './components/widget/Permission'
import Restart from './components/widget/Restart'
import miniToastr from 'mini-toastr'

export default {
    name: 'app',
    data() {
        return {
            ws: null,
            store: {
                debug: true,
                state: {
                    message: false,
                    permission: false,
                    restart: false
                },
                showMessage() {
                    this.debug && console.log("报文显示")
                    this.state.message = true
                    this.state.permission = false
                    this.state.restart = false
                },
                hideMessage() {
                    this.debug && console.log("报文隐藏")
                    this.state.message = false
                    this.state.permission = false
                    this.state.restart = false
                },
                showPermission() {
                    this.debug && console.log("权限不足显示")
                    this.state.message = false
                    this.state.permission = true
                    this.state.restart = false
                },
                hidePermission() {
                    this.debug && console.log("权限不足隐藏")
                    this.state.message = false
                    this.state.permission = false
                    this.state.restart = false
                },
                showRestart() {
                    this.debug && console.log("重启显示")
                    this.state.message = false
                    this.state.permission = false
                    this.state.restart = true
                },
                hideRestart() {
                    this.debug && console.log("重启隐藏")
                    this.state.message = false
                    this.state.permission = false
                    this.state.restart = false
                }
            }
        }
    },
    methods: {
        websocketmessages: function(e) {
            console.log("this----", window.onload);
            var tar = JSON.parse(e.data)
            switch (tar.Cmd) {
                case 'message':
                    window.eventHub.$emit('onmessage', tar.Data)
                    break;
                case 'notice':
                    window.eventHub.$emit('onnotice', tar.Data)
                    break;
            }

        },
        send: function(e) {
            this.ws.send(e)
        },
        toast: function(message) {
          switch (message.level) {
            case 'error':
              miniToastr.error(message.message)
              break;
              case 'success':
                miniToastr.success(message.message)
                break;
                case 'warn':
                  miniToastr.warn(message.message)
                  break;
            default:
              miniToastr.info(message.message)
          }
        }
    },
    created() {
        this.ws = new WebSocket(`ws://${window.location.hostname}:8080/ws`)
        this.ws.onmessage = this.websocketmessages
        window.eventHub.$on('app', this.send)
        console.log("this.$route.params====",this.$route);
        miniToastr.init()
        window.eventHub.$on("toast", this.toast)
    },
    components: {
        Navbar,
        Sidebar,
        Foot,
        Message,
        Permission,
        Restart
    }
}

</script>
