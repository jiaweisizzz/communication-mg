<style scoped>



</style>

<template>

<div v-if="store.state.message" class="jumbotron bg-transparent　margin-none">
    <div class="modal fade slide-down in" style="display: block;">
        <div class="modal-dialog">
            <div class="v-cell">
                <div class="modal-content modal-lg">
                    <div class="modal-header">
                        <button @click="close" type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span><span class="sr-only">Close</span></button>
                        <div class="text-center">
                            <h4 class="modal-title">报文监视</h4>
                            <select v-model="selected">
                                <option v-for="(c, index) in channels" v-bind:value="c.ID">
                                    {{ c.Name}}
                                </option>
                            </select>
                        </div>

                    </div>
                    <div class="modal-body" data-scrollable style="height:400px;">
                        <div v-for="(c, index) in msgs" class="expandable-content">
                            <p> {{index}} 时间:{{nowtime()}} 报文:{{c}}</p>
                            <hr>
                        </div>
                        <div ref="end" style="padding:20px"></div>
                    </div>
                    <div class="modal-footer text-center">
                        <div class="pull-left">
                            <!-- <div class="checkbox checkbox-primary checkbox-inline">
                                <input id="checkbox１" checked="" type="checkbox">
                                <label for="checkbox１">报文解析</label>
                            </div>

                            <div class="checkbox checkbox-primary checkbox-inline">
                                <input id="checkbox２" checked="" type="checkbox">
                                <label for="checkbox２">打印时间</label>
                            </div> -->

                            <!-- <div class="checkbox checkbox-primary checkbox-inline">
                              <input id="checkbox３" checked="" type="checkbox">
                              <label for="checkbox３">开始截取</label>
                          </div>

                          <div class="checkbox checkbox-primary checkbox-inline">
                              <input id="checkbox４" checked="" type="checkbox">
                              <label for="checkbox４">停止截取</label>
                          </div> -->

                            <div class="checkbox checkbox-primary checkbox-inline">
                                <input id="checkbox5" value=true v-model="checkedreceive" type="checkbox">
                                <label for="checkbox5">接收报文</label>
                            </div>

                            <div class="checkbox checkbox-primary checkbox-inline">
                                <input id="checkbox６" value=true v-model="checkedsend" type="checkbox">
                                <label for="checkbox６">发送报文</label>
                            </div>
                        </div>
                        <div class="pull-right">
                            <!-- <button class="btn btn-danger btn-xs"><i class="fa fa-stop"></i></button> -->
                            <button v-if="state == 'run'" @click="stopWatch('stop')" class="btn btn-primary btn-stroke btn-xs" :class="{run}"><i class="fa fa-export"></i>停止监视</button>
                            <button v-if="state == 'stop'" @click="stopWatch('run')" class="btn btn-primary btn-stroke btn-xs" :class="{run}"><i class="fa fa-export"></i>继续监视</button>
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
    name: 'message',
    props: {
        store: null
    },
    data() {
        return {
            checkedsend: true,
            checkedreceive: true,
            state: 'run',
            ws: null,
            msgs: [],
            channels: [{
                ID: 0,
                Name: ''
            }],
            current: {
                ChannelId: 0,
                Temp: null
            },
            selected: null
        }
    },
    watch: {
        selected: function(val) {
            this.current.ChannelId = val
            this.websocketsend()
        }
    },
    methods: {
        insert_flg: function(str, flg, sn) {
            var newstr = "";
            for (var i = 1; i < str.length; i += sn) {
                var tmp = str.substring(i, i + sn);
                newstr += tmp + flg;
            }
            return newstr;
        },
        stopWatch: function(state) {
            this.state = state
        },
        close: function() {
            this.store.hideMessage()
            this.selected = null
            this.msgs = []
                // this.ws.close()
            window.eventHub.$emit('app', JSON.stringify({
                ChannelId: this.current.ChannelId * (-1) + (-1),
                Temp: null
            }))
        },
        websocketmessages: function(e) {
          var star = e.charAt(0)
            e = this.insert_flg(e, " ", 2)
            e = star + e
            if (e != "" && this.state == "run") {
                if (this.checkedsend == true && this.checkedreceive == true) {
                    this.msgs.push(e)
                    this.$refs.end.scrollIntoView()
                } else if (this.checkedsend == true && this.checkedreceive == false) {
                    if (e.charAt(0) == "发") {
                        this.msgs.push(e)
                        this.$refs.end.scrollIntoView()
                    }
                } else if (this.checkedsend == false && this.checkedreceive == true) {
                    if (e.charAt(0) == "收") {
                        this.msgs.push(e)
                        this.$refs.end.scrollIntoView()
                    }
                }
                // this.msgs.push(e)
                // this.$refs.end.scrollIntoView()
            }
        },
        nowtime: function() {
            var myDate = new Date()
            var mytime = myDate.toLocaleTimeString()
            return mytime
        },
        websocketsend: function(e) {
            window.eventHub.$emit('app', JSON.stringify(this.current))
        },
        loadChannels: function() { //更新通道
            this.$http.get('/channels').then(function(res) {
                if (res.data.error) alert(res.data.msg)
                this.channels = res.data.data
            }, function() {
                alert("加载报文通道服务器异常")
            })
        }
    },
    created() {
        this.loadChannels()
            // this.ws = new WebSocket("ws://127.0.0.1:8080/ws")
            // this.ws.onmessage = this.websocketmessages
        window.eventHub.$on('onmessage', this.websocketmessages)
    }
}

</script>
