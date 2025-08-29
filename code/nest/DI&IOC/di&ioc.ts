// 在没有使用IOC的时候，我们需要直接负责创建依赖的对象
class Engine {
  start() {
    console.log('Engine started');
  }
}

class Car {
  private engine: Engine;
  constructor() {
    this.engine = new Engine();
  }

  start() {
    this.engine.start();
  }
}

const car = new Car();
car.start();

// 使用构造函数注入后，我们不需要直接负责创建依赖的对象
class Car2 {
  private engine: Engine;
  constructor(engine: Engine) {
    this.engine = engine;
  }

  start() {
    this.engine.start();
  }
}

const car2 = new Car2(new Engine());
car2.start();

/**
 * IOC 控制反转 Inversion Of Control
 * 使用IOC之后，对象的创建和管理权被反转给了容器，程序不再主动的负责创建对象，而是被动的接收容器注入和对象
 * DI 依赖注入 Depedency Injection 是实现IOC的一种手段，通过DI，我们将类的依赖项注入到类中，而不是在类里面实例化这些依赖
 */
import 'reflect-metadata';
function Injectable(target) {
    //这里面可以不用写任何代码，此装饰器不需要执行任何操作，仅仅用于元数据的生成
}
//@Injectable
class Oil {
  constructor(private num: number) {

  }
}

const dependencies = Reflect.getMetadata('design:paramtypes', Oil)
console.log(dependencies)
