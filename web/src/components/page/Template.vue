<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>模板管理</h1>
    <button @click="createMoban()" class="btn btn-primary"><i class="fa fa-plus"></i> 添加模板 </button>
    <div class="row" data-toggle="isotope">

        </br>
        <div class="item col-md-12 col-xs-12">
            <div class="panel panel-default">
                <table class="table table-striped margin-none">
                    <thead>
                        <tr>
                            <td class="text-left"><strong>模板名称</strong></td>
                            <td class="text-left"><strong>规约</strong></td>
                            <td><strong>说明</strong></td>
                            <td class="text-right width-200"><strong>编辑</strong></td>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(t, index) in mobans">
                            <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                                <router-link :to="{ name: 'templatelist', params: {tname:t.Name}}" class="text-primary"><strong>{{t.Name}}</strong></router-link>
                            </td>
                            <td class="text-left">
                                </strong>
                                <a href="#" class="text-primary"><strong>{{t.Potocol}}</strong></a>
                            </td>
                            <td class="text-danger">说明</td>
                            <td class="text-right">
                                <router-link :to="{ name: 'templatelist', params: {tname:t.Name, tid:t.ID, tchannel:t.Channel, tpotocol:t.Potocol}}" class="btn btn-default btn-xs">查看</router-link>
                                <button @click="deleteTemplate(t)" class="btn btn-danger btn-xs">删除</button>
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
                            <h4 class="modal-title">模板配置</h4>
                        </div>
                    </div>
                    <div class="modal-body">
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">模板名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="模板名称" type="text" v-model="current.Name">
                                </div>
                            </div>
                            <!-- <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">终端号</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="终端号" type="text" v-model.number="current.ID">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">所属通道</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="所属通道" type="text" v-model="current.Channel">
                                </div>
                            </div> -->
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">规约选择</label>
                                <div class="col-sm-9">
                                    <v-select v-model="current.picker" :options="current.organizeMembers" placeholder="选择规约以进行操作"></v-select>
                                </div>
                            </div>
                            <div class="form-group margin-none">
                                <div class="col-sm-offset-3 col-sm-9">
                                    <button　@click="save" type="submit" class="btn btn-primary pull-right">保存</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div </div>
            </div>
        </div>
    </div>

</template>

<script>

import vSelect from 'vue-select'
export default {
    name: 'template',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            show: false,
            mobans: [{
                organizeMembers: [{
                    label: "MODBUS",
                    value: 1
                }, {
                    label: "MODBUSTCP",
                    value: 2
                }],
                picker: null,
                // ID: 0,
                Name: '',
                    // Channel: '',
                Potocol: ''
            }],
            current: {
                organizeMembers: [{
                    label: "MODBUS",
                    value: 1
                }, {
                    label: "MODBUSTCP",
                    value: 2
                }],
                picker: null,
                // ID: 0,
                Name: '',
                    // Channel: '',
                Potocol: ''
            }
        }
    },
    methods: {
        deleteTemplate: function(template) {
          var url = '/templates/' + template.Name
          this.$http.delete(url).then(function(res) {
              if (res.data.error) {
                  alert(res.data.msg)
              } else {
                  this.loadMoban()
              }
          }, function() {
              alert("删除模板时服务器异常")
          })
        },
        createMoban: function() {　 //创建通道
            this.current = this.temp()
            this.show = true
                // this.current.state = "create"
        },
        close: function() {
            this.show = false
            this.current = this.temp()
        },
        save: function() {
            console.log(this.current.picker.label);
            this.current.Potocol = this.current.picker.label
            console.log(this.current);
            this.$http.post('/templates', this.current).then(function(res) {
                if (res.data.error) {
                    alert(res.data.Message)
                } else {
                    this.loadMoban()
                    this.show = false
                    this.current = this.temp()
                }
            }, function() {
                alert("模板服务异常 通讯管理机终止")
            })
        },
        loadMoban: function() {
            this.$http.get('/templates').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.mobans = res.data.data
            }, function() {
                alert("加载模板时服务器异常")
            })
        },
        temp: function() {
            return {
                organizeMembers: [{
                    label: "MODBUS",
                    value: 1
                }, {
                    label: "MODBUSTCP",
                    value: 2
                }],
                picker: null,
                // ID: 0,
                Name: '',
                    // Channel: '',
                Potocol: ''
            }
        }
    },
    components: {
        vSelect
    },
    created() {
        this.loadMoban()
    }
}

</script>
