import 'reflect-metadata';

// 访问器装饰器会是用来装饰访问器的 get set
function LogAccessor(
  target: any,
  propertyKey: string,
  descriptor: PropertyDescriptor
) {
  const originalGet = descriptor.get;
  descriptor.get = function () {
    console.log(`Accessed property: ${propertyKey}`);
    return originalGet?.call(this);
  };
}

class User {
  private _name = 'Alice';

  @LogAccessor
  get name() {
    return this._name;
  }
}

const u = new User();
console.log(u.name); // 访问时会触发装饰器逻辑

// 属性装饰器是用来装饰属性的
function Required(target: any, propertyKey: string) {
  const requiredProps = Reflect.getMetadata('required:props', target) || [];
  Reflect.defineMetadata(
    'required:props',
    [...requiredProps, propertyKey],
    target
  );
}

class User1 {
  @Required
  username: string;

  email: string;
}

function validateRequired(obj: any) {
  const requiredProps = Reflect.getMetadata('required:props', obj) || [];
  for (const key of requiredProps) {
    if (!obj[key]) {
      console.log(`${key} is required`);
    }
  }
}

const u1 = new User1();
u1.username = ''; // 未赋值
validateRequired(u1); // 输出：username is required


// 不同类型的装饰器的执行顺序
/**
 * 1.属性装饰器、方法装饰器、访问器装饰器它们是按照在类中出现的顺序，从上往下依次执行
 * 2.类装饰器最后执行
 * 3.参数装饰器先于方法执行
 */
function classDecorator1(target) {
  console.log('classDecorator1')
}
function classDecorator2(target) {
  console.log('classDecorator2')
}
function propertyDecorator1(target, propertyKey) {
  console.log('propertyDecorator1')
}
function propertyDecorator2(target, propertyKey) {
  console.log('propertyDecorator2')
}
function methodDecorator1(target, propertyKey) {
  console.log('methodDecorator1')
}
function methodDecorator2(target, propertyKey) {
  console.log('methodDecorator2')
}
function accessorDecorator1(target, propertyKey) {
  console.log('accessorDecorator1')
}
function accessorDecorator2(target, propertyKey) {
  console.log('accessorDecorator2')
}
function parametorDecorator4(target, propertyKey, parametorIndex: number) {
  console.log('parametorDecorator4', propertyKey)//propertyKey方法名
}
function parametorDecorator3(target, propertyKey, parametorIndex: number) {
  console.log('parametorDecorator3', propertyKey)//propertyKey方法名
}
function parametorDecorator2(target, propertyKey, parametorIndex: number) {
  console.log('parametorDecorator2', propertyKey)//propertyKey方法名
}
function parametorDecorator1(target, propertyKey, parametorIndex: number) {
  console.log('parametorDecorator1', propertyKey)//propertyKey方法名
}
@classDecorator1
@classDecorator2
class Example {
  @accessorDecorator1
  @accessorDecorator2
  get myProp(){
      return this.prop;
  }

  @propertyDecorator1
  @propertyDecorator2
  prop:string

  @methodDecorator1
  @methodDecorator2
  method(@parametorDecorator4 @parametorDecorator3 param1:any,@parametorDecorator2 @parametorDecorator1 param2:any) {}
}
// 如果一个方法有多个参数，参数装饰器会从右向左执行
// 一个参数也可有会有多个参数装饰 器，这些装饰 器也是从右向左执行的
