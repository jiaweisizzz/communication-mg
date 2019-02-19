<template>
    <div class="jumbotron bg-transparent text-center margin-none">
      <h1>帮助文档</h1>
      <p>{{msg}}</p>
    </div>
</template>

<script>
export default {
  name: 'help',
  data () {
    return {
      msg: '欢迎使用通讯管理机配置工具'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
