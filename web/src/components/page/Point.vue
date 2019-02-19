<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>设备管理</h1>
    <div class="item col-md-12 col-xs-12">
        <div class="panel panel-default">
            <div class="panel-body">
                <div class="row">
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-primary margin-none">设备名称 ：</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.dname}}</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-primary margin-none">终端号 :</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.did}}</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-primary margin-none">所属通道 :</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.dchannel}}</p>
                    </div>
                </div>
            </div>
        </div>


        <!-- Tabbable Widget -->
        <div class="tabbable">

            <!-- Tabs -->
            <ul class="nav nav-tabs">
                <li :class="{active:state=='yc'}"><a href="#" @click.prevent="state='yc'"><i class="fa fa-fw fa-home"></i> 遥测</a></li>
                <li :class="{active:state=='yx'}"><a href="#" @click.prevent="state='yx'"><i class="fa fa-fw fa-user"></i> 遥信</a></li>
                <li :class="{active:state=='ym'}"><a href="#" @click.prevent="state='ym'"><i class="fa fa-fw fa-envelope"></i> 遥脉</a></li>
                <li :class="{active:state=='yk'}"><a href="#" @click.prevent="state='yk'"><i class="fa fa-fw fa-cog"></i> 遥控</a></li>
                <li :class="{active:state=='yt'}"><a href="#" @click.prevent="state='yt'"><i class="fa fa-fw fa-cog"></i> 遥调</a></li>
            </ul>
            <!-- // END Tabs -->

            <!-- Panes -->
            <div class="tab-content">
                <div class="tab-pane " :class="{active:state=='yc'}">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>数据名称</strong></td>
                                <td class="text-left"><strong>点号</strong></td>
                                <td><strong>数据系数</strong></td>
                                <td><strong>数据</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, index) in ycs">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a href="#" class="text-primary"><strong>{{p.Name}}</strong></a>
                                </td>
                                </strong>
                                <td class="text-left">
                                    <a href="#" class="text-primary"><strong>{{p.ID}}</strong></a>
                                    <td class="text-primary">{{p.Scale}}</td>
                                    <td class="text-success">{{p.Data.toFixed(3)}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yx' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>数据名称</strong></td>
                                <td class="text-left"><strong>点号</strong></td>
                                <td><strong>取反系数</strong></td>
                                <td><strong>数据</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, index) in yxs">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a href="#" class="text-primary"><strong>{{p.Name}}</strong></a>
                                </td>
                                </strong>
                                <td class="text-left">
                                    <a href="#" class="text-primary"><strong>{{p.ID}}</strong></a>
                                    <td class="text-primary">{{p.Scale}}</td>
                                    <td :class="[p.Data ? 'text-danger':'text-success']"><span><i class="icon-stop-fill"></span></td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='ym' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>数据名称</strong></td>
                                <td class="text-left"><strong>点号</strong></td>
                                <td><strong>数据系数</strong></td>
                                <td><strong>数据</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, index) in yms">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a href="#" class="text-primary"><strong>{{p.Name}}</strong></a>
                                </td>
                                </strong>
                                <td class="text-left">
                                    <a href="#" class="text-primary"><strong>{{p.ID}}</strong></a>
                                    <td class="text-primary">{{p.Scale}}</td>
                                    <td class="text-success">{{p.Data.toFixed(2)}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yk' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>数据名称</strong></td>
                                <td class="text-left"><strong>点号</strong></td>
                                <td class="text-right"><strong>控制命令</strong></td>
                                <td　class="text-right"><strong>编辑</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, index) in yks">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a href="#" class="text-primary"><strong>{{p.Name}}</strong></a>
                                </td>
                                </strong>
                                <td class="text-left">
                                    <a href="#" class="text-primary"><strong>{{p.ID}}</strong></a>
                                    <td class="text-right">
                                        <div class="radio radio-info radio-inline">
                                            <input type="radio" :id="'on_'+index" :value="true" v-model.boolean="p.Data">
                                            <label :for="'on_'+index">合闸</label>
                                        </div>
                                        <div class="radio radio-info radio-inline">
                                            <input type="radio" :id="'off_'+index" :value="false" v-model.boolean="p.Data">
                                            <label :for="'off_'+index">分闸</label>
                                        </div>
                                    </td>
                                </td>
                                <td class="text-right">
                                    <button @click="execute(p,'yk','select')" class="btn btn-primary btn-xs">选择</button>
                                    <button @click="execute(p,'yk','exe')" class="btn btn-danger btn-xs">执行</button>
                                    <button @click="execute(p,'yk','cancel')" class="btn btn-success btn-xs">撤销</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yt' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>数据名称</strong></td>
                                <td class="text-left"><strong>点号</strong></td>
                                <td class="text-center"><strong>遥调参数</strong></td>
                                <td　class="text-right"><strong>编辑</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(p, index) in yts">
                                <td class="text-left" style="width:30%"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a href="#" class="text-primary"><strong>{{p.Name}}</strong></a>
                                </td>
                                <td class="text-left" style="width:25%">
                                    <a href="#" class="text-primary"><strong>{{p.ID}}</strong></a>
                                </td>
                                <td class="text-right">
                                    <input type="number" class="form-control form-control-w-50 text-right" placeholder="遥调" v-model.Number="p.Data">
                                </td>
                                <td class="text-right" style="width:25%">
                                    <button @click="execute(p,'yt','select')" class="btn btn-primary btn-xs">选择</button>
                                    <button @click="execute(p,'yt','exe')" class="btn btn-danger btn-xs">执行</button>
                                    <button @click="execute(p,'yt','cancel')" class="btn btn-success btn-xs">撤销</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <!-- // END Panes -->

        </div>

    </div>
</div>
</div>
</div>

</template>

<script>

export default {
    name: 'point',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            dname: '',
            dchannel: '',
            did: 0,
            show: false,
            state: "yc",
            showYc: true,
            ycs: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: null
            }],
            yxs: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: null
            }],
            yms: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: null
            }],
            yks: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: null
            }],
            yts: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: null
            }]
        }
    },
    methods: {
        execute: function(item, cmd, subcmd) {
            var data = item.Data
            console.log(typeof data);
            if (cmd == 'yt') {
                data = Number(data)
            }
            if (cmd == 'yk') {
                data = Boolean(data)
            }
            console.log(typeof data);
            var command = {
                Cmd: cmd,
                Type: subcmd,
                Value: data,
                Dchannel: this.dchannel,
                Did: Number(this.did),
                Pid: item.ID
            }
            console.log(command);
            this.$http.put('/coms/commands', command).then(function(res) {
                if (res.data.error) {
                    window.eventHub.$emit("toast", {
                        level: "error",
                        message: res.data.message,
                        title: "title",
                    })
                } else {
                    window.eventHub.$emit("toast", {
                        level: "success",
                        message: res.data.message,
                        title: "title",
                    })
                }
            }, function() {
                alert("执行命令时服务器异常")
            })
        },
        loadPoint: function() {
            var url = '/coms/devices/' + this.dchannel + '/' + this.did
            this.$http.get(url).then(function(res) {
                if (res.data.error) console.log(res.data.msg)
                this.ycs = res.data.data.Ycs
                this.yxs = res.data.data.Yxs
                this.yms = res.data.data.Yms
                this.yts = res.data.data.Yts
                this.yks = res.data.data.Yks
            }, function() {
                console.log("加载点时服务器异常")
            })
        },
        refresh: function() {
            var url = '/coms/devices/' + this.dchannel + '/' + this.did
            this.$http.get(url).then(function(res) {
                if (res.data.error) console.log(res.data.msg)
                this.ycs = res.data.data.Ycs
                this.yxs = res.data.data.Yxs
                this.yms = res.data.data.Yms
            }, function() {
                console.log("更新数据服务器异常")
            })
        },
        temp: function() {
            return {
                Minute: null,
                Combine: null,
                dname: '',
                dchannel: '',
                did: 0
            }
        }
    },
    created() {
        this.dname = this.$route.params.dname
        this.dchannel = this.$route.params.dchannel
        this.did = this.$route.params.did
        this.loadPoint()
        setInterval(this.refresh, 1000)
    }

}

</script>
