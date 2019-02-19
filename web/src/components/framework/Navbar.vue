<style scoped>



</style>

<template>

<div class="navbar navbar-default navbar-fixed-top" role="navigation">
    <div class="container-fluid">
        <div class="navbar-header">
            <router-link to="/home" class="navbar-brand">通讯管理机</router-link>
        </div>
        <div class="navbar-collapse collapse" id="collapse">
            <ul class="nav navbar-nav navbar-right">
                <li class="dropdown notifications">
                    <a @click="jump" class="dropdown-toggle" data-toggle="dropdown"> <i class="fa fa-bell-o"></i> <span v-if="show" class="badge badge-danger">新消息</span> </a>
                </li>
                <li class="dropdown user">
                    <a href="#" @click.prevent.stop="store.showMessage" class="dropdown-toggle">
                        <i class="fa fa-headphones"></i> 报文解析
                    </a>
                </li>
                <!-- <li class="dropdown user">
                    <a href="#" @click.prevent.stop="store.showPermission">
                        <i class="fa fa-unlock-alt"></i> 超级管理员
                    </a>
                </li> -->
                <!-- // END user -->
            </ul>
        </div>
    </div>
</div>

</template>

<script>

export default {
    name: 'navbar',
    props: {
        store: null
    },
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            show: false
        }
    },

    methods: {

        alarm: function(msg) {
            console.log(msg);
            this.show = true
        },
        jump: function() {
            this.show = false
            this.$router.push('/notice')
        }
    },
    created() {
        window.eventHub.$on('onnotice', this.alarm)
    }
}

</script>
