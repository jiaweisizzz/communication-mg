<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>版本管理</h1>
    <button @click="pack()" class="btn btn-primary"><i class="fa fa-plus"></i> 新建版本 </button>
    <div class="row" data-toggle="isotope">

        </br>
        <div class="item col-md-12 col-xs-12">
            <div class="panel panel-default">
                <table class="table table-striped margin-none">
                    <thead>
                        <tr>
                            <td class="text-left"><strong>版本号</strong></td>
                            <td class="text-left"><strong>版本说明</strong></td>
                            <td class="text-right width-200"><strong>编辑</strong></td>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(v, index) in versions">
                            <td class="text-left"><strong class="text-muted ">{{index+1}}.</strong>
                              <a class="text-primary"><strong>{{v.Name}}</strong></a>
                            </td>
                            <td class="text-left">
                                </strong>
                                <a class="text-primary"><strong>{{v.Explan}}</strong></a>
                            </td>
                            <td class="text-right">
                                <button @click="recover(v)" class="btn btn-primary btn-xs">版本还原</button>
                                <button @click="deleteVersion(v)" class="btn btn-danger btn-xs">删除</button>
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
                            <h4 class="modal-title">版本设置</h4>
                        </div>
                    </div>
                    <div class="modal-body" data-scrollable>
                        <form class="form-horizontal" role="form">
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">版本名称</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="版本名称" type="text" v-model="current.Name">
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="inputEmail1" class="col-sm-3 control-label">版本说明</label>
                                <div class="col-sm-9">
                                    <input class="form-control" id="inputEmail1" placeholder="版本说明" type="text" v-model="current.Explan">
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

export default {
    name: 'version',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具',
            show: false,
            versions: {
              Name: '',
              Explan: ''
            },
            current: {
              Name: '',
              Explan: ''
            }
        }
    },
    methods: {
      close: function() {
        this.show = false
        this.current = this.temp()
      },
      save: function() {
        this.$http.post('/version', this.current).then(function(res) {
            if (res.data.error) {
                alert(res.data.Message)
            } else {
                this.loadVersion()
                this.show = false
                this.current = this.temp()
            }
        }, function() {
            alert("保存版本时服务器异常")
        })
      },
      pack: function() {
        this.current = this.temp()
        this.show = true
      },
      recover: function(version) {
        console.log(version);
        this.$http.put('/version', version).then(function(res) {
            if (res.data.error) {
                alert(res.data.msg)
            } else {
                this.loadVersion()
                alert("请重启通讯管理机")
                // this.$dispatch('toast',{
                //   level: 'error',
                //   content: '请重启通讯管理机'
                // })
            }
        }, function() {
            alert("版本重启时服务器异常")
        })
      },
      deleteVersion: function(version) {
        var url = '/version/' + version.Name
        this.$http.delete(url).then(function(res) {
            if (res.data.error) {
                alert(res.data.msg)
            } else {
                this.loadVersion()
            }
        }, function() {
            alert("删除版本时服务器异常")
        })
      },
      loadVersion: function() {
        this.versions = []
        this.$http.get('/version').then(function(res) {
            if (res.data.error) alert(res.data.msg)
            this.versions = res.data.data
        }, function() {
            alert("加载版本时服务器异常")
        })
      },
      temp: function() {
        return {
          Name: '',
          Explan: ''
        }
      }
    },
    created() {
        this.loadVersion()
    }
}

</script>
