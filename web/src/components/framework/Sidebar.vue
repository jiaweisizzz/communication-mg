<style scoped>



</style>

<template>

<div class="sidebar left sidebar-size-2 sidebar-offset-0 sidebar-skin-blue sidebar-visible-desktop sidebar-visible-mobile">
    <div class="split-vertical">
        <div class="split-vertical-body">
            <div class="split-vertical-cell">
                <div class="tab-pane active" id="sidebar-tabs-menu">
                    <div data-scrollable>
                        <ul class="sidebar-menu sm-icons-left" style="margin-top:18px">
                            <li>
                                <router-link to="/"><i class="fa fa-hdd-o" style="color:white"></i> <span>运行状态</span></router-link>
                            </li>
                        </ul>

                        <ul class="sidebar-menu sm-icons-left" style="margin-top:18px">
                            <li>
                                <router-link to="/channel"><i class="fa fa-link" style="color:white"></i> <span>链路编辑</span></router-link>
                            </li>
                            <li>
                                <router-link to="/template"><i class="fa fa-magic" style="color:white"></i> <span>模板管理</span></router-link>
                            </li>
                            <li>
                                <router-link to="/version"><i class="fa fa-inbox" style="color:white"></i> <span>版本控制</span></router-link>
                            </li>
                        </ul>
                        <!-- <ul class="sidebar-menu sm-icons-left" style="margin-top:18px">
                            <li>
                                <router-link to="/user"><i class="fa fa-user" style="color:white"></i> <span>权限管理</span></router-link>
                            </li>

                            <li>
                                <router-link to="/config"><i class="fa fa-wrench" style="color:white"></i> <span>本机配置</span></router-link>
                            </li>

                            <li>
                                <router-link to="/help"><i class="fa fa-book" style="color:white"></i> <span>说明文档</span></router-link>
                            </li>
                        </ul> -->

                        <div class="category" style="margin-top:18px">CPU / 内存</div>
                        <div class="sidebar-block">
                            <div class="progress">
                                <div class="progress-bar progress-bar-success" role="progressbar" aria-valuenow="{{process.Cpu}}" aria-valuemin="0" aria-valuemax="100" :style="temp">
                                    <!-- <span class="sr-only">{{process.Cpu}}% Complete</span> -->
                                </div>
                            </div>
                            <div class="progress">
                                <div class="progress-bar progress-bar-danger" role="progressbar" aria-valuenow="{{process.Mem}}" aria-valuemin="0" aria-valuemax="100" :style="data">
                                    <!-- <span class="sr-only">74% Complete</span> -->
                                </div>
                            </div>
                        </div>
                        <ul class="sidebar-menu sm-icons-left sm-icons-block" style="margin-top:18px">
                            <li><a href="#" @click.prevent.stop="store.showRestart"><i class="fa fa-refresh" style="color:white"></i> <span>重新启动</span></a></li>
                        </ul>
                    </div>
                </div>
                <!-- // END .tab-pane -->
            </div>
            <!-- // END .split-vertical-cell -->
        </div>
        <!-- // END .split-vertical-body -->
    </div>
</div>

</template>

<script>

export default {
    name: 'sidebar',
    props: {
        store: null
    },
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            temp: '',
            data: '',
            process: {
                Cpu: 0,
                Mem: 0
            }
        }
    },
    methods: {
        loadCpu: function() {
            this.$http.get('/process').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.process = res.data.data
                this.temp = 'width:' + this.process.Cpu + '%;'
                this.data = 'width:' + this.process.Mem + '%;'
            }, function() {

            })
          }
    },
    created() {
        setInterval(this.loadCpu, 1000)
    }
}

</script>
