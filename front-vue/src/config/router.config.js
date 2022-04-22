// eslint-disable-next-line
import { UserLayout, BasicLayout, BlankLayout } from '@/layouts'
const RouteView = {
  name: 'RouteView',
  render: h => h('router-view')
}
export const asyncRouterMap = [
  {
    path: '/',
    name: 'index',
    component: BasicLayout,
    meta: { title: 'home' },
    redirect: '/home',
    children: [
      {
        name: 'home',
        meta: { icon: 'home', title: '首页', show: true },
        path: '/home',
        component: () => import('@/views/home')
      },
      {
        name: 'cmdlist',
        meta: { icon: 'code', title: '计划任务' },
        component: () => import('@/views/command/CmdList'),
        path: '/command/cmdlist'
      },
      {
        name: 'cmdexec',
        hidden: true,
        meta: { title: '执行命令' },
        component: () => import('@/views/command/CmdExec'),
        path: '/command/cmdexec'
      },
      {
        name: 'cmdlog',
        hidden: true,
        meta: { title: '命令日志' },
        component: () => import('@/views/command/CmdLog'),
        path: '/command/cmdlog'
      },
      {
        name: 'plist',
        meta: { icon: 'project', title: '进程管理', show: true },
        component: () => import('@/views/process/List'),
        path: '/process/plist'
      },
      {
        name: 'pedit',
        hidden: true,
        meta: { title: '编辑' },
        component: () => import('@/views/process/Edit'),
        path: '/process/pedit'
      },
      {
        name: 'flyfish',
        meta: { icon: 'database', title: 'Flyfish', show: true },
        component: () => import('@/views/db/Flyfish'),
        path: '/db/flyfish'
      },
      {
        name: 'gamemaster',
        meta: { icon: 'database', title: 'GameMaster', show: true },
        component: () => import('@/views/gm/Gm'),
        path: '/gm/gamemaster'
      },
      {
        path: '/system',
        redirect: '/system/usermange',
        component: RouteView,
        meta: { icon: 'control', title: '系统管理', show: true },
        children: [
          {
            path: '/system/usermange',
            name: 'usermange',
            component: () => import('@/views/system/usermange'),
            meta: { title: '用户管理', show: true }
          }
        ]
      }
    ]
  },
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

/**
 * 基础路由
 * @type { *[] }
 */
export const constantRouterMap = [
  {
    path: '/user',
    component: UserLayout,
    redirect: '/user/login',
    hidden: true,
    children: [
      {
        path: '/user/login',
        name: 'login',
        component: () => import(/* webpackChunkName: "user" */ '@/views/user/Login')
      }
    ]
  },
  {
    path: '/404',
    component: () => import(/* webpackChunkName: "fail" */ '@/views/exception/404')
  }
]
