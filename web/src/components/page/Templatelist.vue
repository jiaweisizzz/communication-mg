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
    <h1>模板列表</h1>
    <button @click="buildXlsxData()" class="btn btn-primary"><i class="fa fa-arrow-down"></i>导出模板</button> &nbsp
    <!-- <button @click="importf(obj)" class="btn btn-danger"><i class="fa fa-arrow-up"></i>导入模板</button> -->
    <button class="a-upload fa fa-arrow-up">
        <input type="file" @change="importf" class="form-control">导入模板</input>
    </button> &nbsp
    <button @click="save" class="btn btn-primary">保存</button> &nbsp
    <button @click="reset" class="btn btn-primary">重置</button>
    <div class="item col-md-12 col-xs-12" style="margin-top:20px">
        <!-- <div class="panel panel-default">
            <div class="panel-body">
                <div class="row">
                    <div class="col-xs-2 text-center">
                        <p class="text-body-2  text-primary margin-none">名称 ：</p>
                    </div>
                    <div class="col-xs-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.tname}}</p>
                    </div>
                    <div class="col-xs-2 text-center">
                        <p class="text-body-2  text-primary margin-none">终端号 :</p>
                    </div>
                    <div class="col-xs-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.tid}}</p>
                    </div>
                    <div class="col-xs-2 text-center">
                        <p class="text-body-2  text-primary margin-none">链路 :</p>
                    </div>
                    <div class="col-xs-2 text-center">
                        <p class="text-title  text-success margin-none">{{this.tchannel}}</p>
                    </div>
                </div>
            </div>
        </div> -->


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
                                <td class="text-left "><strong>数据名称</strong></td>
                                <td><strong>点号</strong></td>
                                <td><strong>系数</strong></td>
                                <td><strong>分区</strong></td>
                                <td><strong>功能码</strong></td>
                                <td><strong>寄存器地址</strong></td>
                                <td><strong>比特位</strong></td>
                                <td><strong>数据类型</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in ycs">
                                <td class="text-left "><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                </strong>
                                <td>
                                    <a class="text-primary"><strong>{{c.ID}}</strong></a>
                                    <td class="text-primary">{{c.Scale}}</td>
                                    <td class="text-primary">{{c.Data.Partition}}</td>
                                    <td class="text-primary">{{c.Data.Funcode}}</td>
                                    <td class="text-primary">{{c.Data.Addr}}</td>
                                    <td class="text-primary">{{c.Data.Bit}}</td>
                                    <td class="text-success">{{showLabel(c.Data.Type)}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yx' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left "><strong>数据名称</strong></td>
                                <td><strong>点号</strong></td>
                                <td><strong>系数</strong></td>
                                <td><strong>分区</strong></td>
                                <td><strong>功能码</strong></td>
                                <td><strong>寄存器地址</strong></td>
                                <td><strong>比特位</strong></td>
                                <td><strong>数据类型</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in yxs">
                                <td class="text-left "><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                </strong>
                                <td>
                                    <a class="text-primary"><strong>{{c.ID}}</strong></a>
                                    <td class="text-primary">{{c.Scale}}</td>
                                    <td class="text-primary">{{c.Data.Partition}}</td>
                                    <td class="text-primary">{{c.Data.Funcode}}</td>
                                    <td class="text-primary">{{c.Data.Addr}}</td>
                                    <td class="text-primary">{{c.Data.Bit}}</td>
                                    <td class="text-success">{{showLabel(c.Data.Type)}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='ym' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left "><strong>数据名称</strong></td>
                                <td><strong>点号</strong></td>
                                <td><strong>系数</strong></td>
                                <td><strong>分区</strong></td>
                                <td><strong>功能码</strong></td>
                                <td><strong>寄存器地址</strong></td>
                                <td><strong>比特位</strong></td>
                                <td><strong>数据类型</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in yms">
                                <td class="text-left "><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                </strong>
                                <td>
                                    <a class="text-primary"><strong>{{c.ID}}</strong></a>
                                    <td class="text-primary">{{c.Scale}}</td>
                                    <td class="text-primary">{{c.Data.Partition}}</td>
                                    <td class="text-primary">{{c.Data.Funcode}}</td>
                                    <td class="text-primary">{{c.Data.Addr}}</td>
                                    <td class="text-primary">{{c.Data.Bit}}</td>
                                    <td class="text-success">{{showLabel(c.Data.Type)}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yk' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left "><strong>数据名称</strong></td>
                                <td><strong>点号</strong></td>
                                <td><strong>系数</strong></td>
                                <td><strong>功能码</strong></td>
                                <td><strong>寄存器地址</strong></td>
                                <td><strong>比特位</strong></td>
                                <td><strong>数据类型</strong></td>
                                <td><strong>合闸</strong></td>
                                <td><strong>分闸</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in yks">
                                <td class="text-left "><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                </strong>
                                <td>
                                    <a class="text-primary"><strong>{{c.ID}}</strong></a>
                                    <td class="text-primary">{{c.Scale}}</td>
                                    <td class="text-primary">{{c.Data.Funcode}}</td>
                                    <td class="text-primary">{{c.Data.Addr}}</td>
                                    <td class="text-primary">{{c.Data.Bit}}</td>
                                    <td class="text-success">{{showLabel(c.Data.Type)}}</td>
                                    <td class="text-primary">{{c.Data.On}}</td>
                                    <td class="text-primary">{{c.Data.Off}}</td>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="tab-pane " :class="{active:state=='yt' } ">
                    <table class="table table-striped margin-none">
                        <thead>
                            <tr>
                                <td class="text-left "><strong>数据名称</strong></td>
                                <td><strong>点号</strong></td>
                                <td><strong>系数</strong></td>
                                <td><strong>功能码</strong></td>
                                <td><strong>寄存器地址</strong></td>
                                <td><strong>比特位</strong></td>
                                <td><strong>数据类型</strong></td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(c, index) in yts">
                                <td class="text-left "><strong class="text-muted ">{{index+1}}.</strong>
                                    <a class="text-primary"><strong>{{c.Name}}</strong></a>
                                </td>
                                </strong>
                                <td>
                                    <a class="text-primary"><strong>{{c.ID}}</strong></a>
                                    <td class="text-primary">{{c.Scale}}</td>
                                    <td class="text-primary">{{c.Data.Funcode}}</td>
                                    <td class="text-primary">{{c.Data.Addr}}</td>
                                    <td class="text-primary">{{c.Data.Bit}}</td>
                                    <td class="text-success">{{showLabel(c.Data.Type)}}</td>
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

<!-- xlsx -->
<!-- <div id="xlsxreader">
  <p>asdasd</p>
</div> -->

</template>

<script>

import xlsx from "xlsx"

export default {
    // ready() {
    //         this.$els.file.addEventListener('change', this.importf)
    //     },
    name: 'templatelist',
    data() {
        return {
            // xlsxInstance: null, // xlsx实例
            msg: '欢迎使用通讯管理机配置工具',
            state: 'yc',
            tname: '',
            tid: 0,
            tchannel: '',
            tpotocol: '',
            show: false,
            ycs: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: {
                    Partition: '',
                    Funcode: '',
                    Addr: '',
                    Bit: '',
                    Type: ''
                }
            }],
            yxs: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: {
                    Partition: '',
                    Funcode: '',
                    Addr: '',
                    Bit: '',
                    Type: ''
                }
            }],
            yms: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: {
                    Partition: '',
                    Funcode: '',
                    Addr: '',
                    Bit: '',
                    Type: ''
                }
            }],
            yks: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: {
                    Partition: '',
                    Funcode: '',
                    Addr: '',
                    Bit: '',
                    Type: '',
                    On: '',
                    Off: ''
                }
            }],
            yts: [{
                ID: 0,
                Name: '',
                Scale: 0.0,
                Data: {
                    Partition: '',
                    Funcode: '',
                    Addr: '',
                    Bit: '',
                    Type: ''
                }
            }],
            current: {
                ID: 0,
                Name: '',
                Channel: '',
                Potocol: '',
                ycs: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yxs: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yms: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yks: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: '',
                        On: '',
                        Off: ''
                    }
                }],
                yts: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }]
            }
        }
    },
    methods: {
        reset: function() {
            this.loadData()
            this.current = this.temp()
        },
        save: function() {
            console.log("save");
            console.log(this.current);
            this.$http.put('/templates', this.current).then(function(res) {
                if (res.data.error) {
                    alert(res.data.msg)
                } else {
                  window.eventHub.$emit("toast", {
                      level: "success",
                      message: res.data.message,
                      title: "title",
                  })
                    this.loadData()
                    this.current = this.temp()
                }
            }, function() {
                alert("保存模板服务器异常")
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
            this.ycs = []
            this.yxs = []
            this.yms = []
            this.yks = []
            this.yts = []
            this.current.ycs = []
            this.current.yxs = []
            this.current.yms = []
            this.current.yks = []
            this.current.yts = []
            if (wb.SheetNames[0] === 'YC') {
                for (var i = 0; i < arr0.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr0[i]["数据名称"]
                    dt.ID = Number(arr0[i]["点号"])
                    dt.Scale = Number(arr0[i]["系数"])
                    dt.Data.Partition = arr0[i]["分区"]
                    dt.Data.Funcode = arr0[i]["功能码"]
                    dt.Data.Addr = arr0[i]["寄存器地址"]
                    dt.Data.Bit = arr0[i]["比特位"]
                    dt.Data.Type = arr0[i]["数据类型"]
                    this.ycs.push(dt)
                    this.current.ycs.push(dt)
                }
            }
            if (wb.SheetNames[1] === 'YX') {
                for (var i = 0; i < arr1.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr1[i]["数据名称"]
                    dt.ID = Number(arr1[i]["点号"])
                    dt.Scale = Number(arr1[i]["系数"])
                    dt.Data.Partition = arr1[i]["分区"]
                    dt.Data.Funcode = arr1[i]["功能码"]
                    dt.Data.Addr = arr1[i]["寄存器地址"]
                    dt.Data.Bit = arr1[i]["比特位"]
                    dt.Data.Type = arr1[i]["数据类型"]
                    this.yxs.push(dt)
                    this.current.yxs.push(dt)
                }
            }
            if (wb.SheetNames[2] === 'YM') {
                for (var i = 0; i < arr2.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr2[i]["数据名称"]
                    dt.ID = Number(arr2[i]["点号"])
                    dt.Scale = Number(arr2[i]["系数"])
                    dt.Data.Partition = arr2[i]["分区"]
                    dt.Data.Funcode = arr2[i]["功能码"]
                    dt.Data.Addr = arr2[i]["寄存器地址"]
                    dt.Data.Bit = arr2[i]["比特位"]
                    dt.Data.Type = arr2[i]["数据类型"]
                    this.yms.push(dt)
                    this.current.yms.push(dt)
                }
            }
            if (wb.SheetNames[3] === 'YK') {
                for (var i = 0; i < arr3.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr3[i]["数据名称"]
                    dt.ID = Number(arr3[i]["点号"])
                    dt.Scale = Number(arr3[i]["系数"])
                    dt.Data.Funcode = arr3[i]["功能码"]
                    dt.Data.Addr = arr3[i]["寄存器地址"]
                    dt.Data.Bit = arr3[i]["比特位"]
                    dt.Data.Type = arr3[i]["数据类型"]
                    dt.Data.On = arr3[i]["合闸"]
                    dt.Data.Off = arr3[i]["分闸"]
                    this.yks.push(dt)
                    this.current.yks.push(dt)
                }
            }
            if (wb.SheetNames[4] === 'YT') {
                for (var i = 0; i < arr4.length; i++) {
                    var dt = {
                        Data: {}
                    }
                    dt.Name = arr4[i]["数据名称"]
                    dt.ID = Number(arr4[i]["点号"])
                    dt.Scale = Number(arr4[i]["系数"])
                    dt.Data.Funcode = arr4[i]["功能码"]
                    dt.Data.Addr = arr4[i]["寄存器地址"]
                    dt.Data.Bit = arr4[i]["比特位"]
                    dt.Data.Type = arr4[i]["数据类型"]
                    this.yts.push(dt)
                    this.current.yts.push(dt)
                }
            }
        },
        buildXlsxData: function() {
            // console.log("buildXlsxData");
            let workbook = {
                SheetNames: [],
                Sheets: {}
            }
            if (this.ycs != null) {
                var sheet1 = 'YC'
                workbook.SheetNames.push(sheet1)
                workbook.Sheets[sheet1] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    // I: '数据类型对照表',
                }
                workbook.Sheets[sheet1]['!ref'] = 'A1:' + 'H' + this.ycs.length + 1
                // workbook.Sheets[sheet1]['!ref'] = workbook.Sheets[sheet1]['!ref'] + 'I'+21
                for (var colName in colOne) {
                    workbook.Sheets[sheet1][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.ycs.length; i++) {
                    workbook.Sheets[sheet1]['A' + (i + 1)] = {
                        v: this.ycs[i - 1].Name
                    }
                    workbook.Sheets[sheet1]['B' + (i + 1)] = {
                        v: this.ycs[i - 1].ID
                    }
                    workbook.Sheets[sheet1]['C' + (i + 1)] = {
                        v: this.ycs[i - 1].Scale
                    }
                    workbook.Sheets[sheet1]['D' + (i + 1)] = {
                        v: this.ycs[i - 1].Data.Partition
                    }
                    workbook.Sheets[sheet1]['E' + (i + 1)] = {
                        v: this.ycs[i - 1].Data.Funcode
                    }
                    workbook.Sheets[sheet1]['F' + (i + 1)] = {
                        v: this.ycs[i - 1].Data.Addr
                    }
                    workbook.Sheets[sheet1]['G' + (i + 1)] = {
                        v: this.ycs[i - 1].Data.Bit
                    }
                    workbook.Sheets[sheet1]['H' + (i + 1)] = {
                        v: this.ycs[i - 1].Data.Type
                    }
                }
                // workbook.Sheets[sheet1]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet1]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet1]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet1]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet1]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet1]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet1]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet1]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet1]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet1]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet1]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet1]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet1]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet1]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet1]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet1]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet1]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet1]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet1]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet1]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            } else {
                var sheet1 = 'YC'
                workbook.SheetNames.push(sheet1)
                workbook.Sheets[sheet1] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    I: '数据类型对照表',
                }
                workbook.Sheets[sheet1]['!ref'] = 'A1:' + 'H' + 1
                // workbook.Sheets[sheet1]['!ref'] = workbook.Sheets[sheet1]['!ref'] + 'I' + 21
                for (var colName in colOne) {
                    workbook.Sheets[sheet1][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                // workbook.Sheets[sheet1]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet1]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet1]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet1]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet1]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet1]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet1]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet1]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet1]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet1]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet1]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet1]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet1]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet1]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet1]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet1]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet1]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet1]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet1]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet1]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            }
            if (this.yxs != null) {
                var sheet2 = 'YX'
                workbook.SheetNames.push(sheet2)
                workbook.Sheets[sheet2] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    // I: '数据类型对照表',
                }
                workbook.Sheets[sheet2]['!ref'] = 'A1:' + 'H' + this.yxs.length + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet2][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.yxs.length; i++) {
                    workbook.Sheets[sheet2]['A' + (i + 1)] = {
                        v: this.yxs[i - 1].Name
                    }
                    workbook.Sheets[sheet2]['B' + (i + 1)] = {
                        v: this.yxs[i - 1].ID
                    }
                    workbook.Sheets[sheet2]['C' + (i + 1)] = {
                        v: this.yxs[i - 1].Scale
                    }
                    workbook.Sheets[sheet2]['D' + (i + 1)] = {
                        v: this.yxs[i - 1].Data.Partition
                    }
                    workbook.Sheets[sheet2]['E' + (i + 1)] = {
                        v: this.yxs[i - 1].Data.Funcode
                    }
                    workbook.Sheets[sheet2]['F' + (i + 1)] = {
                        v: this.yxs[i - 1].Data.Addr
                    }
                    workbook.Sheets[sheet2]['G' + (i + 1)] = {
                        v: this.yxs[i - 1].Data.Bit
                    }
                    workbook.Sheets[sheet2]['H' + (i + 1)] = {
                        v: this.yxs[i - 1].Data.Type
                    }
                }
                // workbook.Sheets[sheet2]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet2]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet2]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet2]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet2]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet2]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet2]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet2]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet2]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet2]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet2]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet2]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet2]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet2]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet2]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet2]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet2]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet2]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet2]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet2]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            } else {
                var sheet2 = 'YX'
                workbook.SheetNames.push(sheet2)
                workbook.Sheets[sheet2] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    // I: '数据类型对照表',
                }
                workbook.Sheets[sheet2]['!ref'] = 'A1:' + 'H' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet2][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                // workbook.Sheets[sheet2]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet2]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet2]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet2]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet2]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet2]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet2]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet2]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet2]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet2]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet2]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet2]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet2]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet2]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet2]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet2]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet2]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet2]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet2]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet2]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            }
            if (this.yms != null) {
                var sheet3 = 'YM'
                workbook.SheetNames.push(sheet3)
                workbook.Sheets[sheet3] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    // I: '数据类型对照表',
                }
                workbook.Sheets[sheet3]['!ref'] = 'A1:' + 'H' + this.yms.length + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet3][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.yms.length; i++) {
                    workbook.Sheets[sheet3]['A' + (i + 1)] = {
                        v: this.yms[i - 1].Name
                    }
                    workbook.Sheets[sheet3]['B' + (i + 1)] = {
                        v: this.yms[i - 1].ID
                    }
                    workbook.Sheets[sheet3]['C' + (i + 1)] = {
                        v: this.yms[i - 1].Scale
                    }
                    workbook.Sheets[sheet3]['D' + (i + 1)] = {
                        v: this.yms[i - 1].Data.Partition
                    }
                    workbook.Sheets[sheet3]['E' + (i + 1)] = {
                        v: this.yms[i - 1].Data.Funcode
                    }
                    workbook.Sheets[sheet3]['F' + (i + 1)] = {
                        v: this.yms[i - 1].Data.Addr
                    }
                    workbook.Sheets[sheet3]['G' + (i + 1)] = {
                        v: this.yms[i - 1].Data.Bit
                    }
                    workbook.Sheets[sheet3]['H' + (i + 1)] = {
                        v: this.yms[i - 1].Data.Type
                    }
                }
                // workbook.Sheets[sheet3]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet3]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet3]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet3]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet3]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet3]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet3]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet3]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet3]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet3]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet3]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet3]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet3]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet3]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet3]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet3]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet3]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet3]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet3]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet3]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            } else {
                var sheet3 = 'YM'
                workbook.SheetNames.push(sheet3)
                workbook.Sheets[sheet3] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '分区',
                    E: '功能码',
                    F: '寄存器地址',
                    G: '比特位',
                    H: '数据类型',
                    // I: '数据类型对照表',
                }
                workbook.Sheets[sheet3]['!ref'] = 'A1:' + 'H' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet3][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                // workbook.Sheets[sheet3]['I' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet3]['I' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet3]['I' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet3]['I' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet3]['I' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet3]['I' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet3]['I' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet3]['I' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet3]['I' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet3]['I' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet3]['I' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet3]['I' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet3]['I' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet3]['I' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet3]['I' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet3]['I' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet3]['I' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet3]['I' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet3]['I' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet3]['I' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            }
            if (this.yks != null) {
                var sheet4 = 'YK'
                workbook.SheetNames.push(sheet4)
                workbook.Sheets[sheet4] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '功能码',
                    E: '寄存器地址',
                    F: '比特位',
                    G: '数据类型',
                    H: '合闸',
                    I: '分闸',
                    // J: '数据类型对照表',
                }
                workbook.Sheets[sheet4]['!ref'] = 'A1:' + 'I' + this.yks.length + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet4][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.yks.length; i++) {
                    workbook.Sheets[sheet4]['A' + (i + 1)] = {
                        v: this.yks[i - 1].Name
                    }
                    workbook.Sheets[sheet4]['B' + (i + 1)] = {
                        v: this.yks[i - 1].ID
                    }
                    workbook.Sheets[sheet4]['C' + (i + 1)] = {
                        v: this.yks[i - 1].Scale
                    }
                    workbook.Sheets[sheet4]['D' + (i + 1)] = {
                        v: this.yks[i - 1].Data.Funcode
                    }
                    workbook.Sheets[sheet4]['E' + (i + 1)] = {
                        v: this.yks[i - 1].Data.Addr
                    }
                    workbook.Sheets[sheet4]['F' + (i + 1)] = {
                        v: this.yks[i - 1].Data.Bit
                    }
                    workbook.Sheets[sheet4]['G' + (i + 1)] = {
                        v: this.yks[i - 1].Data.Type
                    }
                    workbook.Sheets[sheet4]['H' + (i + 1)] = {
                        v: this.yks[i - 1].Data.On
                    }
                    workbook.Sheets[sheet4]['I' + (i + 1)] = {
                        v: this.yks[i - 1].Data.Off
                    }
                }
                // workbook.Sheets[sheet4]['J' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet4]['J' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet4]['J' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet4]['J' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet4]['J' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet4]['J' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet4]['J' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet4]['J' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet4]['J' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet4]['J' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet4]['J' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet4]['J' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet4]['J' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet4]['J' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet4]['J' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet4]['J' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet4]['J' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet4]['J' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet4]['J' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet4]['J' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            } else {
                var sheet4 = 'YK'
                workbook.SheetNames.push(sheet4)
                workbook.Sheets[sheet4] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '功能码',
                    E: '寄存器地址',
                    F: '比特位',
                    G: '数据类型',
                    H: '合闸',
                    I: '分闸',
                    // J: '数据类型对照表',
                }
                workbook.Sheets[sheet4]['!ref'] = 'A1:' + 'I' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet4][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                // workbook.Sheets[sheet4]['J' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet4]['J' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet4]['J' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet4]['J' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet4]['J' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet4]['J' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet4]['J' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet4]['J' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet4]['J' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet4]['J' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet4]['J' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet4]['J' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet4]['J' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet4]['J' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet4]['J' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet4]['J' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet4]['J' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet4]['J' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet4]['J' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet4]['J' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            }
            if (this.yts != null) {
                var sheet5 = 'YT'
                workbook.SheetNames.push(sheet5)
                workbook.Sheets[sheet5] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '功能码',
                    E: '寄存器地址',
                    F: '比特位',
                    G: '数据类型',
                    // H: '数据类型对照表',
                }
                workbook.Sheets[sheet5]['!ref'] = 'A1:' + 'G' + this.yts.length + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet5][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                for (var i = 1; i <= this.yts.length; i++) {
                    workbook.Sheets[sheet5]['A' + (i + 1)] = {
                        v: this.yts[i - 1].Name
                    }
                    workbook.Sheets[sheet5]['B' + (i + 1)] = {
                        v: this.yts[i - 1].ID
                    }
                    workbook.Sheets[sheet5]['C' + (i + 1)] = {
                        v: this.yts[i - 1].Scale
                    }
                    workbook.Sheets[sheet5]['D' + (i + 1)] = {
                        v: this.yts[i - 1].Data.Funcode
                    }
                    workbook.Sheets[sheet5]['E' + (i + 1)] = {
                        v: this.yts[i - 1].Data.Addr
                    }
                    workbook.Sheets[sheet5]['F' + (i + 1)] = {
                        v: this.yts[i - 1].Data.Bit
                    }
                    workbook.Sheets[sheet5]['G' + (i + 1)] = {
                        v: this.yts[i - 1].Data.Type
                    }
                }
                // workbook.Sheets[sheet5]['H' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet5]['H' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet5]['H' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet5]['H' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet5]['H' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet5]['H' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet5]['H' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet5]['H' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet5]['H' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet5]['H' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet5]['H' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet5]['H' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet5]['H' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet5]['H' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet5]['H' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet5]['H' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet5]['H' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet5]['H' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet5]['H' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet5]['H' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
            } else {
                var sheet5 = 'YT'
                workbook.SheetNames.push(sheet5)
                workbook.Sheets[sheet5] = {}
                let colOne = {
                    A: '数据名称',
                    B: '点号',
                    C: '系数',
                    D: '功能码',
                    E: '寄存器地址',
                    F: '比特位',
                    G: '数据类型',
                    // H: '数据类型对照表',
                }
                workbook.Sheets[sheet5]['!ref'] = 'A1:' + 'G' + 1
                for (var colName in colOne) {
                    workbook.Sheets[sheet5][colName + 1] = {
                        v: colOne[colName]
                    }
                }
                // workbook.Sheets[sheet5]['H' + 2] = {
                //     v: 'bit位型:0'
                // }
                // workbook.Sheets[sheet5]['H' + 3] = {
                //     v: '16位有符号整型高字节在前:1'
                // }
                // workbook.Sheets[sheet5]['H' + 4] = {
                //     v: '16位有符号整型低字节在前:2'
                // }
                // workbook.Sheets[sheet5]['H' + 5] = {
                //     v: '16位无符号整型高字节在前:3'
                // }
                // workbook.Sheets[sheet5]['H' + 6] = {
                //     v: '16位无符号整型低字节在前:4'
                // }
                // workbook.Sheets[sheet5]['H' + 7] = {
                //     v: '16位浮点型高字节在前:5'
                // }
                // workbook.Sheets[sheet5]['H' + 8] = {
                //     v: '16位浮点型低字节在前:6'
                // }
                // workbook.Sheets[sheet5]['H' + 9] = {
                //     v: '32位有符号整型高寄存器高字节在前:7'
                // }
                // workbook.Sheets[sheet5]['H' + 10] = {
                //     v: '32位有符号整型高寄存器低字节在前:8'
                // }
                // workbook.Sheets[sheet5]['H' + 11] = {
                //     v: '32位有符号整形低寄存器高字节在前:9'
                // }
                // workbook.Sheets[sheet5]['H' + 12] = {
                //     v: '32位有符号整形低寄存器低字节在前:10'
                // }
                // workbook.Sheets[sheet5]['H' + 13] = {
                //     v: '32位无符号整型高寄存器高字节在前:11'
                // }
                // workbook.Sheets[sheet5]['H' + 14] = {
                //     v: '32位无符号整型高寄存器低字节在前:12'
                // }
                // workbook.Sheets[sheet5]['H' + 15] = {
                //     v: '32位无符号整型低寄存器高字节在前:13'
                // }
                // workbook.Sheets[sheet5]['H' + 16] = {
                //     v: '32位无符号整型低寄存器低字节在前:14'
                // }
                // workbook.Sheets[sheet5]['H' + 17] = {
                //     v: '32位浮点型高寄存器高字节在前:15'
                // }
                // workbook.Sheets[sheet5]['H' + 18] = {
                //     v: '32位浮点型高寄存器低字节在前:16'
                // }
                // workbook.Sheets[sheet5]['H' + 19] = {
                //     v: '32位浮点型低寄存器高字节在前:17'
                // }
                // workbook.Sheets[sheet5]['H' + 20] = {
                //     v: '32位浮点型低寄存器低字节在前:18'
                // }
                // workbook.Sheets[sheet5]['H' + 21] = {
                //     v: '64位无符号整型高寄存器高字节在前:19'
                // }
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
            let name = this.tname
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
        showLabel: function(type){
          var label = ""
          switch (type) {
              case "0":
                  label = "bit位型:0"
                  break;
              case "1":
                  label = "16位有符号整型高字节在前:1"
                  break;
              case "2":
                  label = "16位有符号整型低字节在前:2"
                  break;
              case "3":
                  label = "16位无符号整型高字节在前:3"
                  break;
              case "4":
                  label = "16位无符号整型低字节在前:4"
                  break;
              case "5":
                  label = "16位浮点型高字节在前:5"
                  break;
              case "6":
                  label = "16位浮点型低字节在前:6"
                  break;
              case "7":
                  label = "32位有符号整型高寄存器高字节在前:7"
                  break;
              case "8":
                  label = "32位有符号整型高寄存器低字节在前:8"
                  break;
              case "9":
                  label = "32位有符号整型低寄存器高字节在前:9"
                  break;
              case "10":
                  label = "32位有符号整型低寄存器低字节在前:10"
                  break;
              case "11":
                  label = "32位无符号整型高寄存器高字节在前:11"
                  break;
              case "12":
                  label = "32位无符号整型高寄存器低字节在前:12"
                  break;
              case "13":
                  label = "32位无符号整型低寄存器高字节在前:13"
                  break;
              case "14":
                  label = "32位无符号整型低寄存器低字节在前:14"
                  break;
              case "15":
                  label = "32位浮点型高寄存器高字节在前:15"
                  break;
              case "16":
                  label = "32位浮点型高寄存器低字节在前:16"
                  break;
              case "17":
                  label = "32位浮点型低寄存器高字节在前:17"
                  break;
              case "18":
                  label = "32位浮点型低寄存器低字节在前:18"
                  break;
              default:
                  label = "64位无符号整型高寄存器高字节在前:19"
          }
          return label
        },
        loadData: function() {
            var url = '/templates/' + this.tname
            this.$http.get(url).then(function(res) {
                if (res.data.error) console.log(res.data.msg)
                this.ycs = res.data.data.Ycs
                this.yxs = res.data.data.Yxs
                this.yms = res.data.data.Yms
                this.yts = res.data.data.Yts
                this.yks = res.data.data.Yks
            }, function() {
                console.log("加载当前模板服务器异常")
            })
        },
        temp: function() {
            return {
                ID: 0,
                Name: '',
                Channel: '',
                Potocol: '',
                ycs: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yxs: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yms: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }],
                yks: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: '',
                        On: '',
                        Off: ''
                    }
                }],
                yts: [{
                    ID: 0,
                    Name: '',
                    Scale: 0.0,
                    Data: {
                        Partition: '',
                        Funcode: '',
                        Addr: '',
                        Bit: '',
                        Type: ''
                    }
                }]
            }
        }
    },
    created() {
        this.tname = this.$route.params.tname
        this.tid = this.$route.params.tid
        this.tchannel = this.$route.params.tchannel
        this.tpotocol = this.$route.params.tpotocol
        this.current.ID = this.tid
        this.current.Name = this.tname
        this.current.Channel = this.tchannel
        this.current.Potocol = this.tpotocol
            // this.xlsxInstance = xlsx
        this.loadData()
    }
}

</script>
