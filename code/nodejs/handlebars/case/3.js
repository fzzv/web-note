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
`

const data = {
  items: ['apple', 'banana', 'cherry', 'orange']
}

const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
