<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>消息列表</h1>
    <div class="item col-md-12 col-xs-12">
        <div class="panel panel-default">
            <div class="panel-body">
                <div class="row">
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-primary margin-none">报警总数 ：</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-danger margin-none">{{msgs.length}}</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">查询报警 :</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <vue-datepicker v-model="time"></vue-datepicker>
                    </div>
                </div>
            </div>
        </div>
        </br>

        <div class="row" data-toggle="isotope">
            <div class="item col-md-12 col-xs-12">
                <div v-if="msgs.length!=0" class="panel panel-default">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>时间</strong></td>
                                <td class="text-left"><strong>设备名称</strong></td>
                                <td><strong>数据名称</strong></td>
                                <td class="text-right width-200"><strong>类型</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(m, index) in msgs">
                                <td class="text-left text-success"><strong class="text-muted ">{{index+1}}.</strong>
                                    <strong>{{m.Time}}</strong>
                                </td>
                                <td class="text-left">
                                    <strong>{{m.Dname}}</strong>
                                </td>
                                <td class="text-primary margin-none">
                                    <strong>{{m.Pname}}</strong>
                                </td>
                                <td v-bind:class="m.Type=='恢复'? 'text-success': 'text-danger'" class="text-right ">
                                    <strong>{{m.Type}}</strong>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div v-else>
                    <h2>暂无报警消息</h2>
                </div>
            </div>
        </div>
    </div>
</div>

</template>

<script>

import VueDatepicker from 'vue-datepicker-local'

export default {
    name: 'notice',
    data() {
        return {
            msgs: [],
            // time: new Date()
            time: (function() {
                return (new Date()).getFullYear() + '-' + (String((new Date()).getMonth() + 1).length === 1 ? '0' + ((new Date()).getMonth() + 1) : ((new Date()).getMonth() + 1)) + '-' + (String((new Date()).getDate()).length === 1 ? '0' + (new Date()).getDate() : (new Date()).getDate())
            })(),
        }
    },
    watch: {
        time: function(val) {
            console.log(val);
            let time = new Date(val).getFullYear() + '-' + (String((new Date(val)).getMonth() + 1).length === 1 ? '0' + ((new Date(val)).getMonth() + 1) : ((new Date(val)).getMonth() + 1)) + '-' + (String((new Date(val)).getDate()).length === 1 ? '0' + (new Date(val)).getDate() : (new Date(val)).getDate())
            console.log('time====', time);
            var url = '/message/' + time
            this.$http.get(url).then(function(res) {
                if (res.data.error) alert(res.data.msg)
                var msg = res.data.data
                if (msg === null) {
                    this.msgs = []
                } else {
                    this.msgs = []
                    var strs = new Array()
                    var array = new Array()
                    strs = msg.split("\n")
                    for (var i = 0; i < strs.length - 1; i++) {
                        var message = {}
                        array = strs[i].split(",")
                        message.Time = array[0]
                        message.Dname = array[1]
                        message.Pname = array[2]
                        message.Type = array[3]
                        this.msgs.push(message)
                    }
                }
            }, function() {
                alert("加载报警消息服务器异常")
            })
        }
    },
    methods: {
        messageload: function() {
            this.$http.get('/message').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                var msg = res.data.data
                var strs = new Array()
                var array = new Array()
                strs = msg.split("\n")
                for (var i = 0; i < strs.length - 1; i++) {
                    var message = {}
                    array = strs[i].split(",")
                    message.Time = array[0]
                    message.Dname = array[1]
                    message.Pname = array[2]
                    message.Type = array[3]
                    this.msgs.push(message)
                }
            }, function() {
                alert("消息模块异常 通讯管理机终止")
            })
        },
        newmessageload: function() {
            this.$http.get('/message').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                var msg = res.data.data
                var strs = new Array()
                var array = new Array()
                this.msgs = []
                strs = msg.split("\n")
                for (var i = 0; i < strs.length - 1; i++) {
                    var message = {}
                    array = strs[i].split(",")
                    message.Time = array[0]
                    message.Dname = array[1]
                    message.Pname = array[2]
                    message.Type = array[3]
                    this.msgs.push(message)
                }
            }, function() {
                alert("消息服务异常 通讯管理机终止")
            })
        }
    },
    created() {
        window.eventHub.$on('onnotice', this.newmessageload)
        this.messageload()
    },
    components: {
        VueDatepicker
    }
}

</script>
