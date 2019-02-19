import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/page/Home'
import List from '@/components/page/List'
import Templatelist from '@/components/page/Templatelist'
import Point from '@/components/page/Point'
import Config from '@/components/page/Config'
import User from '@/components/page/User'
import Channel from '@/components/page/Channel'
import Version from '@/components/page/Version'
import Template from '@/components/page/Template'
import Help from '@/components/page/Help'
// import Login from '@/components/page/Login'
import Notice from '@/components/page/Notice'
import Pointtable from '@/components/page/Pointtable'

Vue.use(Router)

export default new Router({
  routes: [
    // {
    //   path: '/',
    //   name: 'login',
    //   component: Login
    // },
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/notice',
        name: 'notice',
        component: Notice
    },
    {
      path: '/user',
      name: 'user',
      component: User
    },
    {
      path: '/config',
      name: 'config',
      component: Config
    },
    {
      path: '/point/:dname/:dchannel/:did',
      name: 'point',
      component: Point
    },
    {
      path: '/channel',
      name: 'channel',
      component: Channel
    },
    {
      path: '/list/:cname',
      name: 'list',
      component: List
    },
    {
      path: '/pointtable/:cname',
      name: 'pointtable',
      component: Pointtable
    },
    {
      path: '/templatelist/:tname/:tid/:tchannel',
      name: 'templatelist',
      component: Templatelist
    },
    {
      path: '/version',
      name: 'version',
      component: Version
    },
    {
      path: '/template',
      name: 'template',
      component: Template
    },
    {
      path: '/help',
      name: 'help',
      component: Help
    }
  ]
})
