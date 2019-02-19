<style scoped>



</style>

<template>

<div class="jumbotron bg-transparent text-center margin-none">
    <h1>本机配置</h1>

  <div class="col-sm-8 col-sm-offset-2">
    <h4 class="page-section-heading">模板配置</h4>
    <div class="panel panel-default">
        <div class="panel-body buttons-spacing-vertical">
            <!-- <p>Stacked Buttons with two lines of text and icon. All regular button styles can be used with stacked buttons.</p> -->
            <p>
                <button class="btn btn-success btn-icon-stacked"><i class="fa fa-2x fa-arrow-circle-down"></i> <span>load mould<br>导入模板</span></button>
                <button class="btn btn-info btn-icon-stacked"><i class="fa fa-2x fa-share-square-o"></i> <span>export mould<br>导出模板</span></button>
                <button class="btn btn-default btn-icon-stacked"><i class="fa fa-2x fa-key "></i> <span>key register<br>产品注册</span></button>

            </p>
        </div>
    </div>
    <h4 class="page-section-heading">系统维护</h4>

    <div class="panel panel-default">

        <div class="panel-body buttons-spacing-vertical">
            <!-- <p>Stacked Buttons with two lines of text and icon. All regular button styles can be used with stacked buttons.</p> -->
            <p>
                <button class="btn btn-success btn-icon-stacked"><i class="fa fa-2x fa-arrow-circle-up"></i> <span>system updata <br>系统更新</span></button>
                <button class="btn btn-info btn-icon-stacked"><i class="fa fa-2x  fa-download"></i> <span>system setup<br>开始安装</span></button>
                <button class="btn btn-default btn-icon-stacked"><i class="fa fa-2x fa-refresh"></i> <span>restore setting<br>出厂恢复</span></button>
            </p>
        </div>
    </div>
  </div>

</div>

</template>

<script>

export default {
    name: 'config',
    data() {
        return {
            msg: '欢迎使用通讯管理机配置工具'
        }
    }
}

</script>
