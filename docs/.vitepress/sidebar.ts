export const js = () => {
  return [
    {
      text: '重学Js',
      items: [
        { text: '事件', link: '/JS/events' },
        { text: '异步编程', link: '/JS/async-programming' },
        { text: '高阶编程技巧', link: '/JS/advanced-programming' },
        { text: '面向对象编程', link: '/JS/oop' },
        { text: '浏览器底层渲染', link: '/JS/browser-render' },
        { text: 'http网络通信层' , link: '/JS/http-network'},
        { text: 'Proxy和Reflect' , link: '/JS/proxy-reflect' },
      ]
    }
  ]
}

export const vue3 = () => {
  return [
    {
      text: 'Vue3相关',
      items: [
        { text: 'vue3的setup函数', link: '/vue3/setup-function' },
        { text: 'vue3的响应式API', link: '/vue3/composition-api' },
        { text: 'vue3+ts+elementplus项目搭建', link: '/vue3/project-setup' },
      ]
    }
  ]
}

export const minivue3 = () => {
  return [
    {
      text: '手写minivue3',
      items: [
        { text: 'Vue2 对比 Vue3', link: '/vue3/minivue3/minivue3-v2vsv3' },
        { text: '环境搭建', link: '/vue3/minivue3/minivue3-env' },
        { text: 'Reactivity响应式模块', link: '/vue3/minivue3/minivue3-reactivity' },
      ]
    }
  ]
}

export const interview = () => {
  return [
    {
      text: '面试题',
      items: [
        { text: 'html', link: '/interview/html' },
        { text: 'css', link: '/interview/css' },
        { text: 'js', link: '/interview/js' },
        { text: 'vue', link: '/interview/vue' },
        { text: 'vue3', link: '/interview/vue3' },
        { text: 'node.js', link: '/interview/nodejs' },
        { text: 'MySQL', link: '/interview/MySQL' },
        { text: '移动端和小程序', link: '/interview/移动端和小程序' },
        { text: '计算机网络', link: '/interview/计算机网络' },
        { text: '前端工程化', link: '/interview/前端工程化' },
        { text: '手写代码', link: '/interview/手写代码' },
        { text: '项目描述', link: '/interview/项目描述' },
        { text: '开发经验', link: '/interview/开发经验' },
        { text: '面试技巧', link: '/interview/面试技巧' },
        { text: '面试复盘', link: '/interview/面试复盘' },
      ]
    }
  ]
}

export const flutter = () => {
  return [
    {
      text: 'Flutter基础',
      items: [
        { text: '初始化项目', link: '/flutter/init' },
        { text: '组件Widget', link: '/flutter/widget' },
        { text: '文本', link: '/flutter/text' },
        { text: '页面布局', link: '/flutter/layout' },
        { text: '按钮', link: '/flutter/button' },
        { text: '图片', link: '/flutter/image' },
        { text: '列表', link: '/flutter/list' },
        { text: '第三方库组件', link: '/flutter/otherWidget' },
        { text: '状态管理', link: '/flutter/state' },
        { text: '路由', link: '/flutter/route' },
        { text: '导航', link: '/flutter/nav' },
        { text: '表单', link: '/flutter/form' },
      ]
    },
    {
      text: 'Flutter实战',
      items: [
        { text: 'Tab导航', link: '/flutter/sz_nav' },
        { text: 'Fluro路由', link: '/flutter/sz_route' },
        { text: '状态管理', link: '/flutter/sz_state' },
        { text: '接口调用', link: '/flutter/sz_dio' },
        { text: '首页', link: '/flutter/sz_home' },
        { text: '登录', link: '/flutter/sz_login' },
      ]
    },
  ]
}

export const canvas = () => {
  return [
    {
      text: 'Canvas',
      items: [
        { text: 'canvas入门', link: '/canvas/01.canvas入门' },
      ]
    }
  ]
}

export const node = () => {
  return [
    {
      text: 'Node基础',
      items: [
        { text: '1.REPL', link: '/node/1.REPL' },
        { text: 'htmx', link: '/node/htmx' },
        { text: 'handlebars', link: '/node/handlebars' },
        { text: 'winston', link: '/node/winston' },
      ]
    },
  ]
}

export const nest = () => {
  return [
    {
      text: 'NestJS',
      items: [
        { text: 'NestJS基础', link: '/nest/nest基础' },
        { text: '装饰器', link: '/nest/装饰器' },
        { text: '依赖注入和控制反转', link: '/nest/依赖注入和控制反转' },
        { text: 'SOLID原则', link: '/nest/solid原则' },
        { text: 'TypeORM', link: '/nest/TypeORM' },
        { text: '小案例', link: '/nest/小案例' },
      ]
    },
    {
      text: 'Kysely',
      items: [
        { text: '使用kysely进行简单的curd', link: '/nest/kysely/使用kysely进行简单的curd' },
        { text: 'kysely中常用的method', link: '/nest/kysely/kysely中常用的method' },
      ]
    },
  ]
}

export const java = () => {
  return [
    {
      text: 'SpringBoot',
      items: [
        { text: 'SpringBoot基础', link: '/java/SpringBoot基础' },
        { text: 'SpringBoot基础项目', link: '/java/SpringBoot基础项目' },
        { text: 'SpringBoot秒杀项目', link: '/java/SpringBoot秒杀项目' },
        { text: 'SpringBoot租房项目', link: '/java/SpringBoot租房项目' },
      ]
    },
    {
      text: 'Maven',
      items: [
        { text: 'Maven', link: '/java/maven' },
      ]
    },
    {
      text: '其他',
      items: [
        { text: 'Mybatis Plus 入门', link: '/java/mybatisplus' },
      ]
    },
  ]
}

export const docker = () => {
  return [
    {
      text: 'Docker',
      items: [
        { text: '基本概念', link: '/docker/01.基本概念' },
        { text: 'docker安装', link: '/docker/02.docker安装' },
        { text: 'docker命令', link: '/docker/03.docker命令' },
        { text: 'docker中安装nginx', link: '/docker/docker中安装nginx' },
      ]
    }
  ]
}

export const webpack = () => {
  return [
    {
      text: 'Webpack',
      items: [
        { text: '初识webpack', link: '/frontend-engineering/webpack/1.初始webpack' },
        { text: '处理css资源', link: '/frontend-engineering/webpack/2.处理css资源' },
        { text: '处理其他资源', link: '/frontend-engineering/webpack/3.处理其他资源' },
        { text: '构建现代js环境', link: '/frontend-engineering/webpack/4.构建现代js环境' },
        { text: 'npm库和source-map', link: '/frontend-engineering/webpack/5.npm库和source-map' },
        { text: 'config结构解读', link: '/frontend-engineering/webpack/6.配置结构' },
        { text: '本地开发服务器和HMR', link: '/frontend-engineering/webpack/7.本地开发服务器和HMR' },
        { text: 'webpack模块联邦', link: '/frontend-engineering/webpack/8.webpack模块联邦' },
      ]
    },
  ]
}

export const vite = () => {
  return [
    {
      text: '其他',
      items: [
        { text: '打包成app遇到的问题', link: '/frontend-engineering/vite/打包成app遇到的问题' }
      ]
    },
  ]
}

export const database = () => {
  return [
    {
      text: '数据库',
      items: [
        { text: '数据库设计', link: '/database/mysql/数据库设计' },
      ]
    },
    {
      text: 'MySQL',
      items: [
        { text: 'MySQL的操作', link: '/database/mysql/MySQL操作' },
      ]
    },
  ]
}

