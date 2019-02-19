<style scoped>



</style>

<template>

<div v-if="show" class="jumbotron bg-transparent text-center margin-none">
    <div class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3 class="modal-title modal-inverse">{{msg}}</h3"">
                    </div>
                    <div class="modal-body">
                      <div class="panel panel-default">
                        <div class="panel-body">
                          <form>
                            <div class="form-group form-control-default">
                              <label for="exampleInputEmail1">用户名</label>
                              <input class="form-control" id="exampleInputEmail1" placeholder="用户名" type="text" v-model="current.Name">
                            </div>
                            <div class="form-group form-control-default">
                              <label for="exampleInputPassword1">密码</label>
                              <input type="password" class="form-control" id="exampleInputPassword1" placeholder="密码" v-model="current.Password">
                            </div>
                                <button @click="login" type="button" class="btn btn-success" data-dismiss="modal">登录</button>
                                <button @click="registe" type="button" class="btn btn-info" data-dismiss="modal">注册</button>
                          </form>
                        </div>
                      </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
<div v-else class="jumbotron bg-transparent text-center margin-none">
    <div class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <!-- <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button> -->
                        <h3 class="modal-title modal-inverse">欢迎注册MG2200</h3"">
                    </div>
                    <div class="modal-body">
                      <div class="panel panel-default">
                        <div class="panel-body">
                          <form>
                            <div class="form-group form-control-default">
                              <label for="exampleInputEmail1">用户名</label>
                              <input class="form-control" id="exampleInputEmail1" placeholder="用户名" type="text" v-model="register.Name">
                            </div>
                            <div class="form-group form-control-default">
                              <label for="exampleInputPassword1">密码</label>
                              <input type="password" class="form-control" id="exampleInputPassword1" placeholder="密码" v-model="register.Password">
                            </div>
                            <div class="form-group form-control-default">
                              <label for="exampleInputEmail1">确认密码</label>
                              <input type="password" class="form-control" id="exampleInputEmail1" placeholder="确认密码" v-model="register.Samepassword">
                            </div>
                                <button @click="back" type="button" class="btn btn-success" data-dismiss="modal">返回</button>
                                <button @click="contorlregister" type="button" class="btn btn-info" data-dismiss="modal">注册</button>

                          </form>
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
    name: 'login',
    props: {
        store: null
    },
    data() {
        return {
            msg: '欢迎登录MG2200',
            show: true,
            showregister: false,
            state: false,
            current: {
                Name: '',
                Password: ''
            },
            register: {
                Name: '',
                Password: '',
                Samepassword: ''
            }
        }
    },
    methods: {
        back: function() {
            this.showregister = false
            this.show = true
        },
        registe: function() {
            this.show = false
            this.showregister = true
        },
        contorlregister: function() {
          if ((this.register.Name == "") || (this.register.Password == "") || (this.register.Samepassword == "")) {
            alert("用户名或密码格式错误")
          }
          this.showregister = false
          this.show = true
        },
        login: function() {
            var url = '/users/' + this.current.Name + '/' + this.current.Password
            this.$http.get(url).then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.state = res.data.data
                if (this.state == true) {
                    this.show = false
                    this.$router.push('/home')
                }
            }, function() {
                alert("login服务器异常")
            })
        }
    }
}

</script>
