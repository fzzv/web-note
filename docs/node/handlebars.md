# Handlebars 模板引擎

Handlebars 是一个基于 Mustache 模板语言的 JavaScript 模板引擎，提供了强大的模板功能，广泛用于前端和 Node.js 后端开发。

## 基本概念

Handlebars 使用双花括号 `{{}}` 作为模板语法，支持变量插值、条件语句、循环、部分模板等功能。

## 基础使用

### 1. 简单模板渲染

```javascript
const handlebars = require('handlebars');

// 1. 创建一个模板
const tpl = `
  <h1>{{title}}</h1>
  <p>{{content}}</p>
`;

// 2. 上下文对象
const data = {
  title: 'Hello, Handlebars!',
  content: 'This is a template example.'
};

// 3. 渲染模板
const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
```

**输出结果**：
```html
<h1>Hello, Handlebars!</h1>
<p>This is a template example.</p>
```

### 2. 条件语句

```javascript
const handlebars = require('handlebars');

// if else 语句
const tpl = `
  {{#if isAdmin}}
    <h1>Admin</h1>
  {{else}}
    <h1>User</h1>
  {{/if}}
`;

const data = {
  isAdmin: false
};

const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
```

**输出结果**：
```html
<h1>User</h1>
```

### 3. 循环遍历和自定义 Helper

```javascript
const handlebars = require('handlebars');

// 注册一个全局函数，在模板中可以调用
handlebars.registerHelper('toUpperCase', function (str) {
  return str.toUpperCase();
});

// 循环遍历将所有内容进行大写
const tpl = `
  {{#each items}}
    <h1>{{toUpperCase this}}</h1>
  {{/each}}
`;

const data = {
  items: ['apple', 'banana', 'cherry', 'orange']
};

const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
```

**输出结果**：
```html
<h1>APPLE</h1>
<h1>BANANA</h1>
<h1>CHERRY</h1>
<h1>ORANGE</h1>
```

## Express 集成

### 项目结构
```
handlebars/
├── index.js                 # 主应用文件
├── package.json             # 项目配置
├── views/                   # 视图目录
│   ├── layouts/            # 布局模板
│   │   └── main.handlebars # 主布局
│   ├── partials/           # 部分模板
│   │   ├── header.handlebars
│   │   └── treeNode.handlebars
│   ├── home.handlebars     # 首页模板
│   └── tree.handlebars     # 树形结构模板
└── case/                   # 示例案例
    ├── 1.js               # 基础示例
    ├── 2.js               # 条件语句示例
    └── 3.js               # Helper 示例
```

### Express 配置

```javascript
const express = require('express');
const exphbs = require('express-handlebars');
const path = require('path');

// 创建 Handlebars 实例
const hbs = exphbs.create({});

const app = express();

// 配置 Handlebars 作为模板引擎
app.engine('handlebars', hbs.engine);
app.set('view engine', 'handlebars');
app.set('views', path.join(__dirname, 'views'));

// 路由定义
app.get('/', (req, res) => {
  res.render('home', { title: 'home' });
});

app.get('/tree', (req, res) => {
  const tree = [
    {
      name: 'parent1',
      children: [
        {
          name: 'parent1-child1',
          children: [
            { name: 'parent1-child1-grandson1', children: [] }
          ]
        },
        { name: 'parent1-child2', children: [] }
      ]
    }
  ];
  res.render('tree', { tree });
});

app.listen(3000, () => {
  console.log('Server is running on port 3000');
});
```

## 模板系统

### 1. 主布局模板 (`layouts/main.handlebars`)

```handlebars
<!DOCTYPE html>
<html lang="en">
{{> header}}
<body>
  {{{body}}}
</body>
</html>
```

**特点**：
- `{{> header}}`：引用部分模板
- `{{{body}}}`：三重花括号表示不转义 HTML 内容
- 作为所有页面的基础布局

### 2. 部分模板 (Partials)

#### 头部模板 (`partials/header.handlebars`)
```handlebars
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{title}}</title>
</head>
```

#### 树节点模板 (`partials/treeNode.handlebars`)
```handlebars
<div>
    {{name}}
    {{#if children}}
      <div style="padding-left:{{multiply depth 20}}px">
        {{#each children}}
        {{> treeNode this depth=(increase ../depth)}}
        {{/each}}
      </div>
    {{/if}}
</div>
```

**特点**：
- 递归调用自身 `{{> treeNode this}}`
- 使用自定义 Helper：`multiply` 和 `increase`
- `../depth` 访问父级上下文

### 3. 页面模板

#### 首页模板 (`home.handlebars`)
```handlebars
<h1>{{title}}</h1>
```

#### 树形结构模板 (`tree.handlebars`)
```handlebars
<div>
  {{#each tree}}
    {{> treeNode this depth=1}}
  {{/each}}
</div>
```

## Handlebars 语法详解

### 1. 基础语法

#### 变量插值
```handlebars
{{name}}           <!-- 转义 HTML -->
{{{name}}}         <!-- 不转义 HTML -->
{{../name}}        <!-- 访问父级上下文 -->
{{@index}}         <!-- 内置变量 -->
{{@key}}           <!-- 对象键名 -->
```

#### 路径访问
```handlebars
{{user.name}}      <!-- 对象属性 -->
{{users.[0].name}} <!-- 数组元素 -->
{{lookup users 0}} <!-- 动态查找 -->
```

### 2. 条件语句

```handlebars
<!-- if 语句 -->
{{#if condition}}
  <p>条件为真</p>
{{else}}
  <p>条件为假</p>
{{/if}}

<!-- unless 语句 -->
{{#unless condition}}
  <p>条件为假时显示</p>
{{/unless}}

<!-- 比较操作 -->
{{#if (eq status "active")}}
  <p>状态为活跃</p>
{{/if}}
```

### 3. 循环语句

```handlebars
<!-- each 循环 -->
{{#each items}}
  <div>{{@index}}: {{this}}</div>
{{/each}}

<!-- 对象遍历 -->
{{#each user}}
  <p>{{@key}}: {{this}}</p>
{{/each}}

<!-- 嵌套循环 -->
{{#each categories}}
  <h2>{{name}}</h2>
  {{#each items}}
    <p>{{../name}} - {{this}}</p>
  {{/each}}
{{/each}}
```

### 4. with 语句

```handlebars
{{#with user}}
  <h1>{{name}}</h1>
  <p>{{email}}</p>
{{/with}}
```

## 自定义 Helper

### 1. 基础 Helper

```javascript
// 注册简单 Helper
handlebars.registerHelper('toUpperCase', function (str) {
  return str.toUpperCase();
});

// 注册带参数的 Helper
handlebars.registerHelper('multiply', function (a, b) {
  return a * b;
});

// 注册带选项的 Helper
handlebars.registerHelper('list', function (items, options) {
  let result = '<ul>';
  for (let i = 0; i < items.length; i++) {
    result += '<li>' + options.fn(items[i]) + '</li>';
  }
  result += '</ul>';
  return new handlebars.SafeString(result);
});
```

### 2. 块级 Helper

```javascript
// 注册块级 Helper
handlebars.registerHelper('bold', function (options) {
  return '<strong>' + options.fn(this) + '</strong>';
});

// 使用示例
// {{#bold}}这里的内容会被加粗{{/bold}}
```

### 3. 条件 Helper

```javascript
// 比较 Helper
handlebars.registerHelper('eq', function (a, b) {
  return a === b;
});

handlebars.registerHelper('gt', function (a, b) {
  return a > b;
});

handlebars.registerHelper('lt', function (a, b) {
  return a < b;
});

// 使用示例
// {{#if (gt age 18)}}成年人{{/if}}
```

### 4. 实用 Helper

```javascript
// 日期格式化
handlebars.registerHelper('formatDate', function (date, format) {
  // 使用 moment.js 或其他日期库
  return moment(date).format(format);
});

// JSON 序列化
handlebars.registerHelper('json', function (context) {
  return JSON.stringify(context);
});

// 截取字符串
handlebars.registerHelper('truncate', function (str, length) {
  if (str.length > length) {
    return str.substring(0, length) + '...';
  }
  return str;
});

// 数组长度
handlebars.registerHelper('length', function (array) {
  return array.length;
});
```

## 高级特性

### 1. 预编译模板

```javascript
// 预编译模板以提高性能
const fs = require('fs');
const handlebars = require('handlebars');

// 读取模板文件
const source = fs.readFileSync('template.handlebars', 'utf8');

// 预编译模板
const template = handlebars.precompile(source);

// 保存预编译结果
fs.writeFileSync('template.js', `
  const handlebars = require('handlebars/runtime');
  module.exports = handlebars.template(${template});
`);
```

### 2. 部分模板注册

```javascript
// 注册部分模板
handlebars.registerPartial('header', `
  <head>
    <title>{{title}}</title>
  </head>
`);

// 从文件注册部分模板
const fs = require('fs');
const headerTemplate = fs.readFileSync('partials/header.handlebars', 'utf8');
handlebars.registerPartial('header', headerTemplate);
```

### 3. 安全性

```javascript
// HTML 转义
const safeString = new handlebars.SafeString('<b>Bold</b>');

// 自定义转义函数
handlebars.registerHelper('escape', function (str) {
  return handlebars.Utils.escapeExpression(str);
});
```

## Express-Handlebars 高级配置

### 1. 完整配置选项

```javascript
const exphbs = require('express-handlebars');

const hbs = exphbs.create({
  // 默认布局
  defaultLayout: 'main',
  
  // 布局目录
  layoutsDir: path.join(__dirname, 'views/layouts'),
  
  // 部分模板目录
  partialsDir: path.join(__dirname, 'views/partials'),
  
  // 文件扩展名
  extname: '.handlebars',
  
  // 自定义 Helper
  helpers: {
    toUpperCase: function (str) {
      return str.toUpperCase();
    },
    multiply: function (a, b) {
      return a * b;
    },
    increase: function (value) {
      return parseInt(value) + 1;
    }
  },
  
  // 运行时选项
  runtimeOptions: {
    allowProtoPropertiesByDefault: true,
    allowProtoMethodsByDefault: true
  }
});
```

### 2. 多布局支持

```javascript
// 不同页面使用不同布局
app.get('/admin', (req, res) => {
  res.render('admin', {
    layout: 'admin',  // 使用 admin.handlebars 布局
    title: 'Admin Panel'
  });
});

// 不使用布局
app.get('/api', (req, res) => {
  res.render('api', {
    layout: false,  // 不使用布局
    data: { message: 'API Response' }
  });
});
```

### 3. 动态部分模板

```javascript
// 动态注册部分模板
app.use((req, res, next) => {
  // 根据用户角色动态注册不同的导航模板
  if (req.user && req.user.role === 'admin') {
    hbs.handlebars.registerPartial('navigation', adminNavTemplate);
  } else {
    hbs.handlebars.registerPartial('navigation', userNavTemplate);
  }
  next();
});
```

## 实际应用示例

### 1. 博客系统

```javascript
// 博客文章列表
app.get('/blog', async (req, res) => {
  const posts = await Blog.findAll();
  res.render('blog/index', {
    title: '博客',
    posts: posts,
    helpers: {
      formatDate: function (date) {
        return new Date(date).toLocaleDateString();
      }
    }
  });
});

// 博客文章详情
app.get('/blog/:id', async (req, res) => {
  const post = await Blog.findById(req.params.id);
  res.render('blog/detail', {
    title: post.title,
    post: post
  });
});
```

#### 博客模板 (`blog/index.handlebars`)
```handlebars
<div class="blog-list">
  {{#each posts}}
    <article class="post-preview">
      <h2><a href="/blog/{{id}}">{{title}}</a></h2>
      <p class="meta">发布于 {{formatDate createdAt}}</p>
      <p>{{excerpt}}</p>
    </article>
  {{/each}}
</div>
```

### 2. 用户管理系统

```javascript
// 用户列表
app.get('/users', async (req, res) => {
  const users = await User.findAll();
  res.render('users/index', {
    title: '用户管理',
    users: users
  });
});
```

#### 用户模板 (`users/index.handlebars`)
```handlebars
<table class="users-table">
  <thead>
    <tr>
      <th>ID</th>
      <th>姓名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{id}}</td>
      <td>{{name}}</td>
      <td>{{email}}</td>
      <td>
        {{#if isActive}}
          <span class="status active">活跃</span>
        {{else}}
          <span class="status inactive">非活跃</span>
        {{/if}}
      </td>
      <td>
        <a href="/users/{{id}}/edit">编辑</a>
        <a href="/users/{{id}}/delete">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
```

### 3. 表单处理

```javascript
// 显示表单
app.get('/contact', (req, res) => {
  res.render('contact', {
    title: '联系我们'
  });
});

// 处理表单提交
app.post('/contact', (req, res) => {
  const { name, email, message } = req.body;
  // 处理表单数据...
  
  res.render('contact', {
    title: '联系我们',
    success: '消息发送成功！',
    formData: { name, email, message }
  });
});
```

#### 联系表单模板 (`contact.handlebars`)
```handlebars
{{#if success}}
  <div class="alert alert-success">{{success}}</div>
{{/if}}

{{#if errors}}
  <div class="alert alert-error">
    <ul>
      {{#each errors}}
        <li>{{this}}</li>
      {{/each}}
    </ul>
  </div>
{{/if}}

<form method="POST" action="/contact">
  <div class="form-group">
    <label for="name">姓名</label>
    <input type="text" id="name" name="name" value="{{formData.name}}" required>
  </div>
  
  <div class="form-group">
    <label for="email">邮箱</label>
    <input type="email" id="email" name="email" value="{{formData.email}}" required>
  </div>
  
  <div class="form-group">
    <label for="message">消息</label>
    <textarea id="message" name="message" required>{{formData.message}}</textarea>
  </div>
  
  <button type="submit">发送</button>
</form>
```

## 性能优化

### 1. 模板缓存

```javascript
// 生产环境启用模板缓存
const hbs = exphbs.create({
  defaultLayout: 'main',
  // 生产环境缓存模板
  cache: process.env.NODE_ENV === 'production'
});
```

### 2. 预编译模板

```javascript
// 构建脚本：预编译所有模板
const fs = require('fs');
const path = require('path');
const handlebars = require('handlebars');

function precompileTemplates(dir) {
  const files = fs.readdirSync(dir);
  
  files.forEach(file => {
    if (path.extname(file) === '.handlebars') {
      const source = fs.readFileSync(path.join(dir, file), 'utf8');
      const template = handlebars.precompile(source);
      const outputPath = path.join(dir, file.replace('.handlebars', '.js'));
      
      fs.writeFileSync(outputPath, `
        const handlebars = require('handlebars/runtime');
        module.exports = handlebars.template(${template});
      `);
    }
  });
}

precompileTemplates('./views');
```

### 3. 部分模板优化

```javascript
// 懒加载部分模板
const partialCache = new Map();

handlebars.registerHelper('lazyPartial', function (name, context) {
  if (!partialCache.has(name)) {
    const source = fs.readFileSync(`./views/partials/${name}.handlebars`, 'utf8');
    partialCache.set(name, handlebars.compile(source));
  }
  
  const template = partialCache.get(name);
  return new handlebars.SafeString(template(context));
});
```

## 最佳实践

### 1. 项目结构组织

```
views/
├── layouts/              # 布局模板
│   ├── main.handlebars  # 主布局
│   ├── admin.handlebars # 管理后台布局
│   └── mobile.handlebars # 移动端布局
├── partials/            # 部分模板
│   ├── header.handlebars
│   ├── footer.handlebars
│   ├── navigation.handlebars
│   └── forms/           # 表单相关
│       ├── input.handlebars
│       └── button.handlebars
├── pages/               # 页面模板
│   ├── home.handlebars
│   ├── about.handlebars
│   └── contact.handlebars
└── components/          # 组件模板
    ├── card.handlebars
    ├── modal.handlebars
    └── table.handlebars
```

### 2. Helper 组织

```javascript
// helpers/index.js
const dateHelpers = require('./date');
const stringHelpers = require('./string');
const mathHelpers = require('./math');

module.exports = {
  ...dateHelpers,
  ...stringHelpers,
  ...mathHelpers
};

// helpers/date.js
const moment = require('moment');

module.exports = {
  formatDate: function (date, format = 'YYYY-MM-DD') {
    return moment(date).format(format);
  },
  
  timeAgo: function (date) {
    return moment(date).fromNow();
  }
};

// helpers/string.js
module.exports = {
  truncate: function (str, length = 100) {
    if (str.length > length) {
      return str.substring(0, length) + '...';
    }
    return str;
  },
  
  capitalize: function (str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
  }
};
```

### 3. 错误处理

```javascript
// 全局错误处理
app.use((err, req, res, next) => {
  console.error(err.stack);
  
  res.status(500).render('error', {
    layout: 'error',
    title: '服务器错误',
    message: process.env.NODE_ENV === 'production' 
      ? '服务器内部错误' 
      : err.message,
    error: process.env.NODE_ENV === 'production' ? {} : err
  });
});

// 404 处理
app.use((req, res) => {
  res.status(404).render('404', {
    title: '页面未找到',
    url: req.originalUrl
  });
});
```

### 4. 安全性考虑

```javascript
// 注册安全的 Helper
handlebars.registerHelper('safeUrl', function (url) {
  // 验证 URL 安全性
  if (url && typeof url === 'string') {
    // 简单的 URL 验证
    const urlRegex = /^https?:\/\//;
    if (urlRegex.test(url)) {
      return url;
    }
  }
  return '#';
});

// 转义用户输入
handlebars.registerHelper('escape', function (str) {
  return handlebars.Utils.escapeExpression(str);
});
```

## 总结

Handlebars 是一个功能强大且易于使用的模板引擎，具有以下特点：

1. **语法简洁**：使用双花括号语法，易于学习和使用
2. **功能丰富**：支持条件语句、循环、部分模板等
3. **可扩展性**：支持自定义 Helper 和部分模板
4. **性能优秀**：支持预编译和缓存
5. **安全性**：默认转义 HTML，防止 XSS 攻击
6. **生态完善**：与 Express 等框架良好集成

在实际项目中，合理使用 Handlebars 的各种特性，可以构建出结构清晰、易于维护的模板系统。
