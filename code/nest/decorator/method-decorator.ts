// 方法装饰器 是一个函数，接收三个参数，分别是 target、propertyKey、descriptor
// target 是类的原型
// propertyKey 是方法的名称
// descriptor 是方法的描述符
function logMethod(target: Object, propertyKey: string | Symbol, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value;
  descriptor.value = function (...args: any[]) {
    console.log(`Method ${propertyKey} called with arguments:`, args);
    const result = originalMethod.apply(this, args);
    console.log(`Method ${propertyKey} returned:`, result);
    return result;
  }
}

class Calculator {
  @logMethod
  add(a: number, b: number) {
    return a + b;
  }
}

const calculator = new Calculator();
calculator.add(1, 2);

// 可以用作方法前，判断是否有权限执行该方法
function checkPermission(target: Object, propertyKey: string | Symbol, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value;
  descriptor.value = function (...args: any[]) {
    if (user[args[0]].role.includes('admin')) {
      return originalMethod.apply(this, args);
    } else {
      throw new Error('Permission denied');
    }
  }
  return descriptor;
}
const user = {
  'user1': {
    role: ['admin', 'user'],
  },
  'user2': {
    role: ['user'],
  },
}

class WebSite {
  @checkPermission
  deleteUser(userId: string) {
    console.log(`delete user ${userId}`);
  }
}

const webSite = new WebSite();
webSite.deleteUser('user1');
// webSite.deleteUser('user2');

// 方法装饰器可以实现方法结果的缓存
function cacheMethod(target: Object, propertyKey: string | Symbol, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value;
  const cache = new Map();
  descriptor.value = function (...args: any[]) {
    const key = JSON.stringify(args);
    if (cache.has(key)) {
      return cache.get(key);
    }
    const result = originalMethod.apply(this, args);
    console.log(`cache ${key} = ${result}`);
    cache.set(key, result);
    return result;
  }
  return descriptor;
}

class CacheCalculator {
  @cacheMethod
  add(a: number, b: number) {
    return a + b;
  }
}

const cacheCalculator = new CacheCalculator();
cacheCalculator.add(1, 2);
cacheCalculator.add(1, 2);

// 方法装饰器可以实现方法的性能监控
function performanceMethod(target: Object, propertyKey: string | Symbol, descriptor: PropertyDescriptor) {
  const originalMethod = descriptor.value;
  descriptor.value = function (...args: any[]) {
    const start = performance.now();
    const result = originalMethod.apply(this, args);
    const end = performance.now();
    console.log(`Method ${propertyKey} executed in ${end - start} milliseconds`);
    return result;
  };
  return descriptor;
}

class DataProcessor {
  @performanceMethod
  processLargeData(data: any[]) {
    // 模拟耗时操作
    return data.map(item => item * 2);
  }
}

new DataProcessor().processLargeData([1, 2, 3, 4, 5]);
