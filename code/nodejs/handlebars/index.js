const express = require('express');
const exphbs = require('express-handlebars');
const path = require('path');

const hbs = exphbs.create({});

const app = express();
app.engine('handlebars',hbs.engine);
app.set('view engine','handlebars');
app.set('views',path.join(__dirname,'views'));

const tree = [
  {
    name:'parent1',
    children:[
      {
        name:'parent1-child1',
        children:[
          {name:'parent1-child1-grandson1',children:[]}
        ]
      },
      {name:'parent1-child2',children:[]}
    ]
  }
]

app.get('/',(req,res)=>{
  res.render('home',{title:'home'});
});
app.get('/tree',(req,res)=>{
  res.render('tree',{tree});
});

app.listen(3000, () => {
  console.log('Server is running on port 3000');
});
