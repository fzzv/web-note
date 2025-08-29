import 'reflect-metadata'

class PersonClass {
  private property: string;

  constructor(property: string) {
    this.property = property;
  }

  @Reflect.metadata('methodName', 'methodValue')
  method() {
    console.log('method');
  }
}

/**
 * @Reflect.metadata 和 Reflect.defineMetadata
 * 都是用来定义元数据，@Reflect.metadata 就是一个语法糖，可以简化元数据定义过程。
 */

const person = new PersonClass('property');

// 给 person 的 property 定义一个元数据
Reflect.defineMetadata('name', 'fan', person, 'property');

// 判断是否有元数据
console.log('判断是否有元数据:', Reflect.hasMetadata('name', person, 'property'));
// 获取 person 的 property 的元数据
console.log('获取 person 的 property 的元数据:', Reflect.getMetadata('name', person, 'property'));

// 获取自有元数据
console.log('获取自有元数据:', Reflect.getOwnMetadata('name', person, 'property'));
console.log('获取自有元数据 method:', Reflect.getOwnMetadata('methodName', Reflect.getPrototypeOf(person), 'method'));
console.log('获取自有元数据 method:', Reflect.getOwnMetadata('methodName', PersonClass.prototype, 'method'));

// 删除元数据
Reflect.deleteMetadata('name', person, 'property');
console.log('删除元数据后判断是否有元数据:', Reflect.hasMetadata('name', person, 'property'));
