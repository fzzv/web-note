# nest 基础

nest 核心其实就是 express + 装饰器，express 是用来处理请求的，装饰器是用来处理路由的。

## 创建项目

可以使用nest提供的脚手架进行创建项目

```bash
# 全局安装
npm install -g @nestjs/cli
# 创建项目
nest new 项目名
```

目录结构

```
src
	app.controller.spec.ts
	app.controller.ts
	app.module.ts
	app.service.ts
	main.ts
```

在实际开发的项目中，每个模块都会包含 controller、service、module、dto、entities

- controller 是处理路由，解析请求参数的。
- service 是处理业务逻辑的，比如操作数据库。
- dto 是封装请求参数的。
- entities 是封装对应数据库表的实体的。

## Nest的依赖包

- @nestjs/core Nest.js核心模块，提供构建、启动和管理Nest.js应用程序的基础设施
- @nestjs/common 包含了构建Nest.js应用程序基础设施和常用装饰器，像控制器、服务、中间件、守卫、拦截器、管道、异常过滤器等
- rxjs 用于构建异步和事件驱动程序的库。
- reflect-metadata 实现元编程的库，提供元数据反射API，可以在运行时检查和操作对象的元数据
- @nestjs/platform-express Nest的Express平台适配器，提供中间件、路由等功能

## 脚手架的常见命令

### nest generate

nest 命令除了可以生成整个项目外，还可以生成一些别的代码，比如 controller、service、module 等。

```bash
# 生成 module
nest generate module 模块名称
```

### nest build

打包 nest 项目 生成编译后的代码

### nest start

start 命令会先执行 build，再使用 node main.js 执行项目

### nest-cli.json 配置文件

nest脚手架的一些相关配置

```json
{
  "$schema": "https://json.schemastore.org/nest-cli",
  "collection": "@nestjs/schematics",
  "sourceRoot": "src",
  "compilerOptions": {
    "deleteOutDir": true, // 设置为 true，每次 build 都会都清空 dist 目录
    "webpack": false, // true就相当于nest build --webpack，设置为false，就是tsc了
    "plugins": [ // 插件
      {
        "name": "@nestjs/swagger",
        "options": {
          "classValidatorShim": true,
          "introspectComments": true
        }
      }
    ]
  }
}
```

### nest info

查看项目信息的，包括系统信息、 node、npm 和依赖版本

## 核心架构三要素

- **Controller**：路由定义与DTO验证
- **Provider**：服务层业务逻辑封装
- **Module**：模块化组织与依赖管理

## 常见的装饰器

### 核心架构装饰器

| 装饰器            | 作用描述                                  | 示例代码片段                                                 |
| :---------------- | :---------------------------------------- | :----------------------------------------------------------- |
| **@Controller()** | 定义 REST 控制器，处理 HTTP 请求          | `@Controller('users') export class UserController {}`        |
| **@Injectable()** | 标记类为可注入的 Provider（服务、仓库等） | `@Injectable() export class UserService {}`                  |
| **@Module()**     | 定义模块组织结构，声明依赖关系            | `@Module({ controllers: [UserController], providers: [UserService] })` |

### 路由方法

| 装饰器        | 对应 HTTP 方法     | 示例                               |
| :------------ | :----------------- | :--------------------------------- |
| **@Get()**    | GET                | `@Get('profile') getProfile() {}`  |
| **@Post()**   | POST               | `@Post() createUser() {}`          |
| **@Put()**    | PUT                | `@Put(':id') updateUser() {}`      |
| **@Delete()** | DELETE             | `@Delete(':id') deleteUser() {}`   |
| **@Patch()**  | PATCH              | `@Patch(':id') partialUpdate() {}` |
| **@All()**    | 匹配所有 HTTP 方法 | `@All() handleAll() {}`            |

### 请求参数获取

| 装饰器                      | 获取内容             | 示例                                                         |
| :-------------------------- | :------------------- | :----------------------------------------------------------- |
| **@Req()**                  | 请求对象（Request）  | `@Get() findAll(@Req() req: Request) {}`                     |
| **@Res()**                  | 响应对象（Response） | `@Get() getFile(@Res() res: Response) { res.sendFile(...) }` |
| **@Param(key?: string)**    | 路由参数             | `@Get(':id') findOne(@Param('id') id: string) {}`            |
| **@Query(key?: string)**    | 查询参数             | `@Get() search(@Query('keyword') keyword: string) {}`        |
| **@Body(key?: string)**     | 请求体数据           | `@Post() create(@Body() createUserDto: CreateUserDto) {}`    |
| **@Headers(name?: string)** | 请求头               | `@Get() checkAuth(@Headers('authorization') auth: string) {}` |
| **@Ip()**                   | 客户端 IP            | `@Get() logIp(@Ip() ip: string) {}`                          |
| **@Session()**              | Session 对象         | `@Get() getSession(@Session() session: Record<string, any>) {}` |

### 数据验证与转换

| 装饰器                | 作用                               | 典型用法                                                     |
| :-------------------- | :--------------------------------- | :----------------------------------------------------------- |
| **@UsePipes()**       | 应用管道进行数据转换和验证         | `@Post() @UsePipes(ValidationPipe) create() {}`              |
| **@ParseIntPipe()**   | 将参数转换为整数                   | `@Get(':id') findOne(@Param('id', ParseIntPipe) id: number) {}` |
| **@ParseUUIDPipe()**  | 验证参数是否符合 UUID 格式         | `@Get(':uuid') findByUUID(@Param('uuid', ParseUUIDPipe) uuid: string) {}` |
| **@ValidationPipe()** | 结合 class-validator 进行 DTO 验证 | `@Post() create(@Body(ValidationPipe) dto: CreateDto) {}`    |

### 依赖注入相关

| 装饰器          | 作用                           | 示例                                                         |
| :-------------- | :----------------------------- | :----------------------------------------------------------- |
| **@Inject()**   | 手动注入指定 token 的依赖      | `constructor(@Inject('CACHE') private cache: Cache) {}`      |
| **@Optional()** | 标记依赖为可选                 | `constructor(@Optional() @Inject('SERVICE') private service?: Service) {}` |
| **@Self()**     | 只从当前模块解析依赖           | `constructor(@Self() private config: ConfigService) {}`      |
| **@SkipSelf()** | 跳过当前模块，从父容器解析依赖 | `constructor(@SkipSelf() parentService: ParentService) {}`   |

### 拦截器与守卫

| 装饰器                 | 作用                                     |
| :--------------------- | :--------------------------------------- |
| **@UseInterceptors()** | 应用拦截器（日志、缓存、响应格式统一等） |
| **@UseGuards()**       | 应用守卫（权限验证、角色控制等）         |

### 异常处理

| 装饰器            | 作用                       |
| :---------------- | :------------------------- |
| **@Catch()**      | 定义异常过滤器捕获特定异常 |
| **@UseFilters()** | 应用异常过滤器             |

### WebSocket

| 装饰器                  | 作用                      |
| :---------------------- | :------------------------ |
| **@WebSocketGateway()** | 定义 WebSocket 网关       |
| **@SubscribeMessage()** | 处理客户端消息事件        |
| **@WebSocketServer()**  | 获取 WebSocket 服务器实例 |

配置相关

| 装饰器        | 作用                                  | 示例                                                         |
| :------------ | :------------------------------------ | :----------------------------------------------------------- |
| **@Global()** | 声明全局模块                          | `@Global() @Module({ providers: [ConfigService] })`          |
| **@Config()** | 快速获取配置（需配合 @nestjs/config） | `constructor(@Config('database') private dbConfig: DatabaseConfig) {}` |

### 实用工具

| 装饰器          | 作用             | 示例                                                         |
| :-------------- | :--------------- | :----------------------------------------------------------- |
| **@Header()**   | 设置响应头       | `@Post() @Header('Cache-Control', 'none') create() {}`       |
| **@HttpCode()** | 指定 HTTP 状态码 | `@Delete(':id') @HttpCode(204) delete() {}`                  |
| **@Redirect()** | 重定向响应       | `@Get() @Redirect('https://new-url.com', 301) redirect() {}` |

### 定时任务

| 装饰器          | 作用                   | 示例                                         |
| :-------------- | :--------------------- | :------------------------------------------- |
| **@Cron()**     | 定义 Cron 定时任务     | `@Cron('45 * * * * *') handleCron() { ... }` |
| **@Interval()** | 按间隔执行任务（毫秒） | `@Interval(10000) handleInterval() { ... }`  |
| **@Timeout()**  | 延迟执行一次任务       | `@Timeout(5000) handleTimeout() { ... }`     |

## 工作流程（MVC）

nest 应用跑起来后，会从 AppModule 开始解析，初始化 IoC 容器，加载所有的 service 到容器里，然后解析 controller 里的路由，接下来就可以接收请求了。

![image-20250315154917874](nest%E5%9F%BA%E7%A1%80.assets/image-20250315154917874.png)

## 接口实现多版本共存

应用开发完一版上线之后，还会不断的迭代。

后续可能需要修改已有的接口，但是为了兼容，之前版本的接口还要保留，就需要实现多个版本的接口共存。

创建个项目

```bash
nest new api-multi-version
```

生成一个模块

```bash
nest g resource user --no-spec
```

启动项目并请求 `/user` 接口

![image-20250408092938860](nest%E5%9F%BA%E7%A1%80.assets/image-20250408092938860.png)

### header 方式

添加一个v2版本的findAll方法

```ts
// user.controller
@Controller({
  path: 'user',
  version: '1',
})
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Post()
  create(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }

  @Get()
  findAll() {
    return this.userService.findAll();
  }

  @Version('2')
  @Get()
  findAll2() {
    return this.userService.findAll2();
  }
}

// user.service
@Injectable()
export class UserService {
  findAll2() {
    return 'This action returns all user V2 版本';
  }
}
```

main中启用 header 中传入字段的方式

```ts
async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  app.enableVersioning({
    type: VersioningType.HEADER,
    header: 'version',
  });

  await app.listen(process.env.PORT ?? 3000);
}
```

测试不同版本的返回值

v1:

![image-20250408103537793](nest%E5%9F%BA%E7%A1%80.assets/image-20250408103537793.png)

v2:

![image-20250408103553816](nest%E5%9F%BA%E7%A1%80.assets/image-20250408103553816.png)

其他版本：

![image-20250408103612304](nest%E5%9F%BA%E7%A1%80.assets/image-20250408103612304.png)

会发现如果传入错误的版本，会返回404，如果想要实现传入其他版本或不传，使用默认版本，可以在 Controller 中设置 version 为 VERSION_NEUTRAL

```ts
@Controller({
  path: 'user',
  version: VERSION_NEUTRAL,
})
export class UserController {
  constructor(private readonly userService: UserService) {}

  @Post()
  create(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }
  
  @Version('2')
  @Get()
  findAll2() {
    return this.userService.findAll2();
  }
  
  @Get()
  findAll() {
    return this.userService.findAll();
  }
}
```

> 需要注意的是，它是从上到下匹配，这样会导致v2失效，所以要把需要的版本移动到默认版本之上

也可以将v2版本的controller抽离到单独的controller 中

```bash
nest g controller user/user-v2 --no-spec --flat
```

user-v2.controller

```ts
import { Controller, Get } from '@nestjs/common';
import { UserService } from './user.service';

@Controller({
  path: 'user',
  version: '2',
})
export class UserV2Controller {
  constructor(private readonly userService: UserService) {}

  @Get()
  findAll2() {
    return this.userService.findAll2();
  }
}
```

user.module 中 v2的controller要放在前面

```ts
import { Module } from '@nestjs/common';
import { UserService } from './user.service';
import { UserController } from './user.controller';
import { UserV2Controller } from './user-v2.controller';

@Module({
  controllers: [UserV2Controller, UserController],
  providers: [UserService],
})
export class UserModule {}
```

### MEDIA_TYPE 方式

使用 media_type 的方式，在 Accpet 中传入key和值

```ts
app.enableVersioning({
    type: VersioningType.MEDIA_TYPE,
    key: 'v=',
});
```

![image-20250408114251018](nest%E5%9F%BA%E7%A1%80.assets/image-20250408114251018.png)

### URI 方式

配置URI方式

```ts
app.enableVersioning({
    type: VersioningType.URI,
});
```

请求的时候添加版本

![image-20250408114557371](nest%E5%9F%BA%E7%A1%80.assets/image-20250408114557371.png)

![image-20250408114543063](nest%E5%9F%BA%E7%A1%80.assets/image-20250408114543063.png)

可以看出 URI 方式不支持设置默认版本 `VERSION_NEUTRAL` 方式，需要设置具体的版本。

```ts
// 设置指定的版本，v1 和 v3 就会走这个 controller
@Controller({
  path: 'user',
  version: ['1', '3'],
})
```

### Custom 自定义方式

使用自定义函数，设置某种情况下才使用特定版本

例如url中有 fan 字段，就使用 v2 版本，如果 header 中有 disable-custom 就返回 404

```ts
app.enableVersioning({
    type: VersioningType.CUSTOM,
    extractor: (request: Request) => {
        if (request.headers['disable-custom']) {
            return '';
        }
        return request.url.includes('fan') ? '2' : '1';
    },
});
```

![image-20250408135341477](nest%E5%9F%BA%E7%A1%80.assets/image-20250408135341477.png)

![image-20250408135455731](nest%E5%9F%BA%E7%A1%80.assets/image-20250408135455731.png)

## 权限控制

常见的权限控制有两种：ACL、RBAC

| **维度**       | **ACL**                           | **RBAC**                       |
| :------------- | :-------------------------------- | :----------------------------- |
| **权限主体**   | **直接绑定用户与资源**            | **通过角色间接关联用户与权限** |
| **管理单元**   | 单个资源（如文件、API 端点）      | 角色（如管理员、编辑、游客）   |
| **灵活性**     | 高（细粒度控制）                  | 中（基于角色批量授权）         |
| **数据结构**   | 资源关联访问列表（用户/权限映射） | 三张表：用户 ↔ 角色 ↔ 权限     |
| **数据库设计** | 每个资源维护独立的权限记录        | 中心化管理角色和权限关系       |
| **权限检查**   | 检查用户是否在资源的允许列表中    | 检查用户角色是否具备所需权限   |

各自的适用场景

- ACL： 需要精确控制每个资源的具体访问权限（如文件系统、网络设备、特定 API 端点）
- RBAC：用户角色明确且权限批量分配（如企业后台系统、SaaS 平台）
