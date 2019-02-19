<style scoped>



</style>

<template>

<div>
    <div class="jumbotron bg-transparent text-center margin-none">
        <h1>权限管理</h1>
        <button @click="openadd" class="btn btn-default"><i class="fa fa-plus"></i> 用户添加 </button>
    </div>
    <div class="row">
        <div class="col-md-10 col-lg-8 col-md-offset-1 col-lg-offset-2">
            <div class="table-responsive">
                <table class="table panel panel-default table-pricing">
                    <thead>
                        <tr>
                            <th class="text-left"><span class="plan">用户名称</span></th>
                            <th class="text-center" style="width: 20%">
                                <span class="plan">浏览权</span>
                            </th>
                            <th class="text-center" style="width: 20%">
                                <span class="plan">操作权</span>
                            </th>
                            <th class="text-center" style="width: 20%">
                                <span class="plan">维护权</span>
                            </th>
                            <th class="text-center" style="width: 20%">
                                <span class="plan">编辑</span>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(u, index) in users">
                            <td class="text-left">
                                <a href="#" class="text-primary"><strong>{{u.Name}}</strong></a>
                            </td>
                            <td v-if="u.Browse" class="text-center text-success"><span class="fa fa-check-circle"></span></td>
                            <td v-else class="text-center text-primary"><span class="fa fa-times-circle"></span></td>
                            <td v-if="u.Operate" class="text-center text-success"><span class="fa fa-check-circle"></span></td>
                            <td v-else class="text-center text-primary"><span class="fa fa-times-circle"></span></td>
                            <td v-if="u.Maintain" class="text-center text-success"><span class="fa fa-check-circle"></span></td>
                            <td v-else class="text-center text-primary"><span class="fa fa-times-circle"></span></td>
                            <td class="text-center text-success">
                                <button @click="open(u)" class="btn btn-sm btn-warning btn-circle">编辑</button>
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
                            <h4 class="modal-title">用户配置</h4>
                        </div>
                    </div>
                    <div class="modal-body" data-scrollable>
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">用户名</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="用户名" type="text" v-model="current.Name">
                                </div>
                            </div>
                            <div v-show="password" class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">用户密码</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail3" placeholder="用户密码" type="text" v-model="current.Password">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">浏览权</label>
                                <td class="text-right" style="padding-left:20px">
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio1" :value="true" name="radioInline" checked="" v-model="current.Browse">
                                        <label for="inlineRadio1">有</label>
                                    </div>
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio2" :value="false" name="radioInline" v-model="current.Browse">
                                        <label for="inlineRadio2">无</label>
                                    </div>
                                </td>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail4" class="col-sm-3 control-label">操作权</label>
                                <td class="text-right" style="padding-left:20px">
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio3" :value="true" name="radioInline1" checked="" v-model="current.Operate">
                                        <label for="inlineRadio3">有</label>
                                    </div>
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio4" :value="false" name="radioInline1" v-model="current.Operate">
                                        <label for="inlineRadio4">无</label>
                                    </div>
                                </td>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail3" class="col-sm-3 control-label">维护权</label>
                                <td class="text-right" style="padding-left:20px">
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio5" :value="true" name="radioInline2" checked="" v-model="current.Maintain">
                                        <label for="inlineRadio5">有</label>
                                    </div>
                                    <div class="radio radio-info radio-inline">
                                        <input type="radio" id="inlineRadio6" :value="false" name="radioInline2" v-model="current.Maintain">
                                        <label for="inlineRadio6">无</label>
                                    </div>
                                </td>
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

</div>

</template>

<script>

export default {
    name: 'user',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            show: false,
            password: false,
            users: [{
                Name: '',
                Browse: null,
                Operate: null,
                Maintain: null
            }],
            current: {
                state: '',
                Name: '',
                Password: '',
                Browse: null,
                Operate: null,
                Maintain: null
            }
        }
    },
    methods: {
        open: function(user) {
            this.current = user
            this.show = true
        },
        openadd: function() {
            this.current = this.temp()
            this.password = true
            this.show = true
            this.current.state = "create"
        },
        close: function() {
            this.password = false
            this.show = false
        },
        save: function() {
            if (this.current.state === "create") {
                this.$http.post('/users', this.current).then(function(res) {
                    if (res.data.error) {
                        alert(res.data.Message)
                    } else {
                        this.loadUsers()
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("服务器异常111")
                })

            } else {
                this.$http.put('/users', this.current).then(function(res) {

                    if (res.data.error) {
                        alert(res.data.msg)
                    } else {
                        this.loadUsers()
                        this.show = false
                        this.current = this.temp()
                    }
                }, function() {
                    alert("修改用户时服务器异常")
                })
            }
        },
        loadUsers: function() {
            this.$http.get('/users').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.users = res.data.data
            }, function() {
                alert("服务器异常111")
            })
        },
        temp: function() {
            return {
                state: '',
                Name: '',
                Browse: null,
                Operate: null,
                Maintain: null
            }
        }
    },
    created() {
        this.loadUsers()
    }
}

</script>
