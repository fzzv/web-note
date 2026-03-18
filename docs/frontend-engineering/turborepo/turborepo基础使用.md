# Turborepo 基础使用指南

## 什么是 Turborepo？

Turborepo 是由 Vercel 开发的高性能构建系统，专门用于 JavaScript 和 TypeScript 代码库。它是一个针对 monorepo 的增量构建工具，可以显著提高大型项目的构建和开发效率。

### 核心特性

- **增量构建**：只构建发生变化的包
- **远程缓存**：在团队和 CI/CD 之间共享构建缓存
- **并行执行**：智能任务调度和并行处理
- **任务管道**：定义任务之间的依赖关系
- **代码生成**：内置代码生成器
- **零配置**：开箱即用的合理默认配置

## 安装 Turborepo

### 全局安装

```bash
# 使用 npm
npm install -g turbo

# 使用 yarn
yarn global add turbo

# 使用 pnpm
pnpm add -g turbo
```

### 项目中安装

```bash
# 使用 npm
npm install -D turbo

# 使用 yarn
yarn add -D turbo

# 使用 pnpm
pnpm add -D turbo
```

## 创建 Turborepo 项目

### 使用官方模板创建

```bash
# 创建新的 turborepo 项目
npx create-turbo@latest

# 或者使用特定模板
npx create-turbo@latest --example basic
npx create-turbo@latest --example with-nextjs
npx create-turbo@latest --example with-react
```

### 手动初始化现有项目

如果你已经有一个 monorepo 项目，可以手动添加 Turborepo：

1. 安装 turbo
2. 创建 `turbo.json` 配置文件
3. 更新 package.json 脚本

## 项目结构

典型的 Turborepo 项目结构如下：

```
my-turborepo/
├── apps/
│   ├── web/                 # Next.js 应用
│   │   ├── package.json
│   │   └── ...
│   └── docs/                # 文档站点
│       ├── package.json
│       └── ...
├── packages/
│   ├── ui/                  # 共享 UI 组件库
│   │   ├── package.json
│   │   └── ...
│   ├── eslint-config/       # 共享 ESLint 配置
│   │   ├── package.json
│   │   └── ...
│   └── typescript-config/   # 共享 TypeScript 配置
│       ├── package.json
│       └── ...
├── package.json             # 根 package.json
├── turbo.json              # Turborepo 配置
└── yarn.lock               # 锁定文件
```

## 配置文件详解

### turbo.json 配置

`turbo.json` 是 Turborepo 的核心配置文件：

```json
{
  "$schema": "https://turbo.build/schema.json",
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".next/**", "!.next/cache/**"]
    },
    "test": {
      "dependsOn": ["build"],
      "outputs": ["coverage/**"],
      "inputs": ["src/**/*.tsx", "src/**/*.ts", "test/**/*.ts"]
    },
    "lint": {
      "outputs": []
    },
    "dev": {
      "cache": false,
      "persistent": true
    }
  },
  "globalDependencies": ["**/.env.*local"],
  "globalEnv": ["NODE_ENV"]
}
```

#### 配置选项说明

- **pipeline**：定义任务管道和依赖关系
- **dependsOn**：定义任务依赖
  - `^build`：表示依赖其他包的 build 任务
  - `build`：表示依赖当前包的 build 任务
- **outputs**：指定任务输出的文件或目录
- **inputs**：指定任务输入的文件（用于缓存判断）
- **cache**：是否缓存任务结果
- **persistent**：是否为持久任务（如 dev server）
- **globalDependencies**：全局依赖文件
- **globalEnv**：全局环境变量

### 根 package.json 配置

```json
{
  "name": "my-turborepo",
  "private": true,
  "workspaces": [
    "apps/*",
    "packages/*"
  ],
  "scripts": {
    "build": "turbo run build",
    "dev": "turbo run dev",
    "lint": "turbo run lint",
    "test": "turbo run test",
    "clean": "turbo run clean"
  },
  "devDependencies": {
    "turbo": "latest"
  },
  "packageManager": "yarn@1.22.19"
}
```

## 基本使用命令

### 运行任务

```bash
# 运行所有包的 build 任务
turbo run build

# 运行特定包的任务
turbo run build --filter=web

# 运行多个任务
turbo run build test lint

# 并行运行开发服务器
turbo run dev --parallel

# 运行任务并显示详细输出
turbo run build --verbose
```

### 过滤器（Filters）

Turborepo 提供了强大的过滤器功能：

```bash
# 只运行特定包
turbo run build --filter=web

# 运行多个包
turbo run build --filter=web --filter=api

# 使用通配符
turbo run build --filter="@myorg/*"

# 运行包及其依赖
turbo run build --filter=web...

# 运行包及其依赖者
turbo run build --filter=...web

# 基于 git 变更过滤
turbo run build --filter=[HEAD^1]

# 排除特定包
turbo run build --filter=!docs
```

### 缓存管理

```bash
# 查看缓存统计
turbo run build --dry-run

# 强制重新构建（忽略缓存）
turbo run build --force

# 清理本地缓存
turbo prune

# 查看缓存配置
turbo run build --dry-run=json
```

## 远程缓存配置

### Vercel 远程缓存

```bash
# 登录 Vercel
npx turbo login

# 链接项目
npx turbo link

# 现在构建会自动使用远程缓存
turbo run build
```

### 自定义远程缓存

在 `turbo.json` 中配置：

```json
{
  "remoteCache": {
    "signature": true
  }
}
```

## 代码生成器

### 创建生成器

```bash
# 创建新的生成器
turbo gen workspace --name=my-package --type=package
```

### 自定义生成器

在项目根目录创建 `turbo/generators` 目录：

```typescript
// turbo/generators/config.ts
import type { PlopTypes } from "@turbo/gen";

export default function generator(plop: PlopTypes.NodePlopAPI): void {
  plop.setGenerator("component", {
    description: "创建新的 React 组件",
    prompts: [
      {
        type: "input",
        name: "name",
        message: "组件名称？"
      }
    ],
    actions: [
      {
        type: "add",
        path: "packages/ui/src/{{pascalCase name}}.tsx",
        templateFile: "templates/component.hbs"
      }
    ]
  });
}
```

## 最佳实践

### 1. 合理组织项目结构

```
├── apps/           # 应用程序
├── packages/       # 可复用的包
├── tools/          # 工具和配置
└── docs/           # 文档
```

### 2. 配置合适的任务依赖

```json
{
  "pipeline": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**"]
    },
    "test": {
      "dependsOn": ["build"],
      "outputs": ["coverage/**"]
    }
  }
}
```

### 3. 使用环境变量

```json
{
  "globalEnv": [
    "NODE_ENV",
    "DATABASE_URL"
  ]
}
```

### 4. 优化缓存配置

```json
{
  "pipeline": {
    "build": {
      "inputs": [
        "src/**/*.ts",
        "src/**/*.tsx",
        "!src/**/*.test.ts"
      ],
      "outputs": ["dist/**"]
    }
  }
}
```

### 5. 使用 .turbo 忽略文件

创建 `.turboignore` 文件：

```
# 忽略测试文件
**/*.test.ts
**/*.spec.ts

# 忽略文档
**/README.md
**/docs/**
```

## 常见问题和解决方案

### 1. 缓存问题

如果遇到缓存问题，可以：

```bash
# 清理缓存
turbo prune

# 强制重新构建
turbo run build --force
```

### 2. 依赖问题

确保正确配置 workspace 依赖：

```json
{
  "dependencies": {
    "@myorg/ui": "workspace:*"
  }
}
```

### 3. 性能优化

- 合理设置并行任务数量
- 优化任务依赖关系
- 使用远程缓存
- 配置合适的 inputs 和 outputs

### 4. 调试技巧

```bash
# 查看任务执行计划
turbo run build --dry-run

# 显示详细日志
turbo run build --verbose

# 生成任务图
turbo run build --graph
```

## 与其他工具集成

### ESLint 集成

```json
// packages/eslint-config/package.json
{
  "name": "@myorg/eslint-config",
  "main": "index.js",
  "dependencies": {
    "@typescript-eslint/eslint-plugin": "^5.0.0",
    "@typescript-eslint/parser": "^5.0.0",
    "eslint-config-next": "^12.0.0"
  }
}
```

### TypeScript 集成

```json
// packages/tsconfig/package.json
{
  "name": "@myorg/tsconfig",
  "files": ["base.json", "nextjs.json", "react-library.json"]
}
```

### Jest 集成

```json
{
  "pipeline": {
    "test": {
      "outputs": ["coverage/**"],
      "inputs": [
        "src/**/*.tsx",
        "src/**/*.ts",
        "test/**/*.ts",
        "jest.config.js"
      ]
    }
  }
}
```

## 监控和分析

### 构建分析

```bash
# 生成构建报告
turbo run build --summarize

# 查看任务时间线
turbo run build --timeline
```

### 性能监控

使用 Turborepo 的内置分析工具来监控构建性能：

```bash
# 查看缓存命中率
turbo run build --dry-run | grep "cache hit"

# 分析任务执行时间
turbo run build --timeline
```

## 升级和迁移

### 升级 Turborepo

```bash
# 检查最新版本
npm outdated turbo

# 升级到最新版本
npm update turbo
```

## 总结

Turborepo 是一个强大的 monorepo 构建工具，通过以下特性显著提升开发效率：

1. **智能缓存**：避免重复构建
2. **并行执行**：充分利用系统资源
3. **任务编排**：合理安排任务依赖
4. **远程缓存**：团队协作效率提升
5. **零配置**：开箱即用的体验

掌握 Turborepo 的基础使用后，可以显著提升大型前端项目的开发和构建效率，特别适合多包、多应用的复杂项目场景。
