import { defineConfig } from 'vite'
import { LogPrefixUnplugin } from './unplugin'

export default defineConfig({
  plugins: [
    LogPrefixUnplugin.vite({ prefix: '[DEV1ertert23]' }) // 使用插件并传递参数
  ]
})
