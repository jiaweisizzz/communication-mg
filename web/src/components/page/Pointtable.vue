<style scoped>

/*a  upload */

.a-upload {
    padding: 6px 12px;
    height: 36px;
    line-height: 1.42857143;
    position: relative;
    cursor: pointer;
    color: #888;
    background: #fafafa;
    border: 1px solid #ddd;
    border-radius: 4px;
    overflow: hidden;
    display: inline-block;
    *display: inline;
    *zoom: 1;
}

.a-upload input {
    position: absolute;
    font-size: 14px;
    right: 0;
    top: 0;
    opacity: 0;
    filter: alpha(opacity=0);
    cursor: pointer;
}

.a-upload:hover {
    color: #444;
    background: #eee;
    border-color: #ccc;
    text-decoration: none
}

</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>点表配置</h1>
    <button @click="exportXlsxData()" class="btn btn-primary"><i class="fa fa-arrow-down"></i>导出点表</button> &nbsp
    <button class="a-upload fa fa-arrow-up">
        <input type="file" @change="importf" class="form-control">导入点表</input>
    </button> &nbsp
    <button @click="deletepoint" class="btn btn-danger">选中删除</button> &nbsp
    <button @click="save" class="btn btn-primary">保存</button> &nbsp
    <button @click="loadall" class="btn btn-primary">加载全部</button> &nbsp
    <div class="item col-md-12 col-xs-12" style="margin-top:20px">

        <!-- Tabbable Widget -->
        <div class="tabbable">

            <!-- Tabs -->
            <ul class="nav nav-tabs">
                <li :class="{active:state=='yc'}"><a href="#" @click.prevent="state='yc'"><i class="fa fa-fw fa-home"></i> 遥测</a></li>
                <li :class="{active:state=='yx'}"><a href="#" @click.prevent="state='yx'"><i class="fa fa-fw fa-user"></i> 遥信</a></li>
                <li :class="{active:state=='ym'}"><a href="#" @click.prevent="state='ym'"><i class="fa fa-fw fa-envelope"></i> 遥脉</a></li>
                <li :class="{active:state=='yk'}"><a href="#" @click.prevent="state='yk'"><i class="fa fa-fw fa-cog"></i> 遥控</a></li>
                <li :class="{active:state=='yt'}"><a href="#" @click.prevent="state='yt'"><i class="fa fa-fw fa-cog"></i> 遥调</a></li>
            </ul>
            <!-- // END Tabs -->

            <!-- Panes -->
            <div class="tab-content">
                <div class="tab-pane " :class="{active:state=='yc'}">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td style="width:100px"><strong>选择</strong></td>
                                <td style="width:100px"><strong>数据名称</strong></td>
                                <td><strong>地址[modbus]</strong></td>
                                <td><strong>点号[104/JSON]</strong></td>
                                <td><strong>功能码[modbus]</strong></td>
                                <td><strong>所属设备[勿改]</strong></td>
                                <td><strong>链路号[勿改]</strong></td>
                                <td><strong>映射点[勿改]</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in AddrYc">
                                <td class="text-primary">
                                    <div class="checkbox checkbox-danger" style="display:inline">
                                        <input :id="'tfi_' + c.ID" type="checkbox" checked v-model="c.mark">
                                        <label :for="'tfi_' + c.ID" class=""></label>
                                    </div>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Addr}}</strong></a>
                                    <td class="text-primary">{{c.AcceptID}}</td>
                                    <td class="text-primary">{{c.Funcode}}</td>
                                    <td class="text-primary">{{c.DeviceID}}</td>
                                    <td class="text-primary">{{c.ChannelID}}</td>
                                    <td class="text-primary">{{c.ID}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yx' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td style="width:100px"><strong>选择</strong></td>
                                <td style="width:100px"><strong>数据名称</strong></td>
                                <td><strong>地址[modbus]</strong></td>
                                <td><strong>点号[104/JSON]</strong></td>
                                <td><strong>功能码[modbus]</strong></td>
                                <td><strong>所属设备[勿改]</strong></td>
                                <td><strong>链路号[勿改]</strong></td>
                                <td><strong>映射点[勿改]</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in AddrYx">
                                <td class="text-primary">
                                    <div class="checkbox checkbox-danger" style="display:inline">
                                        <input :id="'tfa_' + c.ID" type="checkbox" checked v-model="c.mark">
                                        <label :for="'tfa_' + c.ID" class=""></label>
                                    </div>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Addr}}</strong></a>
                                    <td class="text-primary">{{c.AcceptID}}</td>
                                    <td class="text-primary">{{c.Funcode}}</td>
                                    <td class="text-primary">{{c.DeviceID}}</td>
                                    <td class="text-primary">{{c.ChannelID}}</td>
                                    <td class="text-primary">{{c.ID}}</td>

                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='ym' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td style="width:100px"><strong>选择</strong></td>
                                <td style="width:100px"><strong>数据名称</strong></td>
                                <td><strong>地址[modbus]</strong></td>
                                <td><strong>点号[104/JSON]</strong></td>
                                <td><strong>功能码[modbus]</strong></td>
                                <td><strong>所属设备[勿改]</strong></td>
                                <td><strong>链路号[勿改]</strong></td>
                                <td><strong>映射点[勿改]</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in AddrYm">
                                <td class="text-primary">
                                    <div class="checkbox checkbox-danger" style="display:inline">
                                        <input :id="'tfb_' + c.ID" type="checkbox" checked v-model="c.mark">
                                        <label :for="'tfb_' + c.ID" class=""></label>
                                    </div>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Addr}}</strong></a>
                                    <td class="text-primary">{{c.AcceptID}}</td>
                                    <td class="text-primary">{{c.Funcode}}</td>
                                    <td class="text-primary">{{c.DeviceID}}</td>
                                    <td class="text-primary">{{c.ChannelID}}</td>
                                    <td class="text-primary">{{c.ID}}</td>

                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yk' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td style="width:100px"><strong>选择</strong></td>
                                <td style="width:100px"><strong>数据名称</strong></td>
                                <td><strong>地址[modbus]</strong></td>
                                <td><strong>点号[104/JSON]</strong></td>
                                <td><strong>功能码[modbus]</strong></td>
                                <td><strong>所属设备[勿改]</strong></td>
                                <td><strong>链路号[勿改]</strong></td>
                                <td><strong>映射点[勿改]</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in AddrYk">
                                <td class="text-primary">
                                    <div class="checkbox checkbox-danger" style="display:inline">
                                        <input :id="'tfc_' + c.ID" type="checkbox" checked v-model="c.mark">
                                        <label :for="'tfc_' + c.ID" class=""></label>
                                    </div>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Addr}}</strong></a>
                                    <td class="text-primary">{{c.AcceptID}}</td>
                                    <td class="text-primary">{{c.Funcode}}</td>
                                    <td class="text-primary">{{c.DeviceID}}</td>
                                    <td class="text-primary">{{c.ChannelID}}</td>
                                    <td class="text-primary">{{c.ID}}</td>

                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yt' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td style="width:100px"><strong>选择</strong></td>
                                <td style="width:100px"><strong>数据名称</strong></td>
                                <td><strong>地址[modbus]</strong></td>
                                <td><strong>点号[104/JSON]</strong></td>
                                <td><strong>功能码[modbus]</strong></td>
                                <td><strong>所属设备[勿改]</strong></td>
                                <td><strong>链路号[勿改]</strong></td>
                                <td><strong>映射点[勿改]</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in AddrYt">
                                <td class="text-primary">
                                    <div class="checkbox checkbox-danger" style="display:inline">
                                        <input :id="'tfd_' + c.ID" type="checkbox" checked v-model="c.mark">
                                        <label :for="'tfd_' + c.ID" class=""></label>
                                    </div>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                <td>
                                    <a class="text-primary"><strong>{{c.Addr}}</strong></a>
                                    <td class="text-primary">{{c.AcceptID}}</td>
                                    <td class="text-primary">{{c.Funcode}}</td>
                                    <td class="text-primary">{{c.DeviceID}}</td>
                                    <td class="text-primary">{{c.ChannelID}}</td>
                                    <td class="text-primary">{{c.ID}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
            <!-- // END Panes -->
        </div>

    </div>
</div>

</template>

<script>

import xlsx from "xlsx"

export default {
    name: 'pointtable',
    data() {
        return {
            deletearray: {
                Name: '',
                AddrYc: [],
                AddrYx: [],
                AddrYm: [],
                AddrYk: [],
                AddrYt: []
            },
            cname: '',
            msg: '欢迎使用通讯管理机配置工具',
            state: 'yc',
            AddrYc: [{
                Addr: 0,
                Name: '',
                ID: 0,
                AcceptID: 0,
                Funcode: 0,
                DeviceID: 0,
                ChannelID: 0
            }],
            AddrYx: [{
                Addr: 0,
                Name: '',
                ID: 0,
                AcceptID: 0,
                Funcode: 0,
                DeviceID: 0,
                ChannelID: 0
            }],
            AddrYm: [{
                Addr: 0,
                Name: '',
                ID: 0,
                AcceptID: 0,
                Funcode: 0,
                DeviceID: 0,
                ChannelID: 0
            }],
            AddrYk: [{
                Addr: 0,
                Name: '',
                ID: 0,
                AcceptID: 0,
                Funcode: 0,
                DeviceID: 0,
                ChannelID: 0
            }],
            AddrYt: [{
                Addr: 0,
                Name: '',
                ID: 0,
                AcceptID: 0,
                Funcode: 0,
                DeviceID: 0,
                ChannelID: 0
            }],
            current: {
                Name: '',
                AddrYc: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYx: [{
                    Addr: 0,
                    Name: '',
                      ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYm: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYk: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYt: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }]
            }
        }
    },
    methods: {
        loadall: function() {
            this.loadPointtable()
        },
        deletepoint: function() {
            console.log(this.deletearray);
            this.deletearray.Name = this.cname
            for (var i = 0; i < this.AddrYc.length; i++) {
                if (this.AddrYc[i].mark == true) {
                    this.deletearray.AddrYc.push(this.AddrYc[i])
                }
            }
            for (var i = 0; i < this.AddrYx.length; i++) {
                if (this.AddrYx[i].mark == true) {
                    this.deletearray.AddrYx.push(this.AddrYx[i])
                }
            }
            for (var i = 0; i < this.AddrYm.length; i++) {
                if (this.AddrYm[i].mark == true) {
                    this.deletearray.AddrYm.push(this.AddrYm[i])
                }
            }
            for (var i = 0; i < this.AddrYk.length; i++) {
                if (this.AddrYk[i].mark == true) {
                    this.deletearray.AddrYk.push(this.AddrYk[i])
                }
            }
            for (var i = 0; i < this.AddrYt.length; i++) {
                if (this.AddrYt[i].mark == true) {
                    this.deletearray.AddrYt.push(this.AddrYt[i])
                }
            }
            console.log(this.deletearray.AddrYc);
            this.$http.put('/pointaddr', this.deletearray).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                    this.loadnewPointtable()
                    this.deletearray.Name = ''
                    this.deletearray.AddrYc = []
                    this.deletearray.AddrYx = []
                    this.deletearray.AddrYm = []
                    this.deletearray.AddrYk = []
                    this.deletearray.AddrYt = []

                }
            }, function() {
                alert("选中删除时服务器异常")
            })
        },
        save: function() {
            this.current.Name = this.cname
            console.log("save");
            console.log(this.current);
            this.$http.post('/pointaddr', this.current).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                    this.loadnewPointtable()
                    this.current = this.temp()
                }
            }, function() {
                alert("保存点时服务器异常")
            })
        },
        loadnewPointtable: function() {
            var url = '/pointaddr/' + this.cname
            this.$http.get(url).then(function(res) {
                if (res.data.error) console.log(res.data.msg)
                console.log("-----", res.data.data);
                if (res.data.data == null) {
                  this.AddrYc = null
                  this.AddrYx = null
                  this.AddrYm = null
                  this.AddrYk = null
                  this.AddrYt = null
                } else {
                  if (res.data.data.AddrYc == null) {
                    this.AddrYc = null
                  } else {
                    var addryc = res.data.data.AddrYc
                    for (var i = 0; i < addryc.length; i++) {
                        addryc[i].mark = false
                    }
                    this.AddrYc = addryc
                  }
                  if (res.data.data.AddrYx == null) {
                    this.AddrYx = null
                  } else {
                    var addryx = res.data.data.AddrYx
                    for (var a = 0; a < addryx.length; a++) {
                        addryx[a].mark = false
                    }
                    this.AddrYx = addryx
                  }
                  if (res.data.data.AddrYm == null) {
                    this.AddrYm = null
                  } else {
                    var addrym = res.data.data.AddrYm
                    for (var b = 0; b < addrym.length; b++) {
                        addrym[b].mark = false
                    }
                    this.AddrYm = addrym
                  }
                  if (res.data.data.AddrYk == null) {
                    this.AddrYk = null
                  } else {
                    var addryk = res.data.data.AddrYk
                    for (var c = 0; c < addryk.length; c++) {
                        addryk[c].mark = false
                    }
                    this.AddrYk = addryk
                  }
                  if (res.data.data.AddrYt == null) {
                    this.AddrYt = null
                  } else {
                    var addryt = res.data.data.AddrYt
                    for (var d = 0; d < addryt.length; d++) {
                        addryt[d].mark = false
                    }
                    this.AddrYt = addryt
                  }
                }
            }, function() {
                console.log("加载新的点服务器异常")
            })
        },
        importf: function(obj) { //导入
            var wb; //读取完成的数据
            // var rABS = false; //是否将文件读取为二进制字符串
            if (!obj.target.files[0]) {
                return;
            }
            var f = obj.target.files[0];
            var reader = new FileReader();
            reader.onload = this.onload

            reader.readAsBinaryString(f);
        },
        onload: function(e) {
            var wb; //读取完成的数据
            // var rABS = false; //是否将文件读取为二进制字符串
            var data = e.target.result;
            wb = xlsx.read(data, {
                type: 'binary'
            });
            var arr0 = JSON.parse(JSON.stringify(xlsx.utils.sheet_to_json(wb.Sheets[wb.SheetNames[0]])))
            var arr1 = JSON.parse(JSON.stringify(xlsx.utils.sheet_to_json(wb.Sheets[wb.SheetNames[1]])))
            var arr2 = JSON.parse(JSON.stringify(xlsx.utils.sheet_to_json(wb.Sheets[wb.SheetNames[2]])))
            var arr3 = JSON.parse(JSON.stringify(xlsx.utils.sheet_to_json(wb.Sheets[wb.SheetNames[3]])))
            var arr4 = JSON.parse(JSON.stringify(xlsx.utils.sheet_to_json(wb.Sheets[wb.SheetNames[4]])))
            this.AddrYc = []
            this.AddrYx = []
            this.AddrYm = []
            this.AddrYk = []
            this.AddrYt = []
            this.current.AddrYc = []
            this.current.AddrYx = []
            this.current.AddrYm = []
            this.current.AddrYk = []
            this.current.AddrYt = []
            if (wb.SheetNames[0] === 'AddrYc') {
                for (var i = 0; i < arr0.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr0[i]["数据名称"]
                    dt.Addr = Number(arr0[i]["地址[modbus]"])
                    dt.AcceptID = Number(arr0[i]["点号[104]"])
                    dt.Funcode = Number(arr0[i]["功能码[modbus]"])
                    dt.DeviceID = Number(arr0[i]["所属设备[勿改]"])
                    dt.ChannelID = Number(arr0[i]["链路号[勿改]"])
                    dt.ID = Number(arr0[i]["映射点[勿改]"])
                    this.AddrYc.push(dt)
                    this.current.AddrYc.push(dt)
                }
            }
            if (wb.SheetNames[1] === 'AddrYx') {
                for (var i = 0; i < arr1.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr1[i]["数据名称"]
                    dt.Addr = Number(arr1[i]["地址[modbus]"])
                    dt.AcceptID = Number(arr1[i]["点号[104]"])
                    dt.Funcode = Number(arr1[i]["功能码[modbus]"])
                    dt.DeviceID = Number(arr1[i]["所属设备[勿改]"])
                    dt.ChannelID = Number(arr1[i]["链路号[勿改]"])
                    dt.ID = Number(arr1[i]["映射点[勿改]"])
                    this.AddrYx.push(dt)
                    this.current.AddrYx.push(dt)
                }
            }
            if (wb.SheetNames[2] === 'AddrYm') {
                for (var i = 0; i < arr2.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr2[i]["数据名称"]
                    dt.Addr = Number(arr2[i]["地址[modbus]"])
                    dt.AcceptID = Number(arr2[i]["点号[104]"])
                    dt.Funcode = Number(arr2[i]["功能码[modbus]"])
                    dt.DeviceID = Number(arr2[i]["所属设备[勿改]"])
                    dt.ChannelID = Number(arr2[i]["链路号[勿改]"])
                    dt.ID = Number(arr2[i]["映射点[勿改]"])
                    this.AddrYm.push(dt)
                    this.current.AddrYm.push(dt)
                }
            }
            if (wb.SheetNames[3] === 'AddrYk') {
                for (var i = 0; i < arr3.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr3[i]["数据名称"]
                    dt.Addr = Number(arr3[i]["地址[modbus]"])
                    dt.AcceptID = Number(arr3[i]["点号[104]"])
                    dt.Funcode = Number(arr3[i]["功能码[modbus]"])
                    dt.DeviceID = Number(arr3[i]["所属设备[勿改]"])
                    dt.ChannelID = Number(arr3[i]["链路号[勿改]"])
                    dt.ID = Number(arr3[i]["映射点[勿改]"])
                    this.AddrYk.push(dt)
                    this.current.AddrYk.push(dt)
                }
            }
            if (wb.SheetNames[4] === 'AddrYt') {
                for (var i = 0; i < arr4.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr4[i]["数据名称"]
                    dt.Addr = Number(arr4[i]["地址[modbus]"])
                    dt.AcceptID = Number(arr4[i]["点号[104]"])
                    dt.Funcode = Number(arr4[i]["功能码[modbus]"])
                    dt.DeviceID = Number(arr4[i]["所属设备[勿改]"])
                    dt.ChannelID = Number(arr4[i]["链路号[勿改]"])
                    dt.ID = Number(arr4[i]["映射点[勿改]"])
                    this.AddrYt.push(dt)
                    this.current.AddrYt.push(dt)
                }
            }
        },
        exportXlsxData: function() {
            // console.log("buildXlsxData");
            let workbook = {
                SheetNames: [],
                Sheets: {}
            }
            if (this.AddrYc != null) {
                var sheet1 = 'AddrYc'
                workbook.SheetNames.push(sheet1)
                workbook.Sheets[sheet1] = {}
                let colOne = {
                    A: '数据名称',
                    B: '地址[modbus]',
                    C: '点号[104]',
                    D: '功能码[modbus]',
                    E: '所属设备[勿改]',
                    F: '链路号[勿改]',
                    G: '映射点[勿改]',
                }
                workbook.Sheets[sheet1]['!ref'] = 'A1:' + 'G' + (this.AddrYc.length + 1)
                for (var colName in colOne) {
                    workbook.Sheets[sheet1][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.AddrYc.length; i++) {
                    workbook.Sheets[sheet1]['A' + (i + 1)] = {
                        v: this.AddrYc[i - 1].Name
                    }
                    workbook.Sheets[sheet1]['B' + (i + 1)] = {
                        v: this.AddrYc[i - 1].Addr
                    }
                    workbook.Sheets[sheet1]['C' + (i + 1)] = {
                        v: this.AddrYc[i - 1].AcceptID
                    }
                    workbook.Sheets[sheet1]['D' + (i + 1)] = {
                        v: this.AddrYc[i - 1].Funcode
                    }
                    workbook.Sheets[sheet1]['E' + (i + 1)] = {
                        v: this.AddrYc[i - 1].DeviceID
                    }
                    workbook.Sheets[sheet1]['F' + (i + 1)] = {
                        v: this.AddrYc[i - 1].ChannelID
                    }
                    workbook.Sheets[sheet1]['G' + (i + 1)] = {
                        v: this.AddrYc[i - 1].ID
                    }
                }
            } else {
                var sheet1 = 'AddrYc'
                workbook.SheetNames.push(sheet1)
                workbook.Sheets[sheet1] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet1]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet1][colName + 1] = {
                        v: colOne[colName]
                    }
                }
            }
            if (this.AddrYx != null) {
                var sheet2 = 'AddrYx'
                workbook.SheetNames.push(sheet2)
                workbook.Sheets[sheet2] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet2]['!ref'] = 'A1:' + 'G' + (this.AddrYx.length + 1)
                for (var colName in colOne) {
                    workbook.Sheets[sheet2][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.AddrYx.length; i++) {
                    workbook.Sheets[sheet2]['A' + (i + 1)] = {
                        v: this.AddrYx[i - 1].Name
                    }
                    workbook.Sheets[sheet2]['B' + (i + 1)] = {
                        v: this.AddrYx[i - 1].Addr
                    }
                    workbook.Sheets[sheet2]['C' + (i + 1)] = {
                        v: this.AddrYx[i - 1].AcceptID
                    }
                    workbook.Sheets[sheet2]['D' + (i + 1)] = {
                        v: this.AddrYx[i - 1].Funcode
                    }
                    workbook.Sheets[sheet2]['E' + (i + 1)] = {
                        v: this.AddrYx[i - 1].DeviceID
                    }
                    workbook.Sheets[sheet2]['F' + (i + 1)] = {
                        v: this.AddrYx[i - 1].ChannelID
                    }
                    workbook.Sheets[sheet2]['G' + (i + 1)] = {
                        v: this.AddrYx[i - 1].ID
                    }
                }
            } else {
                var sheet2 = 'AddrYx'
                workbook.SheetNames.push(sheet2)
                workbook.Sheets[sheet2] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet2]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet2][colName + 1] = {
                        v: colOne[colName]
                    }
                }
            }
            if (this.AddrYm != null) {
                var sheet3 = 'AddrYm'
                workbook.SheetNames.push(sheet3)
                workbook.Sheets[sheet3] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet3]['!ref'] = 'A1:' + 'G' + (this.AddrYm.length + 1)
                for (var colName in colOne) {
                    workbook.Sheets[sheet3][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.AddrYm.length; i++) {
                    workbook.Sheets[sheet3]['A' + (i + 1)] = {
                        v: this.AddrYm[i - 1].Name
                    }
                    workbook.Sheets[sheet3]['B' + (i + 1)] = {
                        v: this.AddrYm[i - 1].Addr
                    }
                    workbook.Sheets[sheet3]['C' + (i + 1)] = {
                        v: this.AddrYm[i - 1].AcceptID
                    }
                    workbook.Sheets[sheet3]['D' + (i + 1)] = {
                        v: this.AddrYm[i - 1].Funcode
                    }
                    workbook.Sheets[sheet3]['E' + (i + 1)] = {
                        v: this.AddrYm[i - 1].DeviceID
                    }
                    workbook.Sheets[sheet3]['F' + (i + 1)] = {
                        v: this.AddrYm[i - 1].ChannelID
                    }
                    workbook.Sheets[sheet3]['G' + (i + 1)] = {
                        v: this.AddrYm[i - 1].ID
                    }
                }
            } else {
                var sheet3 = 'AddrYm'
                workbook.SheetNames.push(sheet3)
                workbook.Sheets[sheet3] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet3]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet3][colName + 1] = {
                        v: colOne[colName]
                    }
                }
            }
            if (this.AddrYk != null) {
                var sheet4 = 'AddrYk'
                workbook.SheetNames.push(sheet4)
                workbook.Sheets[sheet4] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet4]['!ref'] = 'A1:' + 'G' + (this.AddrYk.length + 1)
                for (var colName in colOne) {
                    workbook.Sheets[sheet4][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.AddrYk.length; i++) {
                    workbook.Sheets[sheet4]['A' + (i + 1)] = {
                        v: this.AddrYk[i - 1].Name
                    }
                    workbook.Sheets[sheet4]['B' + (i + 1)] = {
                        v: this.AddrYk[i - 1].Addr
                    }
                    workbook.Sheets[sheet4]['C' + (i + 1)] = {
                        v: this.AddrYk[i - 1].AcceptID
                    }
                    workbook.Sheets[sheet4]['D' + (i + 1)] = {
                        v: this.AddrYk[i - 1].Funcode
                    }
                    workbook.Sheets[sheet4]['E' + (i + 1)] = {
                        v: this.AddrYk[i - 1].DeviceID
                    }
                    workbook.Sheets[sheet4]['F' + (i + 1)] = {
                        v: this.AddrYk[i - 1].ChannelID
                    }
                    workbook.Sheets[sheet4]['G' + (i + 1)] = {
                        v: this.AddrYk[i - 1].ID
                    }
                }
            } else {
                var sheet4 = 'AddrYk'
                workbook.SheetNames.push(sheet4)
                workbook.Sheets[sheet4] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet4]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet4][colName + 1] = {
                        v: colOne[colName]
                    }
                }
            }
            if (this.AddrYt != null) {
                var sheet5 = 'AddrYt'
                workbook.SheetNames.push(sheet5)
                workbook.Sheets[sheet5] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet5]['!ref'] = 'A1:' + 'G' + (this.AddrYt.length + 1)
                for (var colName in colOne) {
                    workbook.Sheets[sheet5][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.AddrYt.length; i++) {
                    workbook.Sheets[sheet5]['A' + (i + 1)] = {
                        v: this.AddrYt[i - 1].Name
                    }
                    workbook.Sheets[sheet5]['B' + (i + 1)] = {
                        v: this.AddrYt[i - 1].Addr
                    }
                    workbook.Sheets[sheet5]['C' + (i + 1)] = {
                        v: this.AddrYt[i - 1].AcceptID
                    }
                    workbook.Sheets[sheet5]['D' + (i + 1)] = {
                        v: this.AddrYt[i - 1].Funcode
                    }
                    workbook.Sheets[sheet5]['E' + (i + 1)] = {
                        v: this.AddrYt[i - 1].DeviceID
                    }
                    workbook.Sheets[sheet5]['F' + (i + 1)] = {
                        v: this.AddrYt[i - 1].ChannelID
                    }
                    workbook.Sheets[sheet5]['G' + (i + 1)] = {
                        v: this.AddrYt[i - 1].ID
                    }                }
            } else {
                var sheet5 = 'AddrYt'
                workbook.SheetNames.push(sheet5)
                workbook.Sheets[sheet5] = {}
                let colOne = {
                  A: '数据名称',
                  B: '地址[modbus]',
                  C: '点号[104]',
                  D: '功能码[modbus]',
                  E: '所属设备[勿改]',
                  F: '链路号[勿改]',
                  G: '映射点[勿改]',
                }
                workbook.Sheets[sheet5]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet5][colName + 1] = {
                        v: colOne[colName]
                    }
                }
            }
            this.export(workbook)
        },
        export: function(workBook) {
            // 导出 Excel
            var wopts = {
                bookType: 'xlsx',
                bookSST: false,
                type: 'binary'
            };
            var wbout = xlsx.write(workBook, wopts);

            function s2ab(s) {
                var buf = new ArrayBuffer(s.length);
                var view = new Uint8Array(buf);
                for (var i = 0; i != s.length; ++i) view[i] = s.charCodeAt(i) & 0xFF;
                console.log("s2ab");
                return buf;
            }
            /* the saveAs call downloads a file on the local machine */
            let name = this.cname
            this.saveAs(new Blob([s2ab(wbout)], {
                type: ""
            }), name + ".xlsx")
        },
        saveAs: function(blob, fileName) {
            var tmpa = document.createElement("a");
            tmpa.download = fileName || "下载";
            tmpa.href = URL.createObjectURL(blob);
            tmpa.click();
            setTimeout(function() {
                URL.revokeObjectURL(blob);
            }, 100);
        },
        loadPointtable: function() {
            this.$http.get('/pointaddr').then(function(res) {
                if (res.data.error) console.log(res.data.msg)
                if (res.data.data.AddrYc == null) {
                  this.current.AddrYc = []
                  this.AddrYc = null
                } else {
                  var addryc = res.data.data.AddrYc
                  for (var i = 0; i < addryc.length; i++) {
                      addryc[i].mark = false
                  }
                  this.current.AddrYc = addryc
                  this.AddrYc = addryc
                }
                if (res.data.data.AddrYx == null) {
                  this.current.AddrYx = []
                  this.AddrYx = null
                } else {
                  var addryx = res.data.data.AddrYx
                  for (var a = 0; a < addryx.length; a++) {
                      addryx[a].mark = false
                  }
                  this.current.AddrYx = addryx
                  this.AddrYx = addryx
                }
                if (res.data.data.AddrYm == null) {
                  this.current.AddrYm = []
                  this.AddrYm = null
                } else {
                  var addrym = res.data.data.AddrYm
                  for (var b = 0; b < addrym.length; b++) {
                      addrym[b].mark = false
                  }
                  this.current.AddrYm = addrym
                  this.AddrYm = addrym
                }
                if (res.data.data.AddrYk == null) {
                  this.current.AddrYk = []
                  this.AddrYk = null
                } else {
                  var addryk = res.data.data.AddrYk
                  for (var c = 0; c < addryk.length; c++) {
                      addryk[c].mark = false
                  }
                  this.current.AddrYk = addryk
                  this.AddrYk = addryk
                }
                if (res.data.data.AddrYt == null) {
                  this.current.AddrYt = []
                  this.AddrYt = null
                } else {
                  var addryt = res.data.data.AddrYt
                  for (var d = 0; d < addryt.length; d++) {
                      addryt[d].mark = false
                  }
                  this.current.AddrYt = addryt
                  this.AddrYt = addryt
                }
            }, function() {
                console.log("加载点时服务器异常")
            })
        },
        temp: function() {
            return {
                Name: '',
                AddrYc: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYx: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYm: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYk: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }],
                AddrYt: [{
                    Addr: 0,
                    Name: '',
                    ID: 0,
                    AcceptID: 0,
                    Funcode: 0,
                    DeviceID: 0,
                    ChannelID: 0
                }]
            }
        }
    },
    created() {
        this.cname = this.$route.params.cname
        this.loadnewPointtable()
    }
}

</script>
