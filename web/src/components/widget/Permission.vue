<style scoped>



</style>

<template>

<div v-if="store.state.permission" class="jumbotron bg-transparent text-center margin-none">
    <div class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>
                        <h3 class="modal-title modal-inverse">{{msg}}</h3"">
                    </div>
                    <div class="modal-body">
                      <!-- <h4 class="page-section-heading">Basic example</h4> -->
                      <div class="panel panel-default">
                        <div class="panel-body">
                          <form>
                            <div class="form-group form-control-default">
                              <label for="exampleInputEmail1">用户名</label>
                              <input type="email" class="form-control" id="exampleInputEmail1" placeholder="用户名">
                            </div>
                            <div class="form-group form-control-default">
                              <label for="exampleInputPassword1">密码</label>
                              <input type="password" class="form-control" id="exampleInputPassword1" placeholder="密码">
                            </div>
                            <!-- <div class="modal-footer"> -->
                                <button @click="close" type="button" class="btn btn-success" data-dismiss="modal">登录</button>
                                <button @click="close" type="button" class="btn btn-info" data-dismiss="modal">重启</button>
                            <!-- </div> -->

                            <!-- <button type="submit" class="btn btn-primary">Submit</button> -->
                          </form>
                        </div>
                      </div>
                        <!-- {{msg}} -->
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

</template>

<script>

export default {
    name: 'permission',
    props: {
        store: null
    },
    data() {
        return {
            msg: '欢迎登录MG2200'
        }
    },
    methods:{
      close: function (){
        this.store.hidePermission()
      }
    }
}

</script>
