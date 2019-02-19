<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>运行状态</h1>
    <div class="item col-md-12 col-xs-12">
        <div class="panel panel-default">
            <div class="panel-body">
                <div class="row">
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-primary margin-none">设备总数 ：</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{devices.length}}</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-success margin-none">运行数：</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{runNum}}</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-body-2  text-danger margin-none">停止数：</p>
                    </div>
                    <div class="col-md-2 text-center">
                        <p class="text-title  text-success margin-none">{{devices.length - runNum}}</p>
                    </div>
                </div>
            </div>
        </div>
        </br>
        <div class="row" data-toggle="isotope">
            <div class="item col-md-12 col-xs-12">
                <div class="panel panel-default">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left"><strong>设备名称</strong></td>
                                <td><strong>终端号</strong></td>
                                <td><strong>运行状态</strong></td>
                                <td class="text-right width-200"><strong>编辑</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(d, index) in devices">
                                <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{d.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{d.Tnumber}}</strong></a>
                                </td>
                                <td class="text-primary margin-none">
                                    <div v-if="d.Status" class="label label-success">运行中</div>
                                    <div v-if="!d.Status" class="label label-danger">已停止</div>
                                </td>
                                <td class="text-right">
                                    <router-link :to="{ name: 'point', params: {dname:d.Name, dchannel:d.Channel, did:d.ID}}" class="btn btn-primary btn-xs">查看</router-link>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
    <!-- <div v-if="show" class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>


                        <div class="text-center">
                            <h4 class="modal-title">链路配置</h4>
                        </div>
                    </div>
                    <div class="modal-body" data-scrollable>
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">链路名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="链路名称" type="text" :value="current.name">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">串口名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="串口名称" type="text" :value="current.port">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">波特率</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="波特率" type="text">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">数据位</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="数据位" type="text">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">停止位</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="停止位" type="text">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">校验位</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="校验位" type="email">
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
    </div> -->
</div>

</template>

<script>

export default {
    name: 'home',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            cname: '',
            show: false,
            devices: [],
            current: {
                state: '',
                ID: 0,
                Tnumber: 0,
                Name: '',
                Channel: '',
                Status: null
            }
        }
    },
    computed: {
        runNum: function() {
            var num = 0
            for (var i in this.devices) {
                if (this.devices[i].Status == true) {
                    num++
                }
            }
            return num
        },
        // orderedUsers: function () {
        //     return this.devices.orderBy(this.users, 'name')
        // }
    },
    methods: {
        look: function() {
            this.$router.go('/Template')
        },
        loadAllDevices: function() {
            this.$http.get('/coms/devices').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                console.log(res.data.data);
                this.devices = res.data.data
                this.devices.sort(function(a,b){
                  if (a.Tnumber - b.Tnumber === 0) {
                    var aName = a.Name.toLowerCase()
                    var bName = b.Name.toLowerCase()
                      if (aName < bName) return -1
                      if (aName > bName) return 1
                      return 0
                  }
                  return a.Tnumber - b.Tnumber
                })
            }, function() {
                alert("通讯管理机异常终止")
            })
        },
        refresh: function() {
            this.loadAllDevices()
        }
    },
    created() {
        this.loadAllDevices()
        this.timer = setInterval(this.refresh, 1000)
    },
    beforeDestory() {
      clearInterval(this.timer)
    }

}

</script>
