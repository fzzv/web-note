# SOLID 原则

SOLID 原则是面向对象设计的五个基本原则，由 Robert C. Martin 提出，旨在使软件设计更加灵活、可维护和可扩展。

## 1. 单一职责原则（Single Responsibility Principle, SRP）

### 定义
一个类应该只有一个引起变化的原因，换句话说，就是一个类应该只有一个职责。

### 违反原则的示例

```typescript
// ❌ 违反单一职责原则
class User {
  constructor(public name: string, public email: string) {}
  
  save() {
    // 数据库操作职责
  }
  
  sendEmail() {
    // 邮件发送职责
  }
}
```

**问题分析**：
- User 类承担了多个职责：用户数据管理、数据持久化、邮件发送
- 当数据库逻辑改变时，User 类需要修改
- 当邮件发送逻辑改变时，User 类也需要修改
- 违反了单一职责原则

### 符合原则的示例

```typescript
// ✅ 符合单一职责原则
class User {
  constructor(public name: string, public email: string) {
    // 只负责用户数据的表示
  }
}

class UserRepository {
  save(user: User) {
    // 只负责用户数据的持久化
    console.log('保存用户到数据库');
  }
}

class EmailService {
  sendEmail(user: User) {
    // 只负责邮件发送
    console.log(`向 ${user.email} 发送邮件`);
  }
}
```

### 在 NestJS 中的应用

```typescript
// 用户实体 - 只负责数据表示
@Entity()
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string;

  @Column()
  email: string;
}

// 用户仓储 - 只负责数据访问
@Injectable()
export class UserRepository {
  constructor(
    @InjectRepository(User)
    private repository: Repository<User>
  ) {}

  async save(user: User): Promise<User> {
    return this.repository.save(user);
  }

  async findById(id: number): Promise<User> {
    return this.repository.findOne({ where: { id } });
  }
}

// 邮件服务 - 只负责邮件发送
@Injectable()
export class EmailService {
  async sendWelcomeEmail(user: User): Promise<void> {
    // 发送欢迎邮件的逻辑
  }
}

// 用户服务 - 协调各个组件
@Injectable()
export class UserService {
  constructor(
    private userRepository: UserRepository,
    private emailService: EmailService
  ) {}

  async createUser(userData: CreateUserDto): Promise<User> {
    const user = new User();
    user.name = userData.name;
    user.email = userData.email;
    
    const savedUser = await this.userRepository.save(user);
    await this.emailService.sendWelcomeEmail(savedUser);
    
    return savedUser;
  }
}
```

## 2. 开闭原则（Open-Closed Principle, OCP）

### 定义
软件实体应该对扩展开放，对修改关闭。应该通过扩展已有代码来实现新功能，而不是修改原有代码。

### 违反原则的示例

```typescript
// ❌ 违反开闭原则
class Rectangle {
  constructor(public width: number, public height: number) {}
}

class Circle {
  constructor(public radius: number) {}
}

class AreaCalculator {
  calculateArea(shape: any) {
    if (shape instanceof Rectangle) {
      return shape.width * shape.height;
    } else if (shape instanceof Circle) {
      return Math.PI * shape.radius * shape.radius;
    }
    // 每次添加新形状都需要修改这个方法
  }
}
```

**问题分析**：
- 每次添加新的形状类型，都需要修改 `AreaCalculator` 类
- 违反了开闭原则，对修改不够封闭

### 符合原则的示例

```typescript
// ✅ 符合开闭原则
interface Shape {
  calculateArea(): number;
}

class Rectangle implements Shape {
  constructor(public width: number, public height: number) {}
  
  calculateArea(): number {
    return this.width * this.height;
  }
}

class Circle implements Shape {
  constructor(public radius: number) {}
  
  calculateArea(): number {
    return Math.PI * this.radius * this.radius;
  }
}

class AreaCalculator {
  static calculateArea(shape: Shape): number {
    return shape.calculateArea();
  }
}

// 使用示例
const result1 = AreaCalculator.calculateArea(new Rectangle(10, 20));
console.log(result1); // 200

const result2 = AreaCalculator.calculateArea(new Circle(10));
console.log(result2); // 314.159...

// 扩展新形状无需修改现有代码
class Triangle implements Shape {
  constructor(public base: number, public height: number) {}
  
  calculateArea(): number {
    return (this.base * this.height) / 2;
  }
}
```

### 在 NestJS 中的应用

```typescript
// 支付策略接口
interface PaymentStrategy {
  pay(amount: number): Promise<PaymentResult>;
}

// 具体支付实现
@Injectable()
export class CreditCardPayment implements PaymentStrategy {
  async pay(amount: number): Promise<PaymentResult> {
    // 信用卡支付逻辑
    return { success: true, transactionId: 'cc_123' };
  }
}

@Injectable()
export class PayPalPayment implements PaymentStrategy {
  async pay(amount: number): Promise<PaymentResult> {
    // PayPal 支付逻辑
    return { success: true, transactionId: 'pp_456' };
  }
}

// 支付服务 - 无需修改即可支持新的支付方式
@Injectable()
export class PaymentService {
  private strategies = new Map<string, PaymentStrategy>();

  constructor(
    private creditCardPayment: CreditCardPayment,
    private paypalPayment: PayPalPayment
  ) {
    this.strategies.set('credit-card', creditCardPayment);
    this.strategies.set('paypal', paypalPayment);
  }

  async processPayment(type: string, amount: number): Promise<PaymentResult> {
    const strategy = this.strategies.get(type);
    if (!strategy) {
      throw new Error(`Unsupported payment type: ${type}`);
    }
    return strategy.pay(amount);
  }
}
```

## 3. 里氏替换原则（Liskov Substitution Principle, LSP）

### 定义
子类必须能够替换掉它们的基类，这意味着子类应该在任何地方都能替换父类，并且不会导致程序出现异常。

### 违反原则的示例

```typescript
// ❌ 违反里氏替换原则
class Bird {
  fly() {
    console.log('鸟儿在飞翔');
  }
}

class Penguin extends Bird {
  fly() {
    throw new Error('企鹅不会飞');
  }
}

function moveBird(bird: Bird) {
  bird.fly(); // 如果传入 Penguin，会抛出异常
}

// moveBird(new Penguin()); // 💥 会抛出异常
```

**问题分析**：
- Penguin 不能完全替换 Bird，因为它改变了基类的行为预期
- 违反了里氏替换原则

### 符合原则的示例

```typescript
// ✅ 符合里氏替换原则
class Bird {
  move() {
    console.log('鸟儿在移动');
  }
}

class FlyingBird extends Bird {
  move() {
    console.log('我通过飞翔移动');
  }
}

class Penguin extends Bird {
  move() {
    console.log('我通过行走移动');
  }
}

function moveBird(bird: Bird) {
  bird.move(); // 任何子类都能正常工作
}

moveBird(new FlyingBird()); // ✅ 正常工作
moveBird(new Penguin());    // ✅ 正常工作
```

### 更好的设计

```typescript
// 更精细的接口设计
interface Movable {
  move(): void;
}

interface Flyable {
  fly(): void;
}

class Bird implements Movable {
  move(): void {
    console.log('鸟儿在移动');
  }
}

class Eagle extends Bird implements Flyable {
  move(): void {
    console.log('老鹰在移动');
  }
  
  fly(): void {
    console.log('老鹰在高空翱翔');
  }
}

class Penguin extends Bird {
  move(): void {
    console.log('企鹅在陆地上行走');
  }
  
  swim(): void {
    console.log('企鹅在水中游泳');
  }
}
```

### 在 NestJS 中的应用

```typescript
// 基础仓储接口
abstract class BaseRepository<T> {
  abstract save(entity: T): Promise<T>;
  abstract findById(id: string): Promise<T>;
  abstract delete(id: string): Promise<void>;
}

// 用户仓储实现
@Injectable()
export class UserRepository extends BaseRepository<User> {
  constructor(
    @InjectRepository(User)
    private repository: Repository<User>
  ) {
    super();
  }

  async save(user: User): Promise<User> {
    return this.repository.save(user);
  }

  async findById(id: string): Promise<User> {
    return this.repository.findOne({ where: { id } });
  }

  async delete(id: string): Promise<void> {
    await this.repository.delete(id);
  }
}

// 产品仓储实现
@Injectable()
export class ProductRepository extends BaseRepository<Product> {
  constructor(
    @InjectRepository(Product)
    private repository: Repository<Product>
  ) {
    super();
  }

  async save(product: Product): Promise<Product> {
    return this.repository.save(product);
  }

  async findById(id: string): Promise<Product> {
    return this.repository.findOne({ where: { id } });
  }

  async delete(id: string): Promise<void> {
    await this.repository.delete(id);
  }
}

// 服务可以使用任何符合 BaseRepository 的实现
@Injectable()
export class GenericService<T> {
  constructor(private repository: BaseRepository<T>) {}

  async createEntity(entity: T): Promise<T> {
    return this.repository.save(entity);
  }
}
```

## 4. 接口隔离原则（Interface Segregation Principle, ISP）

### 定义
类之间的依赖关系应该建立在最小的接口上，不应该强迫一个类依赖于它不使用的方法。

### 违反原则的示例

```typescript
// ❌ 违反接口隔离原则
interface Animal {
  eat(): void;
  fly(): void;
}

class Dog implements Animal {
  eat(): void {
    console.log('狗在吃东西');
  }
  
  fly(): void {
    throw new Error('狗不会飞'); // 被迫实现不需要的方法
  }
}
```

**问题分析**：
- Dog 类被迫实现它不需要的 `fly()` 方法
- 违反了接口隔离原则

### 符合原则的示例

```typescript
// ✅ 符合接口隔离原则
interface Eater {
  eat(): void;
}

interface Flyer {
  fly(): void;
}

class Dog implements Eater {
  eat(): void {
    console.log('狗在吃东西');
  }
}

class Bird implements Eater, Flyer {
  eat(): void {
    console.log('鸟在吃东西');
  }
  
  fly(): void {
    console.log('鸟在飞翔');
  }
}
```

### 更完善的示例

```typescript
// 细分的接口
interface Walkable {
  walk(): void;
}

interface Swimmable {
  swim(): void;
}

interface Flyable {
  fly(): void;
}

interface Eater {
  eat(): void;
}

// 不同动物实现不同的接口组合
class Fish implements Swimmable, Eater {
  swim(): void {
    console.log('鱼在游泳');
  }
  
  eat(): void {
    console.log('鱼在觅食');
  }
}

class Duck implements Walkable, Swimmable, Flyable, Eater {
  walk(): void {
    console.log('鸭子在走路');
  }
  
  swim(): void {
    console.log('鸭子在游泳');
  }
  
  fly(): void {
    console.log('鸭子在飞翔');
  }
  
  eat(): void {
    console.log('鸭子在吃东西');
  }
}

class Dog implements Walkable, Swimmable, Eater {
  walk(): void {
    console.log('狗在走路');
  }
  
  swim(): void {
    console.log('狗在游泳');
  }
  
  eat(): void {
    console.log('狗在吃东西');
  }
}
```

### 在 NestJS 中的应用

```typescript
// 细分的服务接口
interface UserReader {
  findById(id: string): Promise<User>;
  findByEmail(email: string): Promise<User>;
}

interface UserWriter {
  create(userData: CreateUserDto): Promise<User>;
  update(id: string, userData: UpdateUserDto): Promise<User>;
  delete(id: string): Promise<void>;
}

interface UserValidator {
  validateEmail(email: string): boolean;
  validatePassword(password: string): boolean;
}

// 只读用户服务
@Injectable()
export class UserQueryService implements UserReader {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>
  ) {}

  async findById(id: string): Promise<User> {
    return this.userRepository.findOne({ where: { id } });
  }

  async findByEmail(email: string): Promise<User> {
    return this.userRepository.findOne({ where: { email } });
  }
}

// 用户写入服务
@Injectable()
export class UserCommandService implements UserWriter {
  constructor(
    @InjectRepository(User)
    private userRepository: Repository<User>
  ) {}

  async create(userData: CreateUserDto): Promise<User> {
    const user = this.userRepository.create(userData);
    return this.userRepository.save(user);
  }

  async update(id: string, userData: UpdateUserDto): Promise<User> {
    await this.userRepository.update(id, userData);
    return this.findById(id);
  }

  async delete(id: string): Promise<void> {
    await this.userRepository.delete(id);
  }

  private async findById(id: string): Promise<User> {
    return this.userRepository.findOne({ where: { id } });
  }
}

// 用户验证服务
@Injectable()
export class UserValidationService implements UserValidator {
  validateEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  validatePassword(password: string): boolean {
    return password.length >= 8;
  }
}

// 控制器只依赖它需要的接口
@Controller('users')
export class UserController {
  constructor(
    private userQueryService: UserReader,
    private userCommandService: UserWriter,
    private userValidationService: UserValidator
  ) {}

  @Get(':id')
  async getUser(@Param('id') id: string) {
    return this.userQueryService.findById(id);
  }

  @Post()
  async createUser(@Body() userData: CreateUserDto) {
    if (!this.userValidationService.validateEmail(userData.email)) {
      throw new BadRequestException('Invalid email');
    }
    return this.userCommandService.create(userData);
  }
}
```

## 5. 依赖倒置原则（Dependency Inversion Principle, DIP）

### 定义
高层模块不应该依赖低层模块，二者都应该依赖抽象（接口或抽象类）。依赖关系应该通过抽象来实现，而不是通过具体实现。

### 违反原则的示例

```typescript
// ❌ 违反依赖倒置原则
class MySQLDatabase {
  connect() {
    console.log('连接 MySQL 数据库');
  }
  
  save(user: any) {
    console.log('保存用户到 MySQL');
  }
}

class UserRepository {
  private database: MySQLDatabase; // 直接依赖具体实现
  
  constructor() {
    this.database = new MySQLDatabase(); // 硬编码依赖
  }
  
  save(user: any) {
    this.database.connect();
    this.database.save(user);
  }
}
```

**问题分析**：
- UserRepository 直接依赖 MySQLDatabase 具体实现
- 难以测试和扩展
- 违反了依赖倒置原则

### 符合原则的示例

```typescript
// ✅ 符合依赖倒置原则
interface Database {
  connect(): void;
  save(obj: any): void;
}

class UserRepository {
  private database: Database; // 依赖抽象接口
  
  constructor(database: Database) { // 通过构造函数注入
    this.database = database;
  }
  
  save(user: any) {
    this.database.connect();
    this.database.save(user);
  }
}

// 具体实现
class MySQLDatabase implements Database {
  connect(): void {
    console.log('连接 MySQL 数据库');
  }
  
  save(obj: any): void {
    console.log('保存数据到 MySQL');
  }
}

class MongoDatabase implements Database {
  connect(): void {
    console.log('连接 MongoDB 数据库');
  }
  
  save(obj: any): void {
    console.log('保存数据到 MongoDB');
  }
}

// 使用示例 - 可以轻松切换数据库实现
const mysqlDatabase = new MySQLDatabase();
const userRepository = new UserRepository(mysqlDatabase);
userRepository.save({ id: 1, name: 'nick' });

const mongoDatabase = new MongoDatabase();
const userRepository2 = new UserRepository(mongoDatabase);
userRepository2.save({ id: 1, name: 'nick' });
```

### 在 NestJS 中的应用

```typescript
// 抽象接口
interface ILogger {
  log(message: string): void;
  error(message: string, error?: Error): void;
}

interface IEmailService {
  sendEmail(to: string, subject: string, body: string): Promise<void>;
}

interface ICacheService {
  get(key: string): Promise<any>;
  set(key: string, value: any, ttl?: number): Promise<void>;
  delete(key: string): Promise<void>;
}

// 具体实现
@Injectable()
export class ConsoleLogger implements ILogger {
  log(message: string): void {
    console.log(`[LOG] ${message}`);
  }
  
  error(message: string, error?: Error): void {
    console.error(`[ERROR] ${message}`, error);
  }
}

@Injectable()
export class FileLogger implements ILogger {
  log(message: string): void {
    // 写入文件逻辑
  }
  
  error(message: string, error?: Error): void {
    // 写入错误日志文件逻辑
  }
}

@Injectable()
export class SMTPEmailService implements IEmailService {
  async sendEmail(to: string, subject: string, body: string): Promise<void> {
    // SMTP 邮件发送逻辑
  }
}

@Injectable()
export class RedisCache implements ICacheService {
  async get(key: string): Promise<any> {
    // Redis 获取逻辑
  }
  
  async set(key: string, value: any, ttl?: number): Promise<void> {
    // Redis 设置逻辑
  }
  
  async delete(key: string): Promise<void> {
    // Redis 删除逻辑
  }
}

// 高层模块依赖抽象
@Injectable()
export class UserService {
  constructor(
    @Inject('ILogger') private logger: ILogger,
    @Inject('IEmailService') private emailService: IEmailService,
    @Inject('ICacheService') private cacheService: ICacheService,
    @InjectRepository(User) private userRepository: Repository<User>
  ) {}

  async createUser(userData: CreateUserDto): Promise<User> {
    try {
      this.logger.log(`Creating user: ${userData.email}`);
      
      // 检查缓存
      const cached = await this.cacheService.get(`user:${userData.email}`);
      if (cached) {
        this.logger.log('User found in cache');
        return cached;
      }

      // 创建用户
      const user = this.userRepository.create(userData);
      const savedUser = await this.userRepository.save(user);

      // 缓存用户
      await this.cacheService.set(`user:${userData.email}`, savedUser, 3600);

      // 发送欢迎邮件
      await this.emailService.sendEmail(
        userData.email,
        'Welcome!',
        'Welcome to our platform!'
      );

      this.logger.log(`User created successfully: ${savedUser.id}`);
      return savedUser;

    } catch (error) {
      this.logger.error('Failed to create user', error);
      throw error;
    }
  }
}

// 模块配置 - 可以轻松切换实现
@Module({
  providers: [
    UserService,
    {
      provide: 'ILogger',
      useClass: ConsoleLogger, // 可以切换为 FileLogger
    },
    {
      provide: 'IEmailService',
      useClass: SMTPEmailService,
    },
    {
      provide: 'ICacheService',
      useClass: RedisCache,
    },
  ],
})
export class UserModule {}
```

## SOLID 原则总结

### 原则间的关系

1. **SRP** 确保类的职责单一，降低耦合
2. **OCP** 通过抽象实现扩展，避免修改现有代码
3. **LSP** 确保继承关系的正确性，子类可以替换父类
4. **ISP** 避免接口臃肿，确保接口最小化
5. **DIP** 通过依赖抽象实现解耦，提高灵活性

### 实际开发中的应用策略

#### 1. 设计阶段
- 识别职责边界，应用 SRP
- 定义抽象接口，为 OCP 做准备
- 设计继承层次，遵循 LSP
- 拆分大接口，应用 ISP
- 识别依赖关系，应用 DIP

#### 2. 编码阶段
```typescript
// 综合应用 SOLID 原则的示例
interface IUserValidator {
  validate(userData: CreateUserDto): ValidationResult;
}

interface IUserRepository {
  save(user: User): Promise<User>;
  findByEmail(email: string): Promise<User | null>;
}

interface INotificationService {
  sendWelcomeNotification(user: User): Promise<void>;
}

@Injectable()
export class EmailUserValidator implements IUserValidator {
  validate(userData: CreateUserDto): ValidationResult {
    // SRP: 只负责验证逻辑
    const errors: string[] = [];
    
    if (!userData.email?.includes('@')) {
      errors.push('Invalid email format');
    }
    
    if (!userData.password || userData.password.length < 8) {
      errors.push('Password must be at least 8 characters');
    }
    
    return {
      isValid: errors.length === 0,
      errors
    };
  }
}

@Injectable()
export class DatabaseUserRepository implements IUserRepository {
  constructor(
    @InjectRepository(User)
    private repository: Repository<User>
  ) {}

  async save(user: User): Promise<User> {
    // SRP: 只负责数据持久化
    return this.repository.save(user);
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.repository.findOne({ where: { email } });
  }
}

@Injectable()
export class EmailNotificationService implements INotificationService {
  async sendWelcomeNotification(user: User): Promise<void> {
    // SRP: 只负责通知发送
    console.log(`Sending welcome email to ${user.email}`);
  }
}

@Injectable()
export class UserService {
  constructor(
    // DIP: 依赖抽象而非具体实现
    private validator: IUserValidator,
    private repository: IUserRepository,
    private notificationService: INotificationService
  ) {}

  async createUser(userData: CreateUserDto): Promise<User> {
    // 验证数据
    const validationResult = this.validator.validate(userData);
    if (!validationResult.isValid) {
      throw new BadRequestException(validationResult.errors.join(', '));
    }

    // 检查用户是否已存在
    const existingUser = await this.repository.findByEmail(userData.email);
    if (existingUser) {
      throw new ConflictException('User already exists');
    }

    // 创建用户
    const user = new User();
    user.email = userData.email;
    user.name = userData.name;
    
    const savedUser = await this.repository.save(user);

    // 发送通知
    await this.notificationService.sendWelcomeNotification(savedUser);

    return savedUser;
  }
}
```

#### 3. 测试阶段
```typescript
// SOLID 原则使测试变得简单
describe('UserService', () => {
  let userService: UserService;
  let mockValidator: IUserValidator;
  let mockRepository: IUserRepository;
  let mockNotificationService: INotificationService;

  beforeEach(() => {
    mockValidator = {
      validate: jest.fn()
    };
    
    mockRepository = {
      save: jest.fn(),
      findByEmail: jest.fn()
    };
    
    mockNotificationService = {
      sendWelcomeNotification: jest.fn()
    };

    userService = new UserService(
      mockValidator,
      mockRepository,
      mockNotificationService
    );
  });

  it('should create user successfully', async () => {
    // 由于依赖抽象，可以轻松 mock 所有依赖
    mockValidator.validate.mockReturnValue({ isValid: true, errors: [] });
    mockRepository.findByEmail.mockResolvedValue(null);
    mockRepository.save.mockResolvedValue({ id: 1, email: 'test@test.com' });

    const result = await userService.createUser({
      email: 'test@test.com',
      name: 'Test User',
      password: 'password123'
    });

    expect(result).toBeDefined();
    expect(mockNotificationService.sendWelcomeNotification).toHaveBeenCalled();
  });
});
```

SOLID 原则是编写高质量、可维护代码的基础，在 NestJS 开发中正确应用这些原则，可以显著提高代码的质量和项目的可维护性。

