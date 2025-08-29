const handlebars = require('handlebars');

// if else 语句
const tpl = `
  {{#if isAdmin}}
    <h1>Admin</h1>
  {{else}}
    <h1>User</h1>
  {{/if}}
`

const data = {
  isAdmin: false
}

const template = handlebars.compile(tpl);
const result = template(data);
console.log(result);
