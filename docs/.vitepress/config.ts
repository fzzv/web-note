import { nav } from './nav'
import {
  vue3,
  interview,
  flutter,
  js,
  java,
  minivue3,
  webpack,
  vite,
  node,
  docker,
  canvas,
  database,
  nest,
  go,
} from './sidebar'

export default {
  title: 'Web-Notes',
  description: '一些学习笔记.',
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: 'logo.png' }]
  ],
  themeConfig: {
    siteTitle: 'Web-Notes',
    logo: '/logo.png',
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2022-present Fan'
    },
    nav: nav(),
    sidebar: {
      '/JS/': js(),
      '/vue3/': vue3(),
      '/vue3/minivue3/': minivue3(),
      '/interview/': interview(),
      '/flutter': flutter(),
      'canvas': canvas(),
      '/node': node(),
      '/java': java(),
      '/docker': docker(),
      '/frontend-engineering/webpack': webpack(),
      '/frontend-engineering/vite': vite(),
      '/database': database(),
      '/nest': nest(),
      '/go': go(),
    },
    outline: {
      level: 'deep',
      label: '目录'
    }
  },
  // 使用 v-pre 属性来避免代码块的解析
  markdown: {
    config(md) {
      const defaultCodeInline = md.renderer.rules.code_inline!
      // 重写代码块的解析规则
      md.renderer.rules.code_inline = (tokens, idx, options, env, self) => {
        // 添加 v-pre 属性
        tokens[idx].attrSet('v-pre', '')
        return defaultCodeInline(tokens, idx, options, env, self)
      }
    }
  }
}
