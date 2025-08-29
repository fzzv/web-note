// 类装饰器 是一个函数，接收一个参数，参数是类的构造函数，可以修改类的构造函数，或者返回一个新的构造函数
function logClass(constructor: Function) {
  console.log(`Class Created`,constructor.name)
}

@logClass
class HttpClient {
  constructor() {
    console.log('HttpClient constructor');
  }
}

new HttpClient();

// 类装饰器工厂 是一个返回装饰器的函数，可以接收参数来控制装饰器的行为
function logClassWithParams(params: string) {
  return function (constructor: Function) {
    console.log(`Class Created`,constructor.name, params)
  }
}

@logClassWithParams('https://www.xyu.fan')
class HttpClient2 {
  constructor() {
    console.log('HttpClient2 constructor');
  }
}

new HttpClient2();

// 类装饰器扩展类的功能，比如说可以添加新的属性和方法，或者修改类的行为
function addProperty(target: any) {
  target.prototype.url = 'https://www.xyu.fan';
}

interface HttpClient3 {
  url: string;
}

@addProperty
class HttpClient3 {
  constructor() {
    console.log('HttpClient3 constructor');
  }
}

const httpClient3 = new HttpClient3();
console.log(httpClient3.url);

// 类装饰器重写类 可能通过返回一个新的构造函数来替换原有的构造函数
function rewriteClass<T extends { new (...args: any[]) }>(constructor: T) {
  return class extends constructor {
    constructor(...args: any[]) {
      super(...args);
      console.log('rewriteClass 4 constructor');
    }
  }
}

@rewriteClass
class HttpClient4 {
  constructor() {
    console.log('HttpClient4 constructor');
  }
}

new HttpClient4();
