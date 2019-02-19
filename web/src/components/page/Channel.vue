<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>链路列表</h1>
    <button @click="createChannel()" class="btn btn-primary"><i class="fa fa-plus"></i> 添加链路 </button>
    <div class="row" data-toggle="isotope">

        </br>
        <div class="item col-md-12 col-xs-12">
            <div class="panel panel-default">
                <table class="table table-striped margin-none">
                    <thead>
                        <tr>
                            <td class="text-left"><strong>链路名称</strong></td>
                            <td><strong>链路号</strong></td>
                            <td><strong>链路类型</strong></td>
                            <td><strong class="text-right">设备数目</strong></td>
                            <td class="text-right width-200"><strong>编辑</strong></td>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(c, index) in channels">
                            <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                <router-link :to="{ name: 'list', params: {cname:c.Name, cid:c.ID}}" class="text-primary"><strong>{{c.Name}}</strong></router-link>
                            </td>
                            <td>
                                <a class="text-primary"><strong>{{c.ID}}</strong></a>
                            </td>
                            <td>
                                <a class="text-primary"><strong>{{ctype(c.Source)}}</strong></a>
                            </td>
                            <td class="text-danger">{{c.DeviceNum}}</td>
                            <td v-if="ctype(c.Source)=='下行'" class="text-right">
                                <router-link :to="{ name: 'list', params: {cname:c.Name, cid:c.ID}}" class="btn btn-default btn-xs">设备管理</router-link>
                                <button @click="open(c)" class="btn btn-primary btn-xs">修改</button>
                                <button @click="deleteChannel(c)" class="btn btn-danger btn-xs">删除</button>
                            </td>
                            <td v-else class="text-right">
                                <router-link :to="{ name: 'pointtable', params: {cname:c.Name, cid:c.ID}}" class="btn btn-default btn-xs">点表</router-link>
                                <button @click="open(c)" class="btn btn-primary btn-xs">修改</button>
                                <button @click="deleteChannel(c)" class="btn btn-danger btn-xs">删除</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
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
                            <h4 class="modal-title">链路配置</h4>
                            <h5 class="modal-title" style="color: lightgray">新增通道需要添加设备后重启生效</h5>
                            <h5 class="modal-title" style="color: lightgray">通道修改完成后需重启才能生效</h5>
                        </div>
                    </div>
                    <div class="modal-body" data-scrollable>
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">链路名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="链路名称" type="text" v-model="current.Name">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">链路编号</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="链路编号" type="text" v-model.number="current.ID">
                                </div>
                            </div>
                            <!-- <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">链路规约</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="链路规约" type="text" v-model.number="current.Potocol.Name">
                                </div>
                            </div> -->
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">链路规约</label>
                                <div class="col-sm-9">
                                    <v-select v-model="current.picker" :options="current.organizeMembers" placeholder="选择规约以进行操作">
                                    </v-select>
                                </div>
                            </div>

                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">端口类型</label>
                                <!-- <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="端口类型 com 或 tcp" type="text" v-model="current.Port.Ptype">
                                </div> -->
                                <div v-show="current.state == 'create' || current.Port.Ptype == 'com'" class="radio radio-info radio-inline">
                                    <input type="radio" id="com" value="com" v-model="current.Port.Ptype">
                                    <label for="com">COM</label>
                                </div>
                                <div v-show="current.state == 'create' || current.Port.Ptype == 'tcp'" class="radio radio-info radio-inline">
                                    <input type="radio" id="tcp" value="tcp" v-model="current.Port.Ptype">
                                    <label for="tcp">TCP</label>
                                </div>
                            </div>
                        </form>
                        <form v-show="current.Port.Ptype == 'com'" class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail2" class="col-sm-3 control-label">串口名称</label>
                                <div class="col-sm-9">
                                    <!-- <input class="form-control" id="inputEmail2" placeholder="串口名称" type="text" v-model="current.Port.Name"> -->
                                    <v-select v-model="current.comName" :options="current.comNumber" placeholder="请选择端口以进行操作">
                                    </v-select>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">波特率</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="波特率" type="number" v-model.number="current.Port.Baud">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">数据位</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="数据位" type="number" v-model.number="current.Port.Size">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">读取超时</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="读取超时" type="number" v-model.number="current.Port.ReadTimeout">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail5" class="col-sm-3 control-label">停止位</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail5" placeholder="停止位" type="number" v-model.number="current.Port.StopBits">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail6" class="col-sm-3 control-label">校验位</label>
                                <!-- <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail6" placeholder="校验位" type="text" v-model="current.Port.Parity">
                                </div> -->
                                <div class="radio radio-info radio-inline">
                                    <input type="radio" id="n" value="N" v-model="current.Port.Parity">
                                    <label for="n">无校验</label>
                                </div>
                                <div class="radio radio-info radio-inline">
                                    <input type="radio" id="o" value="O" v-model="current.Port.Parity">
                                    <label for="o">奇校验</label>
                                </div>
                                <div class="radio radio-info radio-inline">
                                    <input type="radio" id="e" value="E" v-model="current.Port.Parity">
                                    <label for="e">偶校验</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail5" class="col-sm-3 control-label">crc校验</label>
                                <div class="radio radio-info radio-inline">
                                    <input type="radio" id="h" value="高" v-model="current.Crc">
                                    <label for="h">高字节在前</label>
                                </div>
                                <div class="radio radio-info radio-inline">
                                    <input type="radio" id="l" value="低" v-model="current.Crc">
                                    <label for="l">低字节在前</label>
                                </div>
                            </div>
                        </form>
                        <form v-show="current.Port.Ptype == 'tcp' && current.picker.label != 'JSON'" class="form-horizontal" role="form">
                            <div class="form-group">
                                <div class="form-group">
                                    <label for="inputEmail1" class="col-sm-3 control-label">tcp类型</label>
                                    <div class="col-sm-9">
                                        <!-- <input class="form-control" id="inputEmail1" placeholder="tcp类型 服务端 或 客户端" type="text" v-model="current.Port.Conntype"> -->
                                        <v-select v-model="current.connName" :options="current.connType" placeholder="tcp类型 服务端 或 客户端">
                                        </v-select>
                                    </div>
                                </div>
                                <label for="inputEmail3" class="col-sm-3 control-label">IP地址</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="IP地址" type="text" v-model="current.Port.Ip">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">端口号</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="端口号" type="text" v-model="current.Port.Portnum">
                                </div>
                            </div>
                        </form>

                        <form v-show="current.Port.Ptype == 'tcp' && current.picker.label == 'JSON'" class="form-horizontal" role="form">
                            <div class="form-group">
                                <div class="form-group">
                                    <label for="inputEmail1" class="col-sm-3 control-label">tcp类型</label>
                                    <div class="col-sm-9">
                                        <!-- <input class="form-control" id="inputEmail1" placeholder="tcp类型 服务端 或 客户端" type="text" v-model="current.Port.Conntype"> -->
                                        <v-select v-model="current.connName" :options="current.connType" placeholder="tcp类型 服务端 或 客户端">
                                        </v-select>
                                    </div>
                                </div>
                                <label for="inputEmail3" class="col-sm-3 control-label">IP地址</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="IP地址" type="text" v-model="current.Port.Ip">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">端口号</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="端口号" type="text" v-model="current.Port.Portnum">
                                </div>
                            </div>
                            <div class="form-group" v-if="current.Potocol.Data">
                                <label for="inputEmail4" class="col-sm-3 control-label">账号</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="账号" type="text" v-model="current.Potocol.Data.Key">
                                </div>
                            </div>
                            <div class="form-group" v-if="current.Potocol.Data">
                                <label for="inputEmail4" class="col-sm-3 control-label">密码</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail4" placeholder="密码" type="text" v-model="current.Potocol.Data.Name">
                                </div>
                            </div>
                        </form>


                        <div class="form-group margin-none">
                            <div class="col-sm-offset-3 col-sm-9">
                                <button　@click="save" type="submit" class="btn btn-primary pull-right">保存</button>
                            </div>
                        </div>
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
    name: 'channel',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            show: false,
            channels: [{
                organizeMembers: [{
                    label: "JSON",
                    value: 1
                }, {
                    label: "MODBUS",
                    value: 2
                }, {
                    label: "MODBUSTCP",
                    value: 5
                }, {
                    label: "FORWARD",
                    value: 3
                }, {
                    label: "104",
                    value: 4
                }, {
                    label: "645",
                    value: 5
                }],
                picker: null,
                comNumber: [{
                    label: "COM1",
                    value: 1
                }, {
                    label: "COM2",
                    value: 2
                }, {
                    label: "COM3",
                    value: 3
                }, {
                    label: "COM4",
                    value: 4
                }, {
                    label: "COM5",
                    value: 5
                }, {
                    label: "COM6",
                    value: 6
                }, {
                    label: "COM7",
                    value: 7
                }, {
                    label: "COM8",
                    value: 8
                }, {
                    label: "COM9",
                    value: 9
                }, {
                    label: "COM10",
                    value: 10
                }, {
                    label: "COM11",
                    value: 11
                }, {
                    label: "COM12",
                    value: 12
                }, {
                    label: "COM13",
                    value: 13
                }, {
                    label: "COM14",
                    value: 14
                }, {
                    label: "COM15",
                    value: 15
                }, {
                    label: "COM16",
                    value: 16
                }, {
                    label: "COM17",
                    value: 17
                }, {
                    label: "COM18",
                    value: 18
                }],
                comName: null,
                connType: [{
                    label: "客户端",
                    value: 1
                }, {
                    label: "服务端",
                    value: 2
                }],
                connName: null,
                ID: 0,
                Name: '',
                Crc: '',
                Source: false,
                DeviceNum: 0,
                Potocol: {
                    Name: '',
                    Data: {
                      Key: '',
                      Name: ''
                    }
                },
                Port: {
                    Ptype: '',
                    Conntype: '',
                    Name: '',
                    Baud: 0,
                    Size: 0,
                    ReadTimeout: 0,
                    Parity: '',
                    StopBits: 0,
                    Channel: '',
                    Ip: '',
                    Portnum: ''
                }
            }],
            current: {
                organizeMembers: [{
                    label: "JSON",
                    value: 1
                }, {
                    label: "MODBUS",
                    value: 2
                }, {
                    label: "MODBUSTCP",
                    value: 5
                }, {
                    label: "FORWARD",
                    value: 3
                }, {
                    label: "104",
                    value: 4
                }, {
                    label: "645",
                    value: 5
                }],
                picker: null,
                comNumber: [{
                    label: "COM1",
                    value: 1
                }, {
                    label: "COM2",
                    value: 2
                }, {
                    label: "COM3",
                    value: 3
                }, {
                    label: "COM4",
                    value: 4
                }, {
                    label: "COM5",
                    value: 5
                }, {
                    label: "COM6",
                    value: 6
                }, {
                    label: "COM7",
                    value: 7
                }, {
                    label: "COM8",
                    value: 8
                }, {
                    label: "COM9",
                    value: 9
                }, {
                    label: "COM10",
                    value: 10
                }, {
                    label: "COM11",
                    value: 11
                }, {
                    label: "COM12",
                    value: 12
                }, {
                    label: "COM13",
                    value: 13
                }, {
                    label: "COM14",
                    value: 14
                }, {
                    label: "COM15",
                    value: 15
                }, {
                    label: "COM16",
                    value: 16
                }, {
                    label: "COM17",
                    value: 17
                }, {
                    label: "COM18",
                    value: 18
                }],
                comName: null,
                connType: [{
                    label: "客户端",
                    value: 1
                }, {
                    label: "服务端",
                    value: 2
                }],
                connName: null,
                state: '',
                ID: 0,
                Name: '',
                Crc: '低',
                Potocol: {
                    Name: '',
                    Data: {
                      Key: '',
                      Name: ''
                    }
                },
                DeviceNum: 0,
                Port: {
                    Ptype: '',
                    Conntype: '',
                    Name: '',
                    Baud: 9600,
                    Size: 1024,
                    ReadTimeout: 100,
                    Parity: 'N',
                    StopBits: 2,
                    Channel: '',
                    Ip: '',
                    Portnum: ''
                }
            }
        }
    },
    methods: {
        ctype: function(source) {
            if (source === true) {
                return "下行"
            }
            if (source === false) {
                return "上行"
            }
        },
        open: function(channel) {
          console.log(channel);
            this.current = channel
            this.current.organizeMembers = [{
                label: "JSON",
                value: 1
            }, {
                label: "MODBUS",
                value: 2
            }, {
                label: "FORWARD",
                value: 3
            }, {
                label: "MODBUSTCP",
                value: 5
            }, {
                label: "104",
                value: 4
            }, {
                label: "645",
                value: 5
            }]
            this.current.comNumber = [{
                label: "COM1",
                value: 1
            }, {
                label: "COM2",
                value: 2
            }, {
                label: "COM3",
                value: 3
            }, {
                label: "COM4",
                value: 4
            }, {
                label: "COM5",
                value: 5
            }, {
                label: "COM6",
                value: 6
            }, {
                label: "COM7",
                value: 7
            }, {
                label: "COM8",
                value: 8
            }, {
                label: "COM9",
                value: 9
            }, {
                label: "COM10",
                value: 10
            }, {
                label: "COM11",
                value: 11
            }, {
                label: "COM12",
                value: 12
            }, {
                label: "COM13",
                value: 13
            }, {
                label: "COM14",
                value: 14
            }, {
                label: "COM15",
                value: 15
            }, {
                label: "COM16",
                value: 16
            }, {
                label: "COM17",
                value: 17
            }, {
                label: "COM18",
                value: 18
            }]
            this.current.connType = [{
                label: "客户端",
                value: 0
            }, {
                label: "服务端",
                value: 1
            }]
            this.current.connName = {
                label: channel.Port.Conntype,
                value: null
            }
            this.current.comName = {
                label: channel.Port.Name,
                value: null
            }
            this.current.picker = {
                    label: channel.Potocol.Name,
                    value: null
                }
                // if (channel.Potocol.Name == "MODBUS") {
                //     this.current.picker.value = 2
                // }
            this.show = true
        },
        close: function() {
            this.show = false
            this.current = this.temp()
        },
        save: function() {
            if (this.current.state === "create") {
                // var bean = {
                //     Name: this.current.Name,
                //     ID: this.current.ID
                // }
                if (this.current.Port.Ptype == "com") {
                    console.log("---====");
                    this.current.Potocol.Name = this.current.picker.label
                    this.current.Port.Name = this.current.comName.label
                } else {
                    console.log("====----");
                    console.log("lable", this.current);
                    this.current.Potocol.Name = this.current.picker.label
                    this.current.Port.Conntype = this.current.connName.label
                    this.current.Port.Portnum = this.current.Port.Portnum
                }
                console.log("传下去的内容",this.current);
                this.$http.post('/channels', this.current).then(function(res) {
                    if (res.data.error) {
                        alert(res.data.Message)
                    } else {
                        this.loadChannels()
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("通道异常 通讯管理机终止")
                })

            } else {
                // console.log(this.current.comName.label);
                if (this.current.Port.Ptype == "com") {
                    this.current.Potocol.Name = this.current.picker.label
                    this.current.Port.Name = this.current.comName.label
                } else {
                    this.current.Potocol.Name = this.current.picker.label
                    this.current.Port.Conntype = this.current.connName.label
                    this.current.Port.Portnum = this.current.Port.Portnum
                }
                this.$http.put('/channels', this.current).then(function(res) {

                    if (res.data.error) {
                        alert(res.data.msg)
                    } else {
                        // this.show = false
                        this.loadChannels()
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("修改通道时服务器异常")
                })
            }
        },
        loadChannels: function() { //更新通道
            this.$http.get('/channels').then(function(res) {
                console.log(this.channels);
                if (res.data.error) alert(res.data.msg)
                this.channels = res.data.data
                console.log(this.channels);
            }, function() {
                alert("加载通道时服务器异常")
            })
        },
        createChannel: function() {　 //创建通道
            this.current = this.temp()
            this.show = true
            this.current.state = "create"
        },
        updateChannel: function() { //更新通道
            this.$http.put('/channels', this.current).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                    this.loadChannels()
                }
            }, function() {
                alert("更新通道时服务器异常")
            })
        },
        deleteChannel: function(channel) { //更新通道
            var url = '/channels/' + channel.Name
            this.$http.delete(url).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                    this.loadChannels()
                }
            }, function() {
                alert("删除通道时服务器异常")
            })
        },
        temp: function() {
            return {
                organizeMembers: [{
                    label: "JSON",
                    value: 1
                }, {
                    label: "MODBUS",
                    value: 2
                }, {
                    label: "MODBUSTCP",
                    value: 5
                }, {
                    label: "FORWARD",
                    value: 3
                }, {
                    label: "104",
                    value: 4
                }, {
                    label: "645",
                    value: 5
                }],
                picker: null,
                comNumber: [{
                    label: "COM1",
                    value: 1
                }, {
                    label: "COM2",
                    value: 2
                }, {
                    label: "COM3",
                    value: 3
                }, {
                    label: "COM4",
                    value: 4
                }, {
                    label: "COM5",
                    value: 5
                }, {
                    label: "COM6",
                    value: 6
                }, {
                    label: "COM7",
                    value: 7
                }, {
                    label: "COM8",
                    value: 8
                }, {
                    label: "COM9",
                    value: 9
                }, {
                    label: "COM10",
                    value: 10
                }, {
                    label: "COM11",
                    value: 11
                }, {
                    label: "COM12",
                    value: 12
                }, {
                    label: "COM13",
                    value: 13
                }, {
                    label: "COM14",
                    value: 14
                }, {
                    label: "COM15",
                    value: 15
                }, {
                    label: "COM16",
                    value: 16
                }, {
                    label: "COM17",
                    value: 17
                }, {
                    label: "COM18",
                    value: 18
                }],
                comName: null,
                connType: [{
                    label: "客户端",
                    value: 1
                }, {
                    label: "服务端",
                    value: 2
                }],
                connName: null,
                ID: 0,
                Name: '',
                Crc: '低',
                DeviceNum: 0,
                Potocol: {
                    Name: '',
                    Data: {
                      Key: '',
                      Name: ''
                    }
                },
                Port: {
                    Ptype: '',
                    Conntype: '',
                    Name: '',
                    Baud: 9600,
                    Size: 1024,
                    ReadTimeout: 100,
                    Parity: 'N',
                    StopBits: 2,
                    Channel: '',
                    Ip: '',
                    Portnum: ''
                }
            }
        }
    },
    components: {
        vSelect
    },
    created() {
        this.loadChannels()
    }

}

</script>
