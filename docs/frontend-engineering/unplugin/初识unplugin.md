# 初识 unplugin

> unplugin 用于构建兼容多平台的插件（如 Vite 插件、Webpack 插件等），你只写一次逻辑，它会自动适配不同构建系统

## 安装

```bash
npm install -D unplugin
```

## 基本使用

```ts
import { defineConfig } from 'unplugin'
import { unplugin } from 'unplugin/config'

export default defineConfig({
  plugins: [
    unplugin({
      name: 'my-plugin',
      entry: 'src/index.ts',
      output: 'dist/index.js',
    }),
  ],
})
```





