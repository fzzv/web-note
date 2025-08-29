const handlebars = require('handlebars');

// 1. 创建一个模板
const tpl = `
  <h1>{{title}}</h1>
  <p>{{content}}</p>
`

// 2. 上下文对象
const data = {
  title: 'Hello, Handlebars!',
  content: 'This is a template example.'
};

// 3. 渲染模板
const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
