<style scoped>



</style>

<template>

<div v-if="store.state.restart" class="jumbotron bg-transparent text-center margin-none">
    <div class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content">
                    <div class="modal-header">
                        <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>
                        <h4 class="modal-title">系统重启</h4>
                    </div>
                    <div class="modal-body">
                        {{msg}}
                    </div>
                    <div class="modal-footer">
                        <button @click="close" @click.prevent.stop="store.hideRestart" type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                        <button @click="restart" type="button" class="btn btn-primary" data-dismiss="modal">确认</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

</template>

<script>

export default {
    name: 'restart',
    props: {
        store: null
    },
    data() {
        return {
            msg: '您确认系统重启吗？'
        }
    },
    methods:{
      restart: function() {
        this.$http.get('/restart').then(function(res) {
            if (res.data.error) alert(res.data.msg)
            this.store.hideRestart()
        }, function() {
            alert("系统重启服务器异常")
        })
      },
      close: function (){
        this.store.hideRestart()
      }
    }
}

</script>
