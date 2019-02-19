<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>设备列表</h1>
    <button @click="createDevice()" class="btn btn-primary"><i class="fa fa-plus"></i> 添加设备 </button>
    <div class="row" data-toggle="isotope">

        </br>
        <div class="row" data-toggle="isotope">
            <div class="item col-md-12 col-xs-12">
                <div class="panel panel-default">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>设备名称</strong></td>
                                <td><strong>设备地址</strong></td>
                                <td class="text-center"><strong>终端号</strong></td>
                                <td class="text-right width-200"><strong>编辑</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(d, index) in devices">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{d.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{d.ID}}</strong></a>
                                </td>
                                <td class="text-center">
                                    <a class="text-primary"><strong>{{d.Tnumber}}</strong></a>
                                </td>
                                <td class="text-right">
                                    <button @click="open(d)" class="btn btn-primary btn-xs">修改</button>
                                    <button @click="deleteDevice(d)" class="btn btn-danger btn-xs">删除</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
    <div v-if="show" class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>


                        <div class="text-center">
                            <h4 class="modal-title">设备配置</h4>
                            <h5 class="modal-title" style="color: lightgray">设备全部添加完成后需重启生效</h5>
                            <h5 class="modal-title" style="color: lightgray">设备修改完成后需重启生效</h5>
                        </div>
                    </div>
                    <div class="modal-body" data-scrollable>
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">设备名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="设备名称" type="text" v-model="current.Name">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">设备地址</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="设备地址" type="text" v-model.number="current.ID">
                                </div>
                            </div>
                            <div v-if="change" class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">终端号(唯一)</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="终端号" type="text" v-model.number="current.Tnumber">
                                </div>
                            </div>
                            <!-- <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">所属通道</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="所属通道" type="text" v-model="current.Channel">
                                </div>
                            </div> -->
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">选择模板</label>
                                <div class="col-sm-9">
                                    <!-- <input class="form-control" id="inputEmail3" placeholder="选择模板" type="text" v-model="current.Template"> -->
                                    <v-select v-model="current.picker" :options="current.organizeMembers" placeholder="选择模板"></v-select>
                                </div>
                            </div>
                            <div class="form-group margin-none">
                                <div class="col-sm-offset-3 col-sm-9">
                                    <button　@click="save" type="submit" class="btn btn-primary pull-right">保存</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

</template>

<script>

import vSelect from 'vue-select'

export default {
    name: 'list',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            cname: '',
            cid: 0,
            show: false,
            change: false,
            mobans: [],
            devices: [{
                ID: 0,
                Tnumber: 0,
                Name: '',
                Channel: '',
                ChannelID: 0
            }],
            current: {
                organizeMembers: [],
                picker: null,
                state: '',
                ID: 0,
                Tnumber: 0,
                Name: '',
                Channel: '',
                ChannelID: 0,
                Template: ''
            }
        }
    },
    methods: {
        loadMoban: function() {
            this.$http.get('/templates').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.mobans = res.data.data
                for (var i = 0; i < this.mobans.length; i++) {
                    this.current.organizeMembers.push({
                        label: this.mobans[i].Name,
                        value: i
                    })
                }
            }, function() {
                alert("加载模板时服务器异常")
            })
        },
        open: function(device) {
            this.current = this.temp()
            for (var variable in device) {
                if (device.hasOwnProperty(variable)) {
                    this.current[variable] = device[variable]
                }
            }
            for (var i = 0; i < this.mobans.length; i++) {
                this.current.organizeMembers.push({
                    label: this.mobans[i].Name,
                    value: i
                })
            }
            this.change = true
            this.show = true
        },
        close: function() {
            this.change = false
            this.show = false
                // this.current = this.temp()
        },
        save: function() {
            if (this.current.state === "create") {
                this.current.Template = this.current.picker.label
                this.current.Channel = this.cname
                this.current.ChannelID = this.cid
                this.$http.post('/devices', this.current).then(function(res) {
                    if (res.data.error) {
                        alert(res.data.Message)
                    } else {
                        this.loadDevices()
                        this.change = false
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("服务器异常333")
                })

            } else {
              this.current.Template = this.current.picker.label
              this.current.Channel = this.cname
              this.current.ChannelID = this.cid
                this.$http.put('/devices', this.current).then(function(res) {
                    if (res.data.error) {
                        alert(res.data.msg)
                    } else {
                        // this.show = false
                        this.loadDevices()
                        this.change = false
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("修改列表时服务器异常")
                })
            }
        },
        loadDevices: function() {
            var url = '/devices/' + this.cname
            this.$http.get(url).then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.devices = res.data.data
                console.log(this.devices);
            }, function() {
                alert("设备异常 通讯管理机终止")
            })
        },
        createDevice: function() {　 //创建设备
            this.current = this.temp()
            for (var i = 0; i < this.mobans.length; i++) {
                this.current.organizeMembers.push({
                    label: this.mobans[i].Name,
                    value: i
                })
            }
            this.show = true
            this.current.state = "create"
        },
        deleteDevice: function(device) {
            console.log(device.Channel, device.ID);
            var url = '/devices/' + device.Channel + '/' + device.ID
            this.$http.delete(url).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                    this.loadDevices()
                }
            }, function() {
                alert("服务器异常222")
            })
        },
        temp: function() {
            return {
                organizeMembers: [],
                picker: null,
                ID: 0,
                Name: '',
                Channel: '',
                ChannelID: 0,
                Template: ''
            }
        }
    },
    components: {
        vSelect
    },
    created() {
        this.cname = this.$route.params.cname
        this.cid = this.$route.params.cid
        this.loadMoban()
        this.loadDevices()
    }

}

</script>
