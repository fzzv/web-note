# HTMX

HTMX 是一个现代化的 JavaScript 库，允许你直接在 HTML 中使用 AJAX、CSS 过渡、WebSocket 和服务器发送事件，无需编写 JavaScript 代码。它通过 HTML 属性扩展了 HTML 的功能。

## 基本概念

HTMX 的核心思想是：
- **HTML 优先**：通过 HTML 属性控制行为
- **服务器驱动**：服务器返回 HTML 片段而非 JSON
- **渐进增强**：在现有 HTML 基础上添加交互性
- **简化开发**：减少前端 JavaScript 代码量

## 使用

### HTML 引入

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HTMX Demo</title>
    <script src="https://unpkg.com/htmx.org@2.0.2"></script>
</head>
<body>
    <!-- HTMX 内容 -->
</body>
</html>
```

## 基础使用

### 1. 发送 GET 请求

```html
<!-- 前端 HTML -->
<button type="button" hx-get="/get">发送Get请求</button>
```

```javascript
// 后端路由
app.get('/get', (req, res) => {
  res.send('get响应');
});
```

**特点**：
- `hx-get`：指定 GET 请求的 URL
- 点击按钮时自动发送请求
- 响应内容会替换按钮的 innerHTML

### 2. 发送 POST 请求

```html
<!-- 前端 HTML -->
<button type="button" hx-post="/post">发送Post请求</button>
```

```javascript
// 后端路由
app.post('/post', (req, res) => {
  res.send('post响应');
});
```

### 3. 表单数据提交

```html
<form>
  <input type="text" name="name" placeholder="Name">
  <input type="email" name="email" placeholder="Email">
  <button hx-post="/payload" hx-trigger="click" hx-target="#result">Send</button>
</form>
<div id="result"></div>
```

```javascript
// 后端处理
app.post('/payload', (req, res) => {
  const user = req.body;
  user.id = Date.now();
  res.send(`
    <p>用户名:${user.name}</p>
    <p>邮箱:${user.email}</p>
  `);
});
```

**特点**：
- `hx-target="#result"`：指定响应内容插入的目标元素
- 自动收集表单数据并发送
- 服务器返回 HTML 片段

## 核心属性详解

### 1. HTTP 方法属性

```html
<!-- 不同的 HTTP 方法 -->
<button hx-get="/api/users">GET 请求</button>
<button hx-post="/api/users">POST 请求</button>
<button hx-put="/api/users/1">PUT 请求</button>
<button hx-patch="/api/users/1">PATCH 请求</button>
<button hx-delete="/api/users/1">DELETE 请求</button>
```

### 2. 触发器 (hx-trigger)

```html
<!-- 不同的触发事件 -->
<input hx-get="/search" hx-trigger="keyup">           <!-- 键盘释放时 -->
<button hx-get="/data" hx-trigger="click">            <!-- 点击时（默认） -->
<div hx-get="/time" hx-trigger="every 1s">           <!-- 每秒轮询 -->
<div hx-get="/load_polling" hx-trigger="load delay:1s"> <!-- 加载后延迟1秒 -->
```

#### 轮询示例

```html
<div hx-get="/time" hx-trigger="every 1s"></div>
<div hx-get="/load_polling" hx-trigger="load delay:1s" hx-swap="outerHTML"></div>
```

```javascript
// 后端实现
app.get('/time', (req, res) => {
  res.send(new Date().toLocaleString());
});

let count = 1;
app.get('/load_polling', (req, res) => {
  if (count++ < 10) {
    res.send(`<div hx-get="/load_polling" hx-trigger="load delay:1s" hx-swap="outerHTML">已加载${count}0%</div>`);
  } else {
    res.send('加载完成');
  }
});
```

### 3. 目标选择器 (hx-target)

```html
<!-- 指定不同的目标元素 -->
<button hx-get="/data" hx-target="#result">更新结果区域</button>
<button hx-get="/data" hx-target="body">更新整个页面</button>
<button hx-get="/data" hx-target="closest div">更新最近的 div</button>
<button hx-get="/data" hx-target="next .content">更新下一个 content 元素</button>

<div id="result"></div>
```

### 4. 内容替换方式 (hx-swap)

```html
<!-- 不同的替换方式 -->
<button hx-get="/time" hx-target="#target" hx-swap="beforebegin">beforebegin-变成target的哥哥</button>
<button hx-get="/time" hx-target="#target" hx-swap="afterbegin">afterbegin-变成target第一子元素之前</button>
<button hx-get="/time" hx-target="#target" hx-swap="beforeend">beforeend-变成target最后一个子元素</button>
<button hx-get="/time" hx-target="#target" hx-swap="afterend">afterend-变成target的弟弟</button>
<button hx-get="/time" hx-target="#target" hx-swap="delete">delete-删除目标元素</button>
<button hx-get="/time" hx-target="#target" hx-swap="none">none-不替换</button>
<button hx-get="/time" hx-target="#target" hx-swap="innerHTML">innerHTML-替换内容</button>
<button hx-get="/time" hx-target="#target" hx-swap="outerHTML">outerHTML-替换元素本身</button>

<div id="target">
  <div>子元素1</div>
  <div>子元素2</div>
</div>
```

**替换方式说明**：
- `innerHTML`：替换目标元素的内容（默认）
- `outerHTML`：替换目标元素本身
- `beforebegin`：在目标元素之前插入
- `afterbegin`：在目标元素内部的开始位置插入
- `beforeend`：在目标元素内部的结束位置插入
- `afterend`：在目标元素之后插入
- `delete`：删除目标元素
- `none`：不进行任何替换

## 高级特性

### 1. 越界替换 (Out of Band Swaps)

```html
<!-- 前端 HTML -->
<button hx-get="/oob" hx-target="#mainTarget">获取内容</button>
<div id="mainTarget">主要目标</div>
<div id="otherTarget">其他目标</div>
```

```javascript
// 后端返回多个目标的内容
app.get('/oob', (req, res) => {
  res.send(`
    <span>这是我返回的主要内容1</span>
    <span>这是我返回的主要内容2</span>
    <div id="otherTarget" hx-swap-oob="true">其它目标内容</div>
  `);
});
```

**特点**：
- `hx-swap-oob="true"`：允许同时更新多个目标元素
- 主要内容更新指定目标，带有 `hx-swap-oob` 的内容更新对应 ID 的元素

### 2. 文件上传

```html
<form hx-post="/upload" hx-encoding="multipart/form-data" hx-target="#result">
  <input type="file" name="file">
  <button type="submit">上传</button>
</form>
<div id="result"></div>
```

```javascript
// 后端文件上传处理
const multer = require('multer');
const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    cb(null, path.join(__dirname, 'uploads'));
  },
  filename: (req, file, cb) => {
    cb(null, file.fieldname + '-' + Date.now() + path.extname(file.originalname));
  }
});
const upload = multer({ storage });

app.post('/upload', upload.single('file'), (req, res) => {
  const filePath = req.file.path;
  console.log('上传成功', filePath);
  res.send(`<b>上传成功</b>:${filePath}`);
});
```

### 3. 请求指示器

```html
<style>
.htmx-request .spinner {
  display: inline-block;
}
.spinner {
  display: none;
}
</style>

<button hx-get="/delay" hx-target="#result">
  发送请求
  <span class="spinner">⏳ 加载中...</span>
</button>
<div id="result"></div>
```

```javascript
// 模拟延迟响应
app.get('/delay', (req, res) => {
  setTimeout(() => {
    res.send('延迟响应完成');
  }, 3000);
});
```

**特点**：
- HTMX 在请求期间自动添加 `htmx-request` 类
- 可以通过 CSS 显示加载指示器

### 4. 请求取消

```html
<button hx-get="/firstRequest" hx-target="#result">第一个请求</button>
<button hx-get="/secondRequest" hx-target="#result">第二个请求</button>
<div id="result"></div>
```

```javascript
app.get('/firstRequest', (req, res) => {
  setTimeout(() => {
    res.send('firstRequestResponse');
  }, 6000);
});

app.get('/secondRequest', (req, res) => {
  setTimeout(() => {
    res.send('secondRequestResponse');
  }, 3000);
});
```

**特点**：
- 当新请求发起时，HTMX 会自动取消同一目标的旧请求
- 防止竞态条件

## 事件系统

### 1. 请求生命周期事件

```html
<button hx-get="/time" hx-target="#target">获取时间</button>
<div id="target"></div>

<script>
// 请求前事件
htmx.on('htmx:beforeRequest', (event) => {
  console.log('htmx:beforeRequest');
  console.log(event.detail.requestConfig.headers);
});

// 请求后事件
htmx.on('htmx:afterRequest', (event) => {
  console.log('htmx:afterRequest');
  console.log(event.detail.xhr.response);
});

// 请求错误事件
htmx.on('htmx:responseError', (event) => {
  console.log('htmx:responseError');
  console.log(event.detail.xhr.response);
});
</script>
```

### 2. 常用事件列表

```javascript
// 请求相关事件
htmx.on('htmx:beforeRequest', handler);    // 请求发送前
htmx.on('htmx:afterRequest', handler);     // 请求完成后
htmx.on('htmx:sendError', handler);        // 发送错误
htmx.on('htmx:responseError', handler);    // 响应错误
htmx.on('htmx:timeout', handler);          // 请求超时

// 内容相关事件
htmx.on('htmx:beforeSwap', handler);       // 内容交换前
htmx.on('htmx:afterSwap', handler);        // 内容交换后
htmx.on('htmx:beforeSettle', handler);     // 内容稳定前
htmx.on('htmx:afterSettle', handler);      // 内容稳定后

// 其他事件
htmx.on('htmx:load', handler);             // 元素加载完成
htmx.on('htmx:configRequest', handler);    // 配置请求
```

### 3. 事件修改请求行为

```html
<button hx-get="/data" hx-target="#result">获取数据</button>
<div id="result"></div>

<script>
// 修改请求头
htmx.on('htmx:configRequest', (event) => {
  event.detail.headers['Authorization'] = 'Bearer token123';
  event.detail.headers['X-Custom-Header'] = 'custom-value';
});

// 修改交换行为
htmx.on('htmx:beforeSwap', (event) => {
  if (event.detail.xhr.status === 404) {
    event.detail.shouldSwap = true;
    event.detail.target.innerHTML = '<p>内容未找到</p>';
  }
});
</script>
```

## 实际应用示例

### 1. 动态搜索

```html
<input type="text" 
       name="search" 
       hx-get="/search" 
       hx-trigger="keyup changed delay:300ms" 
       hx-target="#search-results"
       placeholder="搜索...">
<div id="search-results"></div>
```

```javascript
app.get('/search', (req, res) => {
  const keyword = req.query.search || '';
  if (keyword.length === 0) {
    res.send('');
    return;
  }
  
  // 模拟搜索结果
  const results = [
    `${keyword}的搜索结果1`,
    `${keyword}的搜索结果2`,
    `${keyword}的搜索结果3`
  ];
  
  res.send(`
    <ul>
      ${results.map(result => `<li>${result}</li>`).join('')}
    </ul>
  `);
});
```

### 2. 表单验证

```html
<form hx-post="/validate" hx-target="#validation-result">
  <input type="email" name="email" 
         hx-post="/validate-email" 
         hx-trigger="blur" 
         hx-target="#email-validation">
  <div id="email-validation"></div>
  
  <input type="password" name="password" 
         hx-post="/validate-password" 
         hx-trigger="blur" 
         hx-target="#password-validation">
  <div id="password-validation"></div>
  
  <button type="submit">提交</button>
</form>
<div id="validation-result"></div>
```

```javascript
app.post('/validate-email', (req, res) => {
  const email = req.body.email;
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  
  if (!email) {
    res.send('<span style="color: red;">邮箱不能为空</span>');
  } else if (!emailRegex.test(email)) {
    res.send('<span style="color: red;">邮箱格式不正确</span>');
  } else {
    res.send('<span style="color: green;">邮箱格式正确</span>');
  }
});

app.post('/validate-password', (req, res) => {
  const password = req.body.password;
  
  if (!password) {
    res.send('<span style="color: red;">密码不能为空</span>');
  } else if (password.length < 6) {
    res.send('<span style="color: red;">密码至少6位</span>');
  } else {
    res.send('<span style="color: green;">密码强度良好</span>');
  }
});
```

### 3. 购物车功能

```html
<div id="cart">
  <h3>购物车</h3>
  <div id="cart-items"></div>
  <div id="cart-total">总计: $0.00</div>
</div>

<div class="products">
  <div class="product">
    <h4>商品 1</h4>
    <p>价格: $10.00</p>
    <button hx-post="/cart/add" 
            hx-vals='{"product": "product1", "price": 10.00}'
            hx-target="#cart"
            hx-swap="outerHTML">
      加入购物车
    </button>
  </div>
  
  <div class="product">
    <h4>商品 2</h4>
    <p>价格: $15.00</p>
    <button hx-post="/cart/add" 
            hx-vals='{"product": "product2", "price": 15.00}'
            hx-target="#cart"
            hx-swap="outerHTML">
      加入购物车
    </button>
  </div>
</div>
```

```javascript
let cart = [];

app.post('/cart/add', (req, res) => {
  const { product, price } = req.body;
  
  // 查找是否已存在
  const existingItem = cart.find(item => item.product === product);
  if (existingItem) {
    existingItem.quantity += 1;
  } else {
    cart.push({ product, price: parseFloat(price), quantity: 1 });
  }
  
  // 计算总计
  const total = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  
  // 返回更新的购物车
  res.send(`
    <div id="cart">
      <h3>购物车</h3>
      <div id="cart-items">
        ${cart.map(item => `
          <div class="cart-item">
            ${item.product} x ${item.quantity} = $${(item.price * item.quantity).toFixed(2)}
            <button hx-post="/cart/remove" 
                    hx-vals='{"product": "${item.product}"}'
                    hx-target="#cart"
                    hx-swap="outerHTML">删除</button>
          </div>
        `).join('')}
      </div>
      <div id="cart-total">总计: $${total.toFixed(2)}</div>
    </div>
  `);
});

app.post('/cart/remove', (req, res) => {
  const { product } = req.body;
  cart = cart.filter(item => item.product !== product);
  
  const total = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  
  res.send(`
    <div id="cart">
      <h3>购物车</h3>
      <div id="cart-items">
        ${cart.map(item => `
          <div class="cart-item">
            ${item.product} x ${item.quantity} = $${(item.price * item.quantity).toFixed(2)}
            <button hx-post="/cart/remove" 
                    hx-vals='{"product": "${item.product}"}'
                    hx-target="#cart"
                    hx-swap="outerHTML">删除</button>
          </div>
        `).join('')}
      </div>
      <div id="cart-total">总计: $${total.toFixed(2)}</div>
    </div>
  `);
});
```

### 4. 无限滚动

```html
<div id="content">
  <!-- 初始内容 -->
</div>
<div hx-get="/load-more?page=2" 
     hx-trigger="revealed" 
     hx-swap="outerHTML">
  <p>加载更多...</p>
</div>
```

```javascript
app.get('/load-more', (req, res) => {
  const page = parseInt(req.query.page) || 1;
  const hasMore = page < 10; // 假设总共10页
  
  // 生成内容
  const content = Array.from({ length: 5 }, (_, i) => 
    `<p>第${page}页的内容 ${i + 1}</p>`
  ).join('');
  
  if (hasMore) {
    res.send(`
      ${content}
      <div hx-get="/load-more?page=${page + 1}" 
           hx-trigger="revealed" 
           hx-swap="outerHTML">
        <p>加载更多...</p>
      </div>
    `);
  } else {
    res.send(`${content}<p>没有更多内容了</p>`);
  }
});
```

## 性能优化

### 1. 请求去抖动

```html
<!-- 使用 delay 修饰符减少请求频率 -->
<input hx-get="/search" 
       hx-trigger="keyup changed delay:300ms" 
       hx-target="#results">
```

### 2. 请求缓存

```html
<!-- 缓存 GET 请求结果 -->
<button hx-get="/static-data" hx-target="#result" hx-cache="true">
  获取静态数据
</button>
```

### 3. 选择性更新

```html
<!-- 使用 hx-select 只更新响应的特定部分 -->
<button hx-get="/full-page" 
        hx-target="#result" 
        hx-select="#content">
  只获取内容部分
</button>
```

```javascript
app.get('/full-page', (req, res) => {
  res.send(`
    <html>
      <head><title>完整页面</title></head>
      <body>
        <div id="header">头部</div>
        <div id="content">这是需要的内容</div>
        <div id="footer">底部</div>
      </body>
    </html>
  `);
});
```

## 与其他技术集成

### 1. CSS 动画

```html
<style>
.fade-in {
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}
.fade-in.htmx-added {
  opacity: 1;
}
</style>

<button hx-get="/animated-content" hx-target="#result">
  获取动画内容
</button>
<div id="result"></div>
```

```javascript
app.get('/animated-content', (req, res) => {
  res.send('<div class="fade-in">这是带动画的内容</div>');
});
```

### 2. WebSocket 集成

```html
<div hx-ws="connect:/ws">
  <form hx-ws="send:submit">
    <input name="message" placeholder="输入消息">
    <button type="submit">发送</button>
  </form>
  <div id="messages"></div>
</div>
```

### 3. 服务器发送事件 (SSE)

```html
<div hx-sse="connect:/events">
  <div hx-sse="swap:message" hx-target="#notifications"></div>
</div>
<div id="notifications"></div>
```

```javascript
app.get('/events', (req, res) => {
  res.writeHead(200, {
    'Content-Type': 'text/event-stream',
    'Cache-Control': 'no-cache',
    'Connection': 'keep-alive'
  });
  
  const interval = setInterval(() => {
    res.write(`data: <p>消息时间: ${new Date().toLocaleString()}</p>\n\n`);
  }, 1000);
  
  req.on('close', () => {
    clearInterval(interval);
  });
});
```

## 调试和开发工具

### 1. 启用日志

```html
<script>
// 启用 HTMX 调试日志
htmx.logAll();
</script>
```

### 2. 请求拦截器

```html
<script>
// 请求拦截器，用于调试
htmx.on('htmx:configRequest', (event) => {
  console.log('请求配置:', event.detail);
  
  // 添加调试信息
  event.detail.headers['X-Debug'] = 'true';
  event.detail.headers['X-Timestamp'] = Date.now();
});

// 响应拦截器
htmx.on('htmx:beforeSwap', (event) => {
  console.log('响应内容:', event.detail.xhr.response);
  console.log('目标元素:', event.detail.target);
});
</script>
```

### 3. 错误处理

```html
<script>
// 全局错误处理
htmx.on('htmx:responseError', (event) => {
  const status = event.detail.xhr.status;
  const response = event.detail.xhr.response;
  
  if (status === 404) {
    event.detail.target.innerHTML = '<p>页面未找到</p>';
  } else if (status === 500) {
    event.detail.target.innerHTML = '<p>服务器错误，请稍后重试</p>';
  } else {
    event.detail.target.innerHTML = `<p>请求失败: ${status}</p>`;
  }
  
  // 阻止默认错误处理
  event.detail.shouldSwap = true;
});

// 网络错误处理
htmx.on('htmx:sendError', (event) => {
  event.detail.target.innerHTML = '<p>网络连接错误</p>';
});
</script>
```

## 最佳实践

### 1. HTML 结构设计

```html
<!-- 良好的结构设计 -->
<div class="user-profile" id="user-profile">
  <div class="user-info" id="user-info">
    <!-- 用户信息 -->
  </div>
  
  <div class="user-actions">
    <button hx-get="/user/edit" hx-target="#user-info">编辑</button>
    <button hx-delete="/user" hx-target="#user-profile" 
            hx-confirm="确定要删除用户吗？">删除</button>
  </div>
</div>
```

### 2. 服务器响应设计

```javascript
// 返回完整的 HTML 片段
app.get('/user/edit', (req, res) => {
  res.send(`
    <form hx-put="/user" hx-target="#user-info">
      <input type="text" name="name" value="当前用户名">
      <input type="email" name="email" value="current@email.com">
      <button type="submit">保存</button>
      <button type="button" hx-get="/user/view" hx-target="#user-info">取消</button>
    </form>
  `);
});

app.put('/user', (req, res) => {
  // 更新用户信息
  const user = updateUser(req.body);
  
  res.send(`
    <div class="user-display">
      <h3>${user.name}</h3>
      <p>${user.email}</p>
    </div>
  `);
});
```

### 3. 渐进增强

```html
<!-- 确保没有 JavaScript 时也能工作 -->
<form action="/user" method="POST" hx-post="/user" hx-target="#result">
  <input type="text" name="name" required>
  <button type="submit">提交</button>
</form>

<!-- 提供加载状态 -->
<button hx-get="/slow-endpoint" hx-target="#result" hx-indicator="#spinner">
  获取数据
  <span id="spinner" class="htmx-indicator">⏳</span>
</button>
```

### 4. 安全考虑

```javascript
// 服务器端验证和转义
app.post('/comment', (req, res) => {
  const comment = escapeHtml(req.body.comment);
  const userId = req.session.userId; // 验证用户身份
  
  if (!userId) {
    return res.status(401).send('<p>请先登录</p>');
  }
  
  // 保存评论...
  res.send(`<div class="comment">${comment}</div>`);
});

function escapeHtml(text) {
  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return text.replace(/[&<>"']/g, m => map[m]);
}
```

## 总结

HTMX 是一个强大而简洁的库，具有以下特点：

1. **简化开发**：通过 HTML 属性控制 AJAX 行为，减少 JavaScript 代码
2. **服务器驱动**：服务器返回 HTML 片段，简化前后端交互
3. **渐进增强**：在现有 HTML 基础上添加交互性
4. **丰富功能**：支持多种 HTTP 方法、触发器、替换方式
5. **事件系统**：完整的请求生命周期事件
6. **高性能**：内置请求取消、缓存等优化功能

HTMX 特别适合：
- 需要快速开发的项目
- 服务器端渲染的应用
- 不需要复杂前端框架的场景
- 渐进式增强现有应用

通过合理使用 HTMX 的各种特性，可以构建出交互丰富且性能良好的 Web 应用。
