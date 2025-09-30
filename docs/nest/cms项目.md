# cms

## 初始化项目

```bash
npm install -g @nestjs/cli
nest new cms
```

## 创建模块

```bash
nest generate module admin
nest generate module api
nest generate module shared
```

### app.module.ts

在app.module.ts中引入模块

```typescript
import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module'; // [!code ++]
import { ApiModule } from './api/api.module'; // [!code ++]
import { SharedModule } from './shared/shared.module'; // [!code ++]

@Module({
  imports: [AdminModule, ApiModule, SharedModule], // [!code ++]
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
```

### eslint.config.mjs

配置取消eslint换行符警告

```js
// @ts-check
import eslint from '@eslint/js';
import eslintPluginPrettierRecommended from 'eslint-plugin-prettier/recommended';
import globals from 'globals';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  {
    ignores: ['eslint.config.mjs'],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommendedTypeChecked,
  eslintPluginPrettierRecommended,
  {
    languageOptions: {
      globals: {
        ...globals.node,
        ...globals.jest,
      },
      sourceType: 'commonjs',
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },
  {
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-floating-promises': 'warn',
      '@typescript-eslint/no-unsafe-argument': 'warn',
      'linebreak-style': ['error', 'auto'], // [!code ++]
    },
  },
);
```

## 支持会话

安装所需库

```bash
npm install express-session cookie-parser @nestjs/platform-express
```

### main.ts

```ts
import { NestFactory } from '@nestjs/core';
import session from 'express-session'; // [!code ++]
import cookieParser from 'cookie-parser'; // [!code ++]
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  // 配置 cookie 解析器
  app.use(cookieParser()); // [!code ++]
  // 配置 session
  app.use( // [!code ++]
    session({ // [!code ++]
      secret: 'secret-key', // [!code ++]
      resave: true, // 是否每次都重新保存 // [!code ++]
      saveUninitialized: true, // 是否保存未初始化的会话 // [!code ++]
      cookie: { // [!code ++]
        maxAge: 1000 * 60 * 60 * 24 * 7, // 7天 // [!code ++]
      }, // [!code ++]
    }), // [!code ++]
  ); // [!code ++]
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

## 模板

使用的是 handlebars 和 bootstrap

bootstrap 等静态资源放在 public 目录下

### 安装 handlebars 相关库

```bash
npm i express-handlebars
```

### 控制器

```bash
nest generate controller admin/controllers/dashboard --no-spec --flat
```

### dashboard.hbs

```handlebars
{{!-- views/dashboard.hbs --}}
<h1>{{title}}</h1>
```

### dashboard.controller.ts

```ts
import { Controller, Get, Render } from '@nestjs/common';

@Controller('admin')
export class DashboardController {
  @Get()
  @Render('dashboard')
  dashboard() {
    return { title: 'dashboard' }
  }
}
```

### 页面布局

`/views/partials/header.hbs`

```handlebars
<!-- 导航栏，使用navbar类来定义基本样式，navbar-expand-lg使其在大屏幕上展开，bg-light设置背景为浅色 -->
<nav class="navbar navbar-expand-lg bg-light">
  <!-- 流体容器，使导航栏在大屏幕上全宽展开 -->
  <div class="container-fluid">
    <!-- 导航栏品牌，链接到首页 -->
    <a class="navbar-brand" href="#">CMS</a>
    <!-- 折叠导航栏内容，navbar-collapse用于折叠和展开导航栏 -->
    <div class="collapse navbar-collapse">
      <!-- 导航栏菜单，使用ms-auto类使其自动右对齐 -->
      <ul class="navbar-nav ms-auto">
        <!-- 导航项，包含下拉菜单 -->
        <li class="nav-item dropdown">
          <!-- 下拉菜单的触发链接，使用dropdown-toggle类使其具有下拉功能，data-bs-toggle属性用于触发Bootstrap的下拉菜单插件 -->
          <a class="nav-link dropdown-toggle" href="#" data-bs-toggle="dropdown">
            欢迎
          </a>
          <!-- 下拉菜单内容，使用dropdown-menu类定义 -->
          <ul class="dropdown-menu">
            <!-- 下拉菜单项，使用dropdown-item类定义 -->
            <li><a class="dropdown-item">退出登录</a></li>
          </ul>
        </li>
      </ul>
    </div>
  </div>
</nav>
```

`/views/partials/sidebar.hbs`

```handlebars
<!-- 定义一个列，宽度在中等屏幕及以上为3，在大屏幕及以上为2，并且没有内边距 -->
<div class="col-md-3 col-lg-2 p-0">
  <!-- 定义一个手风琴组件，id为sidebarMenu -->
  <div class="accordion" id="sidebarMenu">
    <!-- 定义一个手风琴项目 -->
    <div class="accordion-item">
      <!-- 定义手风琴的头部，id为动态生成 -->
      <h2 class="accordion-header" id="heading{{id}}">
        <!-- 定义一个按钮，点击时折叠或展开手风琴内容 -->
        <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapse{{id}}">
          <!-- 按钮文本内容 -->
          权限管理
        </button>
      </h2>
      <!-- 定义手风琴的折叠内容，id为动态生成 -->
      <div id="collapse{{id}}" class="accordion-collapse collapse">
        <!-- 定义手风琴的主体内容 -->
        <div class="accordion-body">
          <!-- 定义一个列表组 -->
          <ul class="list-group">
            <!-- 定义一个列表项 -->
            <li class="list-group-item">
              <!-- 定义一个链接 -->
              <a href="">用户管理</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
```

`/views/layouts/main.hbs`

```handlebars
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>CMS后台管理页面</title>
  <link href="/css/bootstrap.min.css" rel="stylesheet" />
  <link href="/css/bootstrap-icons.min.css" rel="stylesheet">
  <script src="/js/jquery.min.js"></script>
  <script src="/js/bootstrap.bundle.min.js"></script>
</head>

<body>
  {{> header}}
  <div class="container-fluid">
    <div class="row">
      {{> sidebar}}
      <!-- 右侧管理页面 -->
      <div class="col-md-9 col-lg-10">
        <div class="container mt-4">
          {{{body}}}
        </div>
      </div>
    </div>
  </div>
</body>

</html>
```

### 配置静态资源目录和视图引擎

`main.ts`

```ts
import { NestFactory } from '@nestjs/core';
import session from 'express-session';
import cookieParser from 'cookie-parser';
import { join } from 'node:path'; // [!code ++]
import { engine } from 'express-handlebars'; // [!code ++]
import { NestExpressApplication } from '@nestjs/platform-express'; // [!code ++]
import { AppModule } from './app.module';

async function bootstrap() {
  // 使用 NestFactory 创建一个 NestExpressApplication 实例
  const app = await NestFactory.create<NestExpressApplication>(AppModule); // [!code ++]
  // 配置静态资源目录
  app.useStaticAssets(join(__dirname, '..', 'public')); // [!code ++]
  // 设置视图文件的基本目录
  app.setBaseViewsDir(join(__dirname, '..', 'views')); // [!code ++]
  // 设置视图引擎为 hbs（Handlebars）
  app.set('view engine', 'hbs'); // [!code ++]
  // 配置 Handlebars 引擎
  app.engine('hbs', engine({ // [!code ++]
    // 设置文件扩展名为 .hbs
    extname: '.hbs', // [!code ++]
    // 配置运行时选项
    runtimeOptions: { // [!code ++]
      // 允许默认情况下访问原型属性
      allowProtoPropertiesByDefault: true, // [!code ++]
      // 允许默认情况下访问原型方法
      allowProtoMethodsByDefault: true, // [!code ++]
    }, // [!code ++]
  })); // [!code ++]
  // 配置 cookie 解析器
  app.use(cookieParser());
  // 配置 session
  app.use(
    session({
      secret: 'secret-key',
      resave: true, // 是否每次都重新保存
      saveUninitialized: true, // 是否保存未初始化的会话
      cookie: {
        maxAge: 1000 * 60 * 60 * 24 * 7, // 7天
      },
    }),
  );
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

## 连接数据库

```bash
npm install @nestjs/config @nestjs/typeorm mysql2
```

### 用户实体 user.entity.ts

`src/shared/entities/user.entity.ts`

```ts
import { Entity, Column, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 50, unique: true })
  username: string;

  @Column()
  password: string;

  @Column({ length: 15, nullable: true })
  mobile: string;

  @Column({ length: 100, nullable: true })
  email: string;

  @Column({ default: 1 })
  status: number;

  @Column({ default: false })
  is_super: boolean;

  @Column({ default: 100 })
  sort: number;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  createdAt: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP', onUpdate: 'CURRENT_TIMESTAMP' })
  updatedAt: Date;
}
```

### configuration.service

```ts
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ConfigurationService {
  constructor(private configService: ConfigService) { }
  get mysqlHost(): string {
    return this.configService.get<string>('MYSQL_HOST')!;
  }
  get mysqlPort(): number {
    return this.configService.get<number>('MYSQL_PORT')!;
  }
  get mysqlDb(): string {
    return this.configService.get<string>('MYSQL_DB')!;
  }
  get mysqlUser(): string {
    return this.configService.get<string>('MYSQL_USER')!;
  }
  get mysqlPass(): string {
    return this.configService.get<string>('MYSQL_PASSWORD')!;
  }
  get mysqlConfig() {
    return {
      host: this.mysqlHost,
      port: this.mysqlPort,
      database: this.mysqlDb,
      username: this.mysqlUser,
      password: this.mysqlPass,
    };
  }
}
```

### 环境变量 .env

```
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=cms
MYSQL_USER=root
MYSQL_PASSWORD=root
```

### 配置数据库连接

`share.module.ts`

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';

@Global()
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    TypeOrmModule.forFeature([User]),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigurationService],
      useFactory: (configService: ConfigurationService) => ({
        type: 'mysql',
        ...configService.mysqlConfig,
        entities: [User],
        synchronize: true,
        autoLoadEntities: true,
        logging: false
      }),
    }),
  ],
  providers: [ConfigurationService],
  exports: [ConfigurationService],
})
export class ShareModule {}
```

## 用户接口

### 生成控制器

```bash
nest generate service share/services/user --no-spec --flat
nest generate controller admin/controllers/user --no-spec --flat
```

### 基础的curd

`mysql-base.service.ts`

```ts
import { Injectable } from '@nestjs/common';
import { Repository, FindOneOptions, ObjectLiteral, DeepPartial } from 'typeorm';
import { QueryDeepPartialEntity } from 'typeorm/query-builder/QueryPartialEntity.js';

@Injectable()
export abstract class MysqlBaseService<T extends ObjectLiteral> {
  constructor(private readonly repository: Repository<T>) {}

  async findAll(): Promise<T[]> {
    return this.repository.find();
  }
  async findOne(options: FindOneOptions<T>): Promise<T | null> {
    return this.repository.findOne(options);
  }
  async create(createDto: DeepPartial<T>): Promise<T | T[]> {
    const entity = this.repository.create(createDto);
    return this.repository.save(entity);
  }
  async update(id: number, updateDto: QueryDeepPartialEntity<T>) {
    return await this.repository.update(id, updateDto);
  }
  async delete(id: number) {
    return await this.repository.delete(id);
  }
}
```

### controller

`user.controller.ts`

定义一个接口用于获取所有用户

```ts
import { Controller, Get } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';

@Controller('admin/user')
export class UserController {

  constructor(private readonly userService: UserService) {}

  @Get()
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }
}
```

### service

`user.service.ts`

```ts
import { Injectable } from '@nestjs/common';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>
  ) {
    super(userRepository);
  }
}
```

## 用户curd

### 安装依赖

```bash
npm i class-validator  class-transformer
```

### class-validator

| 装饰器方法名          | 介绍                                                         |
| :-------------------- | :----------------------------------------------------------- |
| `@IsString()`         | 验证属性是否为字符串类型。                                   |
| `@IsInt()`            | 验证属性是否为整数类型。                                     |
| `@IsBoolean()`        | 验证属性是否为布尔值类型。                                   |
| `@IsNumber()`         | 验证属性是否为数字类型，可以指定选项，如整数、浮点数等。     |
| `@IsArray()`          | 验证属性是否为数组类型。                                     |
| `@IsEmail()`          | 验证属性是否为合法的电子邮件地址。                           |
| `@IsEnum()`           | 验证属性是否为指定枚举类型中的值。                           |
| `@IsDate()`           | 验证属性是否为日期类型。                                     |
| `@IsOptional()`       | 如果属性存在则进行验证，否则跳过此验证。                     |
| `@IsNotEmpty()`       | 验证属性是否不为空（不为 `null` 或 `undefined` 且不为空字符串）。 |
| `@IsEmpty()`          | 验证属性是否为空（`null` 或 `undefined` 或为空字符串）。     |
| `@IsDefined()`        | 验证属性是否已定义（不为 `undefined`）。                     |
| `@Min()`              | 验证属性的值是否大于或等于指定的最小值。                     |
| `@Max()`              | 验证属性的值是否小于或等于指定的最大值。                     |
| `@MinLength()`        | 验证字符串属性的长度是否大于或等于指定的最小长度。           |
| `@MaxLength()`        | 验证字符串属性的长度是否小于或等于指定的最大长度。           |
| `@Length()`           | 验证字符串属性的长度是否在指定的范围内。                     |
| `@Matches()`          | 验证字符串属性是否符合指定的正则表达式。                     |
| `@IsUUID()`           | 验证属性是否为合法的 UUID 格式。                             |
| `@IsUrl()`            | 验证属性是否为合法的 URL 格式。                              |
| `@IsIn()`             | 验证属性是否为给定值数组中的一个。                           |
| `@IsNotIn()`          | 验证属性是否不在给定值数组中。                               |
| `@IsPositive()`       | 验证数字属性是否为正数。                                     |
| `@IsNegative()`       | 验证数字属性是否为负数。                                     |
| `@IsLatitude()`       | 验证属性是否为合法的纬度值（范围：-90 到 90）。              |
| `@IsLongitude()`      | 验证属性是否为合法的经度值（范围：-180 到 180）。            |
| `@IsPhoneNumber()`    | 验证属性是否为合法的电话号码，支持不同国家的格式。           |
| `@IsCreditCard()`     | 验证属性是否为有效的信用卡号。                               |
| `@IsISO8601()`        | 验证属性是否为合法的 ISO 8601 日期格式。                     |
| `@IsJSON()`           | 验证属性是否为合法的 JSON 字符串。                           |
| `@IsIP()`             | 验证属性是否为合法的 IP 地址，可以指定版本（`IPv4` 或 `IPv6`）。 |
| `@IsPostalCode()`     | 验证属性是否为合法的邮政编码，支持不同国家的格式。           |
| `@IsHexColor()`       | 验证属性是否为合法的十六进制颜色代码。                       |
| `@IsCurrency()`       | 验证属性是否为合法的货币金额格式。                           |
| `@IsAlphanumeric()`   | 验证属性是否仅包含字母和数字。                               |
| `@IsAlpha()`          | 验证属性是否仅包含字母。                                     |
| `@IsLowercase()`      | 验证属性是否全部为小写字母。                                 |
| `@IsUppercase()`      | 验证属性是否全部为大写字母。                                 |
| `@IsBase64()`         | 验证属性是否为合法的 Base64 编码字符串。                     |
| `@IsDateString()`     | 验证属性是否为合法的日期字符串。                             |
| `@IsFQDN()`           | 验证属性是否为合法的完全合格域名（FQDN）。                   |
| `@IsMilitaryTime()`   | 验证属性是否为合法的 24 小时时间格式（军事时间）。           |
| `@IsMongoId()`        | 验证属性是否为合法的 MongoDB ObjectId。                      |
| `@IsPort()`           | 验证属性是否为合法的端口号（范围：0 到 65535）。             |
| `@IsISBN()`           | 验证属性是否为合法的 ISBN 格式。                             |
| `@IsISSN()`           | 验证属性是否为合法的 ISSN 格式。                             |
| `@IsRFC3339()`        | 验证属性是否为合法的 RFC 3339 日期格式。                     |
| `@IsBIC()`            | 验证属性是否为合法的银行标识代码（BIC）。                    |
| `@IsJWT()`            | 验证属性是否为合法的 JSON Web Token（JWT）。                 |
| `@IsEAN()`            | 验证属性是否为合法的欧洲商品编号（EAN）。                    |
| `@IsMACAddress()`     | 验证属性是否为合法的 MAC 地址。                              |
| `@IsHexadecimal()`    | 验证属性是否为合法的十六进制数值。                           |
| `@IsTimeZone()`       | 验证属性是否为合法的时区名称。                               |
| `@IsStrongPassword()` | 验证属性是否为强密码，支持自定义验证条件（如长度、字符类型）。 |
| `@IsISO31661Alpha2()` | 验证属性是否为合法的 ISO 3166-1 Alpha-2 国家代码。           |
| `@IsISO31661Alpha3()` | 验证属性是否为合法的 ISO 3166-1 Alpha-3 国家代码。           |
| `@IsEAN13()`          | 验证属性是否为合法的 EAN-13 格式。                           |
| `@IsEAN8()`           | 验证属性是否为合法的 EAN-8 格式。                            |
| `@IsISRC()`           | 验证属性是否为合法的国际标准录音代码（ISRC）。               |
| `@IsISO4217()`        | 验证属性是否为合法的 ISO 4217 货币代码。                     |
| `@IsIBAN()`           | 验证属性是否为合法的国际银行帐号（IBAN）。                   |
| `@IsRFC4180()`        | 验证属性是否为合法的 RFC 4180 CSV 格式。                     |
| `@IsISO6391()`        | 验证属性是否为合法的 ISO 639-1 语言代码。                    |
| `@IsISIN()`           | 验证属性是否为合法的国际证券识别码（ISIN）。                 |

| 名称                           | 介绍                                                         |
| :----------------------------- | :----------------------------------------------------------- |
| `ValidatorConstraint`          | 装饰器，用于定义自定义验证器。可以指定验证器名称和是否为异步。 |
| `ValidatorConstraintInterface` | 接口，用于实现自定义验证器的逻辑。需要实现 `validate` 和 `defaultMessage` 方法。 |
| `ValidationArguments`          | 类，用于传递给验证器的参数信息，包括当前被验证的对象、属性、约束和目标对象等。 |
| `registerDecorator`            | 函数，用于注册自定义装饰器，可以指定目标对象、属性、验证器和其他选项。 |
| `ValidationOptions`            | 接口，用于指定验证选项，如消息、组、每个属性的条件等。       |

### @nestjs/mapped-types

| 方法名             | 介绍                                                         |
| :----------------- | :----------------------------------------------------------- |
| `PartialType`      | 用于将给定类型的所有属性设置为可选属性，通常用于更新操作。   |
| `PickType`         | 用于从给定类型中选择特定的属性来构建一个新类型，只包含选中的属性。 |
| `OmitType`         | 用于从给定类型中排除特定的属性来构建一个新类型，排除指定的属性。 |
| `IntersectionType` | 用于将多个类型合并成一个新类型，包含所有类型的属性。         |
| `MappedType`       | 是一个抽象类型，允许对 DTO 进行进一步扩展或自定义。通常与其他工具一起使用，直接使用较少。 |

### @nestjs/swagger

| 装饰器名称              | 介绍                                                         |
| :---------------------- | :----------------------------------------------------------- |
| `@ApiTags`              | 用于给控制器或模块添加标签，用于对 API 进行分类。            |
| `@ApiOperation`         | 用于描述单个操作的目的和功能，通常用于描述控制器中的方法。   |
| `@ApiResponse`          | 用于指定 API 响应的状态码及其描述，支持定义多个响应。        |
| `@ApiParam`             | 用于描述路径参数，包括名称、类型和描述。                     |
| `@ApiQuery`             | 用于描述查询参数（即 URL 中的 `?key=value` 部分），包括名称、类型和描述。 |
| `@ApiBody`              | 用于描述请求体的结构，通常用于 `POST` 和 `PUT` 请求。        |
| `@ApiHeader`            | 用于描述 HTTP 头信息，包括名称、类型和描述。                 |
| `@ApiBearerAuth`        | 用于描述使用 Bearer Token 的身份验证方式。                   |
| `@ApiCookieAuth`        | 用于描述基于 Cookie 的身份验证方式。                         |
| `@ApiBasicAuth`         | 用于描述基本身份验证方式。                                   |
| `@ApiExcludeEndpoint`   | 用于从 Swagger 文档中排除某个特定的控制器方法。              |
| `@ApiProduces`          | 用于指定 API 方法返回的数据格式，如 `application/json`。     |
| `@ApiConsumes`          | 用于指定 API 方法可以消费的数据格式，如 `application/json`。 |
| `@ApiExtraModels`       | 用于引入额外的模型类，通常用于复杂的响应或嵌套对象。         |
| `@ApiHideProperty`      | 用于从模型类中排除某些属性，使其不在 Swagger 文档中显示。    |
| `@ApiSecurity`          | 用于为控制器方法指定安全机制，如 OAuth2。                    |
| `@ApiExcludeController` | 用于从 Swagger 文档中排除整个控制器。                        |
| `@ApiImplicitParam`     | （已弃用）用于描述隐式的路径参数，建议使用 `@ApiParam` 代替。 |
| `@ApiImplicitQuery`     | （已弃用）用于描述隐式的查询参数，建议使用 `@ApiQuery` 代替。 |
| `@ApiImplicitHeader`    | （已弃用）用于描述隐式的头信息，建议使用 `@ApiHeader` 代替。 |
| `@ApiImplicitBody`      | （已弃用）用于描述隐式的请求体，建议使用 `@ApiBody` 代替。   |

### class-transformer 

| 装饰器名称                 | 介绍                                                         |
| :------------------------- | :----------------------------------------------------------- |
| `@Exclude()`               | 将目标属性从序列化输出中排除，使其不被包含在最终的序列化结果中。 |
| `@Expose()`                | 将目标属性包括在序列化输出中，或者重命名序列化结果中的属性。 |
| `@Transform()`             | 提供自定义的转换逻辑，可以在序列化或反序列化过程中对属性进行转换。 |
| `@Type()`                  | 显式指定属性的类型，通常用于在序列化或反序列化过程中确保正确的类型转换，尤其是在数组或对象中。 |
| `@TransformPlainToClass()` | 将普通对象转换为类实例，使用此装饰器可以自动执行该转换。     |
| `@TransformClassToPlain()` | 将类实例转换为普通对象，使用此装饰器可以自动执行该转换。     |
| `@TransformClassToClass()` | 将一个类实例转换为另一个类实例，通常用于创建副本并在转换过程中应用特定规则。 |

### ClassSerializerInterceptor

`ClassSerializerInterceptor` 是一个内置的拦截器，用于在数据响应之前对数据进行序列化处理。它利用了 `class-transformer` 库，能够根据类定义中的装饰器（例如 `@Exclude` 和 `@Expose`）来自动转换类实例。这对确保敏感数据不会在 API 响应中暴露非常有用。

功能和用途：

- **自动序列化**：拦截控制器方法的返回值，并将类实例序列化为普通对象。
- **属性控制**：通过使用 `class-transformer` 装饰器（如 `@Exclude`、`@Expose`），可以精细控制哪些属性会被序列化和暴露。
- **安全性**：能够防止敏感数据（如密码）在 API 响应中被不小心暴露。
- **嵌套处理**：能够处理嵌套的对象和数组，保证整个数据结构的序列化规则一致。

### SerializeOptions

`SerializeOptions` 是一个装饰器，通常与 `ClassSerializerInterceptor` 一起使用。它允许你为整个控制器或特定的控制器方法设置序列化选项，进一步定制序列化行为。

功能和用途：

- **定制化策略**：你可以为序列化设置不同的策略，例如 `exposeAll` 或 `excludeAll`，来决定默认情况下是包含还是排除类的所有属性。
- **分组控制**：可以为不同的序列化场景设置不同的组（groups），使得同一个类在不同场景下可以以不同的方式序列化。

### 生成控制器

```bash
nest generate controller api/controllers/user --no-spec --flat
```

### 自定义装饰器

`alidation-and-transform.decorators.ts`

```ts
import { applyDecorators } from "@nestjs/common";
import { Type } from "class-transformer";
import { IsBoolean, IsEmail, IsNumber, IsOptional, IsString } from "class-validator";

// 可选字符串
export function IsOptionalString() {
  return applyDecorators(IsOptional(), IsString())
}

// 可选邮箱
export function IsOptionalEmail() {
  return applyDecorators(IsOptional(), IsEmail())
}

// 可选数字 并转换为数字
export function IsOptionalNumber() {
  return applyDecorators(IsOptional(), IsNumber(), Type(() => Number))
}

// 可选布尔值 并转换为布尔值
export function IsOptionalBoolean() {
  return applyDecorators(IsOptional(), IsBoolean(), Type(() => Boolean))
}
```

### 自定义验证器

`user-validators.ts`

```ts
import { Injectable } from "@nestjs/common";
import { registerDecorator, ValidationArguments, ValidationOptions, ValidatorConstraint, ValidatorConstraintInterface, } from "class-validator";

// 定义一个自定义验证器，名为 'startsWith'，不需要异步验证
@ValidatorConstraint({ name: 'startsWith', async: false })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 StartsWithConstraint 类并实现 ValidatorConstraintInterface 接口
export class StartsWithConstraint implements ValidatorConstraintInterface {
  // 定义验证逻辑，检查值是否以指定的前缀开头
  validate(value: any, args: ValidationArguments) {
    const [prefix] = args.constraints;
    return typeof value === 'string' && value.startsWith(prefix);
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    const [prefix] = args.constraints;
    return `${args.property} must start with ${prefix}`;
  }
}

// 定义一个自定义验证器，名为 'isUsernameUnique'，需要异步验证
@ValidatorConstraint({ name: 'isUsernameUnique', async: true })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 IsUsernameUniqueConstraint 类并实现 ValidatorConstraintInterface 接口
export class IsUsernameUniqueConstraint implements ValidatorConstraintInterface {
  // 定义验证逻辑，检查用户名是否唯一
  async validate(value: any, args: ValidationArguments) {
    const existingUsernames = ['ADMIN', 'USER', 'GUEST']; // 模拟已存在的用户名列表
    return !existingUsernames.includes(value);
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    return `${args.property} must be unique`;
  }
}

// 创建 StartsWith 装饰器工厂函数，用于给属性添加 'startsWith' 验证逻辑
export function StartsWith(prefix: string, validationOptions?: ValidationOptions) {
  return function (object: Object, propertyName: string) {
    registerDecorator({
      target: object.constructor, // 目标类
      propertyName: propertyName, // 目标属性名
      options: validationOptions, // 验证选项
      constraints: [prefix], // 传递给验证器的参数，如前缀
      validator: StartsWithConstraint, // 指定使用的验证器类
    });
  };
}

// 创建 IsUsernameUnique 装饰器工厂函数，用于给属性添加 'isUsernameUnique' 验证逻辑
export function IsUsernameUnique(validationOptions?: ValidationOptions) {
  return function (object: Object, propertyName: string) {
    registerDecorator({
      target: object.constructor, // 目标类
      propertyName: propertyName, // 目标属性名
      options: validationOptions, // 验证选项
      constraints: [], // 传递给验证器的参数，这里不需要
      validator: IsUsernameUniqueConstraint, // 指定使用的验证器类
    });
  };
}
```

### 返回结果共用vo

`vo/result.ts`

```ts
import { ApiProperty } from '@nestjs/swagger';

export class Result {
  
  @ApiProperty({ description: '操作是否成功', example: true })
  public success: boolean;

  @ApiProperty({ description: '操作的消息或错误信息', example: '操作成功' })
  public message: string;
  constructor(success: boolean, message?: string) {
    this.success = success;
    this.message = message || '';
  }

  static success(message: string) {
    return new Result(true, message);
  }

  static fail(message: string) {
    return new Result(false, message);
  }
}
```

### 配置 swagger 文档

`main.ts`

```ts
import { NestFactory } from '@nestjs/core';
import session from 'express-session';
import cookieParser from 'cookie-parser';
import { join } from 'node:path';
import { engine } from 'express-handlebars';
import { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger'; // [!code ++]
import { ValidationPipe } from '@nestjs/common'; // [!code ++]

async function bootstrap() {
  // 使用 NestFactory 创建一个 NestExpressApplication 实例
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  // 配置静态资源目录
  app.useStaticAssets(join(__dirname, '..', 'public'));
  // 设置视图文件的基本目录
  app.setBaseViewsDir(join(__dirname, '..', 'views'));
  // 设置视图引擎为 hbs（Handlebars）
  app.set('view engine', 'hbs');
  // 配置 Handlebars 引擎
  app.engine('hbs', engine({
    // 设置文件扩展名为 .hbs
    extname: '.hbs',
    // 配置运行时选项
    runtimeOptions: {
      // 允许默认情况下访问原型属性
      allowProtoPropertiesByDefault: true,
      // 允许默认情况下访问原型方法
      allowProtoMethodsByDefault: true,
    },
  }));
  // 配置 cookie 解析器
  app.use(cookieParser());
  // 配置 session
  app.use(
    session({
      secret: 'secret-key',
      resave: true, // 是否每次都重新保存
      saveUninitialized: true, // 是否保存未初始化的会话
      cookie: {
        maxAge: 1000 * 60 * 60 * 24 * 7, // 7天
      },
    }),
  );
  // 配置全局管道
  app.useGlobalPipes(new ValidationPipe({ transform: true })); // [!code ++]
  // 配置 Swagger
  const config = new DocumentBuilder() // [!code ++]
    // 设置标题
    .setTitle('CMS API') // [!code ++]
    // 设置描述
    .setDescription('CMS API Description') // [!code ++]
    // 设置版本
    .setVersion('1.0') // [!code ++]
    // 设置标签
    .addTag('CMS') // [!code ++]
    // 设置Cookie认证
    .addCookieAuth('connect.sid') // [!code ++]
    // 设置Bearer认证
    .addBearerAuth({ type: 'http', scheme: 'bearer' }) // [!code ++]
    // 构建配置
    .build(); // [!code ++]
  // 使用配置对象创建 Swagger 文档
  const document = SwaggerModule.createDocument(app, config); // [!code ++]
  // 设置 Swagger 模块的路径和文档对象，将 Swagger UI 绑定到 '/api-doc' 路径上
  SwaggerModule.setup('api-doc', app, document); // [!code ++]
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

### 给控制器和实体增加一些 swagger 描述

`admin/controller/user.controller`

```ts
import { Controller, Get } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger'; // [!code ++]

@ApiTags('admin/user') // [!code ++]
@Controller('admin/user')
export class UserController {

  constructor(private readonly userService: UserService) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回用户列表' }) // [!code ++]
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }
}
```

`admin/controller/dashboard.controller`

```ts
import { Controller, Get, Render } from '@nestjs/common';
import { ApiCookieAuth, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger'; // [!code ++]

@ApiTags('admin/dashboard') // [!code ++]
@Controller('admin')
export class DashboardController {
  @Get()
  @ApiCookieAuth() // [!code ++]
  @ApiOperation({ summary: '管理后台仪表盘' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回仪表盘页面' }) // [!code ++]
  @Render('dashboard')
  dashboard() {
    return { title: 'dashboard' }
  }
}
```

`entities/user.entity`

```ts
import { ApiHideProperty, ApiProperty } from '@nestjs/swagger'; // [!code ++]
import { Exclude, Transform } from 'class-transformer'; // [!code ++]
import { Entity, Column, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  @ApiProperty({ description: '用户ID', example: 1 }) // [!code ++]
  id: number;

  @Column({ length: 50, unique: true })
  @ApiProperty({ description: '用户名', example: 'admin' }) // [!code ++]
  username: string;

  @Column()
  @Exclude() // 在序列化时排除密码字段，不返回给前端 // [!code ++]
  @ApiHideProperty() // 隐藏密码字段，不在Swagger文档中显示 // [!code ++]
  password: string;

  @Column({ length: 15, nullable: true })
  @ApiProperty({ description: '手机号', example: '13124567890', format: '手机号码会被部分隐藏' }) // [!code ++]
  @Transform(({ value }) => value ? value.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') : value) // [!code ++]
  mobile: string;

  @Column({ length: 100, nullable: true })
  @ApiProperty({ description: '邮箱', example: 'admin@example.com' }) // [!code ++]
  email: string;

  @Column({ default: 1 })
  @ApiProperty({ description: '状态', example: 1, enum: [1, 2] }) // [!code ++]
  status: number;

  @Column({ default: false })
  @ApiProperty({ description: '是否超级管理员', example: false }) // [!code ++]
  is_super: boolean;

  @Column({ default: 100 })
  @ApiProperty({ description: '排序', example: 100 }) // [!code ++]
  sort: number;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '创建时间', example: '2021-01-01 00:00:00' }) // [!code ++]
  createdAt: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP', onUpdate: 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '更新时间', example: '2021-01-01 00:00:00' }) // [!code ++]
  updatedAt: Date;
}
```

## 用户列表

`views/user/user-list.hbs`

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>用户名</th>
      <th>邮箱</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
    </tr>
    {{/each}}
  </tbody>
</table>
```

修改 `admin/controllers/user.controller.ts` 控制器渲染用户列表

```ts
import { Controller, Get, Render } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(private readonly userService: UserService) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list') // [!code ++]
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }
}
```

## 新增用户

### utility.service

安装

```bash
npm install bcrypt
```

`utility.service.ts`

```ts
import { Injectable } from '@nestjs/common';
// 导入 bcrypt 库，用于处理密码哈希和验证
import bcrypt from 'bcrypt';
// 使用 Injectable 装饰器将类标记为可注入的服务
@Injectable()
export class UtilityService {
  // 定义一个异步方法，用于生成密码的哈希值
  async hashPassword(password: string): Promise<string> {
    // 生成一个盐值，用于增强哈希的安全性
    const salt = await bcrypt.genSalt();
    // 使用生成的盐值对密码进行哈希，并返回哈希结果
    return bcrypt.hash(password, salt);
  }
  // 定义一个异步方法，用于比较输入的密码和存储的哈希值是否匹配
  async comparePassword(password: string, hash: string): Promise<boolean> {
    // 使用 bcrypt 的 compare 方法比较密码和哈希值，返回比较结果（true 或 false）
    return bcrypt.compare(password, hash);
  }
}
```

### error.hbs

```handlebars
<h1>发生错误</h1>
<p>{{message}}</p>
<p>3秒后将自动跳转回上一个页面...</p>
<script>
  setTimeout(function () {
    window.history.back();
  }, 3000);
</script>
```

### user-form.hbs

```handlebars
<h1>添加用户</h1>
<form action="/admin/user" method="POST">
  <div class="mb-3">
    <label for="username" class="form-label">用户名</label>
    <input type="text" class="form-control" id="username" name="username" value="">
  </div>
  <div class="mb-3">
    <label for="username" class="form-label">密码</label>
    <input type="text" class="form-control" id="password" name="password" value="">
  </div>
  <div class="mb-3">
    <label for="email" class="form-label">邮箱</label>
    <input type="email" class="form-control" id="email" name="email" value="">
  </div>
  <div class="mb-3">
    <label for="status" class="form-label">状态</label>
    <select class="form-control" id="status" name="status">
      <option value="1">激活</option>
      <option value="0">未激活</option>
    </select>
  </div>
  <button type="submit" class="btn btn-primary">保存</button>
</form>
```

### user-validators

修改`user-validators.ts`  `IsUsernameUniqueConstraint` 从数据库中读取用户

```ts
import { Injectable } from "@nestjs/common";
import { registerDecorator, ValidationArguments, ValidationOptions, ValidatorConstraint, ValidatorConstraintInterface, } from "class-validator";
import { User } from "../entities/user.entity";
import { InjectRepository } from "@nestjs/typeorm"; // [!code ++]
import { Repository } from "typeorm"; // [!code ++]

// 定义一个自定义验证器，名为 'startsWith'，不需要异步验证
@ValidatorConstraint({ name: 'startsWith', async: false })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 StartsWithConstraint 类并实现 ValidatorConstraintInterface 接口
export class StartsWithConstraint implements ValidatorConstraintInterface {
  // 定义验证逻辑，检查值是否以指定的前缀开头
  validate(value: any, args: ValidationArguments) {
    const [prefix] = args.constraints;
    return typeof value === 'string' && value.startsWith(prefix);
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    const [prefix] = args.constraints;
    return `${args.property} must start with ${prefix}`;
  }
}

// 定义一个自定义验证器，名为 'isUsernameUnique'，需要异步验证
@ValidatorConstraint({ name: 'isUsernameUnique', async: true })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 IsUsernameUniqueConstraint 类并实现 ValidatorConstraintInterface 接口
export class IsUsernameUniqueConstraint implements ValidatorConstraintInterface {
  constructor( // [!code ++]
    @InjectRepository(User) private readonly repository: Repository<User> // [!code ++]
  ) { } // [!code ++]
  // 定义验证逻辑，检查用户名是否唯一
  async validate(value: any, args: ValidationArguments) {
    const existingUsernames = ['ADMIN', 'USER', 'GUEST']; // 模拟已存在的用户名列表 // [!code --]
    return !existingUsernames.includes(value); // [!code --]
    const user = await this.repository.findOne({ where: { username: value } }); // [!code ++]
    return !user; // [!code ++]
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    return `${args.property} must be unique`;
  }
}
```

使用 useContainer 配置依赖注入容器

```ts
import { NestFactory } from '@nestjs/core';
import session from 'express-session';
import cookieParser from 'cookie-parser';
import { join } from 'node:path';
import { engine } from 'express-handlebars';
import { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { ValidationPipe } from '@nestjs/common';
import { useContainer } from 'class-validator'; // [!code ++]

async function bootstrap() {
  // 使用 NestFactory 创建一个 NestExpressApplication 实例
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  // 使用 useContainer 配置依赖注入容器
  useContainer(app.select(AppModule), { fallbackOnErrors: true }); // [!code ++]
  // 配置静态资源目录
  app.useStaticAssets(join(__dirname, '..', 'public'));
  // 设置视图文件的基本目录
  app.setBaseViewsDir(join(__dirname, '..', 'views'));
  // 设置视图引擎为 hbs（Handlebars）
  app.set('view engine', 'hbs');
  // 配置 Handlebars 引擎
  app.engine('hbs', engine({
    // 设置文件扩展名为 .hbs
    extname: '.hbs',
    // 配置运行时选项
    runtimeOptions: {
      // 允许默认情况下访问原型属性
      allowProtoPropertiesByDefault: true,
      // 允许默认情况下访问原型方法
      allowProtoMethodsByDefault: true,
    },
  }));
  // 配置 cookie 解析器
  app.use(cookieParser());
  // 配置 session
  app.use(
    session({
      secret: 'secret-key',
      resave: true, // 是否每次都重新保存
      saveUninitialized: true, // 是否保存未初始化的会话
      cookie: {
        maxAge: 1000 * 60 * 60 * 24 * 7, // 7天
      },
    }),
  );
  // 配置全局管道
  app.useGlobalPipes(new ValidationPipe({ transform: true }));
  // 配置 Swagger
  const config = new DocumentBuilder()
    // 设置标题
    .setTitle('CMS API')
    // 设置描述
    .setDescription('CMS API Description')
    // 设置版本
    .setVersion('1.0')
    // 设置标签
    .addTag('CMS')
    // 设置Cookie认证
    .addCookieAuth('connect.sid')
    // 设置Bearer认证
    .addBearerAuth({ type: 'http', scheme: 'bearer' })
    // 构建配置
    .build();
  // 使用配置对象创建 Swagger 文档
  const document = SwaggerModule.createDocument(app, config);
  // 设置 Swagger 模块的路径和文档对象，将 Swagger UI 绑定到 '/api-doc' 路径上
  SwaggerModule.setup('api-doc', app, document);
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

`dto/user.dto.ts`

```ts
import { IsString, Validate } from "class-validator";
import { StartsWithConstraint, IsUsernameUniqueConstraint } from "../validators/user-validators";
import { ApiProperty, ApiPropertyOptional, PartialType } from "@nestjs/swagger"
import { IsOptionalString, IsOptionalEmail, IsOptionalNumber, IsOptionalBoolean } from "../decorators/alidation-and-transform.decorators";

export class CreateUserDto {
  @ApiProperty({ description: '用户名，必须唯一且以指定前缀开头', example: 'user_john_doe' })
  @IsString()
  @Validate(StartsWithConstraint, ['user_'], {
    message: `用户名必须以 "user_" 开头`,
  })
  @Validate(IsUsernameUniqueConstraint, { message: '用户名已存在' })
  username: string;

  @ApiProperty({ description: '密码', example: 'securePassword123' })
  @IsString()
  password: string;

  @ApiPropertyOptional({ description: '手机号', example: '13124567890' })
  @IsOptionalString()
  mobile?: string;

  @ApiPropertyOptional({ description: '邮箱地址', example: 'john.doe@example.com' })
  @IsOptionalEmail()
  email?: string;

  @ApiPropertyOptional({ description: '用户状态', example: 1 })
  @IsOptionalNumber()
  status?: number;

  @ApiPropertyOptional({ description: '是否为超级管理员', example: true })
  @IsOptionalBoolean()
  is_super?: boolean;
}

export class UpdateUserDto extends PartialType(CreateUserDto) {
  @ApiProperty({ description: '用户ID', example: 1 })
  @IsOptionalNumber()
  id: number;
}
```

### share.module.ts

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { UtilityService } from './services/utility.service'; // [!code ++]
import { IsUsernameUniqueConstraint } from './validators/user-validators'; // [!code ++]

@Global()
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    TypeOrmModule.forFeature([User]),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigurationService],
      useFactory: (configService: ConfigurationService) => ({
        type: 'mysql',
        ...configService.mysqlConfig,
        entities: [User],
        synchronize: true,
        autoLoadEntities: true,
        logging: false
      }),
    }),
  ],
  providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint], // [!code ++]
  exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint], // [!code ++]
})
export class ShareModule {}
```

### 自定义异常过滤器 admin-exception.filter

```ts
import { ExceptionFilter, Catch, ArgumentsHost, HttpException, BadRequestException } from '@nestjs/common';
// 导入 express 的 Response 对象，用于构建 HTTP 响应
import { Response } from 'express';
// 使用 @Catch 装饰器捕获所有 HttpException 异常
@Catch(HttpException)
export class AdminExceptionFilter implements ExceptionFilter {
  // 实现 catch 方法，用于处理捕获的异常
  catch(exception: HttpException, host: ArgumentsHost) {
    // 获取当前 HTTP 请求上下文
    const ctx = host.switchToHttp();
    // 获取 HTTP 响应对象
    const response = ctx.getResponse<Response>();
    // 获取异常的 HTTP 状态码
    const status = exception.getStatus();
    // 初始化错误信息，默认为异常的消息
    let errorMessage = exception.message;
    // 如果异常是 BadRequestException 类型，进一步处理错误信息
    if (exception instanceof BadRequestException) {
      // 获取异常的响应体
      const responseBody: any = exception.getResponse();
      // 检查响应体是否是对象并且包含 message 属性
      if (typeof responseBody === 'object' && responseBody.message) {
        // 如果 message 是数组，则将其拼接成字符串，否则直接使用 message
        errorMessage = Array.isArray(responseBody.message)
          ? responseBody.message.join(', ')
          : responseBody.message;
      }
    }
    // 使用响应对象构建并发送错误页面，包含错误信息和重定向 URL
    response.status(status).render('error', {
      message: errorMessage,
      redirectUrl: ctx.getRequest().url,
    });
  }
}
```

注入过滤器

```ts
import { Module } from '@nestjs/common';
import { DashboardController } from './controllers/dashboard.controller';
import { UserController } from './controllers/user.controller';
import { AdminExceptionFilter } from './filters/admin-exception.filter'; // [!code ++]

@Module({
  controllers: [DashboardController, UserController],
  providers: [{ // [!code ++]
    provide: 'APP_FILTER', // [!code ++]
    useClass: AdminExceptionFilter, // [!code ++]
  }], // [!code ++]
})
export class AdminModule {}
```

### user.controller.ts

增加新增用户表单页面和新增用户接口

```ts
import { Body, Controller, Get, Post, Redirect, Render } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service'; // [!code ++]
import { CreateUserDto } from 'src/share/dtos/user.dto'; // [!code ++]

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService // [!code ++]
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }

  @Get('create') // [!code ++]
  @ApiOperation({ summary: '添加用户(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' }) // [!code ++]
  @Render('user/user-form') // [!code ++]
  async create() { // [!code ++]
    return { user: {} }; // [!code ++]
  } // [!code ++]

  @Post() // [!code ++]
  @Redirect('/admin/user') // [!code ++]
  @ApiOperation({ summary: '添加用户(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' }) // [!code ++]
  async createUser(@Body() createUserDto: CreateUserDto) { // [!code ++]
    console.log(createUserDto, 'createUserDto') // [!code ++]
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password); // [!code ++]
    await this.userService.create({ ...createUserDto, password: hashedPassword }); // [!code ++]
    return { url: '/admin/user', success: true, message: '用户添加成功' }; // [!code ++]
  } // [!code ++]
}

```

## 编辑用户

### 中间件

```ts
import { NextFunction, Request, Response } from "express";

/**
 * HTML 的 <form> 标签默认只支持 GET 和 POST
 * 但 RESTful API 常常需要 PUT、PATCH、DELETE 等方法
 * 为了绕过这个限制，前端可以在表单里加一个隐藏字段 _method，把要真正使用的 HTTP 方法放进去。
 * example:
 * <form action="/users/1" method="POST">
 *   <input type="hidden" name="_method" value="DELETE">
 *   <button type="submit">Delete User</button>
 * </form>
 */
function methodOverride(req: Request, res: Response, next: NextFunction) {
  if (req.body && typeof req.body === 'object' && '_method' in req.body) {
    req.method = req.body._method.toUpperCase();
    delete req.body._method;
  }
  next();
}

export default methodOverride;
```

配置中间件

```ts
import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common'; // [!code ++]
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module';
import { ApiModule } from './api/api.module';
import { ShareModule } from './share/share.module';
import methodOverride from './share/middlewares/method-override'; // [!code ++]

@Module({
  imports: [AdminModule, ApiModule, ShareModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule implements NestModule { // [!code ++]
  configure(consumer: MiddlewareConsumer) { // [!code ++]
    consumer.apply(methodOverride).forRoutes('*'); // [!code ++]
  } // [!code ++]
}
```

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>用户名</th>
      <th>邮箱</th>
      <th>操作</th> // [!code ++]
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td> // [!code ++]
        <a href="/admin/user/edit/{{this.id}}">编辑</a> // [!code ++]
      </td> // [!code ++]
    </tr>
    {{/each}}
  </tbody>
</table>
```

### user-form.hbs

```handlebars
<h1>{{#if user.id}}编辑用户{{else}}添加用户{{/if}}</h1> // [!code ++]
<form action="{{#if user.id}}/admin/user/{{user.id}}{{else}}/admin/user{{/if}}" method="POST"> // [!code ++]
  {{#if user.id}} // [!code ++]
    <input type="hidden" name="_method" value="PUT"> // [!code ++]
  {{/if}} // [!code ++]
  <div class="mb-3">
    <label for="username" class="form-label">用户名</label>
    <input type="text" class="form-control" id="username" name="username" value="{{user.username}}"> // [!code ++]
  </div>
  <div class="mb-3">
    <label for="username" class="form-label">密码</label>
    <input type="text" class="form-control" id="password" name="password" value="">
  </div>
  <div class="mb-3">
    <label for="email" class="form-label">邮箱</label>
    <input type="email" class="form-control" id="email" name="email" value="{{user.email}}"> // [!code ++]
  </div>
  <div class="mb-3">
    <label for="status" class="form-label">状态</label>
    <select class="form-control" id="status" name="status">
      <option value="1" {{#if user.status}}selected{{/if}}>激活</option> // [!code ++]
      <option value="0" {{#unless user.status}}selected{{/unless}}>未激活</option> // [!code ++]
    </select>
  </div>
  <button type="submit" class="btn btn-primary">保存</button>
</form>
```

### user.controller

```ts
import { Body, Controller, Get, NotFoundException, Param, Post, Put, ParseIntPipe, Redirect, Render } from '@nestjs/common'; // [!code ++]
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto'; // [!code ++]

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }

  @Get('create')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  @Render('user/user-form')
  async create() {
    return { user: {} };
  }

  @Post()
  @Redirect('/admin/user')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  async createUser(@Body() createUserDto: CreateUserDto) {
    console.log(createUserDto, 'createUserDto')
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id') // [!code ++]
  @ApiOperation({ summary: '编辑用户(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' }) // [!code ++]
  @Render('user/user-form') // [!code ++]
  async edit(@Param('id', ParseIntPipe) id: number) { // [!code ++]
    const user = await this.userService.findOne({ where: { id } }); // [!code ++]
    if (!user) { // [!code ++]
      throw new NotFoundException('用户不存在'); // [!code ++]
    } // [!code ++]
    return { user }; // [!code ++]
  } // [!code ++]

  @Put(':id') // [!code ++]
  @Redirect('/admin/user') // [!code ++]
  @ApiOperation({ summary: '编辑用户(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' }) // [!code ++]
  async updateUser(@Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto) { // [!code ++]
    if (updateUserDto.password) { // [!code ++]
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password); // [!code ++]
    } else { // [!code ++]
      delete updateUserDto.password; // [!code ++]
    } // [!code ++]
    await this.userService.update(id, updateUserDto); // [!code ++]
    return { url: '/admin/user', success: true, message: '用户更新成功' }; // [!code ++]
  } // [!code ++]
}
```

## 查看用户信息

### user-detail.hbs

```handlebars
<h1>用户详情</h1>
<div class="mb-3">
  <label class="form-label">用户名:</label>
  <p class="form-control-plaintext">{{user.username}}</p>
</div>
<div class="mb-3">
  <label class="form-label">邮箱:</label>
  <p class="form-control-plaintext">{{user.email}}</p>
</div>
<div class="mb-3">
  <label class="form-label">状态:</label>
  <p class="form-control-plaintext">{{#if user.status}}激活{{else}}未激活{{/if}}</p>
</div>
<a href="/admin/user/edit/{{user.id}}" class="btn btn-warning">编辑</a>
<a href="/admin/user" class="btn btn-secondary">返回列表</a>
```

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>用户名</th>
      <th>邮箱</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a> // [!code ++]
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
```

### user.controller.ts

```ts
import { Body, Controller, Get, NotFoundException, Param, Post, Put, ParseIntPipe, Redirect, Render } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }

  @Get('create')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  @Render('user/user-form')
  async create() {
    return { user: {} };
  }

  @Post()
  @Redirect('/admin/user')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  async createUser(@Body() createUserDto: CreateUserDto) {
    console.log(createUserDto, 'createUserDto')
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  @Render('user/user-form')
  async edit(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Put(':id')
  @Redirect('/admin/user')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  async updateUser(@Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto) {
    if (updateUserDto.password) {
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password);
    } else {
      delete updateUserDto.password;
    }
    await this.userService.update(id, updateUserDto);
    return { url: '/admin/user', success: true, message: '用户更新成功' };
  }

  @Get(':id') // [!code ++]
  @ApiOperation({ summary: '获取用户详情(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回用户详情' }) // [!code ++]
  @Render('user/user-detail') // [!code ++]
  async findOne(@Param('id', ParseIntPipe) id: number) { // [!code ++]
    const user = await this.userService.findOne({ where: { id } }); // [!code ++]
    if (!user) { // [!code ++]
      throw new NotFoundException('用户不存在'); // [!code ++]
    } // [!code ++]
    return { user }; // [!code ++]
  } // [!code ++]
}
```

## 删除用户

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>用户名</th>
      <th>邮箱</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a> // [!code ++]
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<script> // [!code ++]
  function deleteUser(id) { // [!code ++]
    if (confirm('确定要删除该用户吗？')) { // [!code ++]
      $.ajax({ // [!code ++]
        url: '/admin/user/' + id, // [!code ++]
        type: 'DELETE', // [!code ++]
        success: function (res) { // [!code ++]
          if (res.success) { // [!code ++]
            window.location.reload() // [!code ++]
          } // [!code ++]
        } // [!code ++]
      }) // [!code ++]
    } // [!code ++]
  } // [!code ++]
</script> // [!code ++]
```

### user.controller.ts

```ts
import { Body, Controller, Delete, Get, NotFoundException, Param, ParseIntPipe, Post, Put, Redirect, Render } from '@nestjs/common'; // [!code ++]
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }

  @Get('create')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  @Render('user/user-form')
  async create() {
    return { user: {} };
  }

  @Post()
  @Redirect('/admin/user')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  async createUser(@Body() createUserDto: CreateUserDto) {
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  @Render('user/user-form')
  async edit(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Put(':id')
  @Redirect('/admin/user')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  async updateUser(@Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto) {
    if (updateUserDto.password) {
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password);
    } else {
      delete updateUserDto.password;
    }
    await this.userService.update(id, updateUserDto);
    return { url: '/admin/user', success: true, message: '用户更新成功' };
  }

  @Get(':id')
  @ApiOperation({ summary: '获取用户详情(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户详情' })
  @Render('user/user-detail')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Delete(':id') // [!code ++]
  @ApiOperation({ summary: '删除用户(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回删除用户页面' }) // [!code ++]
  async deleteUser(@Param('id', ParseIntPipe) id: number) { // [!code ++]
    await this.userService.delete(id); // [!code ++]
    return { success: true, message: '用户删除成功' }; // [!code ++]
  } // [!code ++]
}
```

## 修改用户状态

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th> // [!code ++]
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td> // [!code ++]
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}"> // [!code ++]
          {{#if this.status}} // [!code ++]
          <i class="bi bi-check-circle-fill text-success"></i> // [!code ++]
          {{else}} // [!code ++]
          <i class="bi bi-x-circle-fill text-danger"></i> // [!code ++]
          {{/if}} // [!code ++]
        </span> // [!code ++]
      </td> // [!code ++]
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<script>
  $(function () { // [!code ++]
    $('.status-toggle').on('click', function () { // [!code ++]
      const $this = $(this); // [!code ++]
      const userId = $this.data('id'); // [!code ++]
      const currentStatus = $this.data('status'); // [!code ++]
      const newStatus = currentStatus === 1 ? 0 : 1; // [!code ++]
      $.ajax({ // [!code ++]
        url: `/admin/user/${userId}`, // [!code ++]
        type: 'PUT', // [!code ++]
        contentType: 'application/json', // [!code ++]
        headers: { // [!code ++]
          'accept': 'application/json' // [!code ++]
        }, // [!code ++]
        data: JSON.stringify({ status: newStatus }), // [!code ++]
        success: function (response) { // [!code ++]
          if (response.success) { // [!code ++]
            $this.data('status', newStatus); // [!code ++]
            $this.html(`<i class="bi ${newStatus ? "bi-check-circle-fill" : "bi-x-circle-fill"} ${newStatus ? "text-success" : "text-danger"}"></i>`); // [!code ++]
          } // [!code ++]
        }, // [!code ++]
        error: function (error) { // [!code ++]
          const { responseJSON } = error; // [!code ++]
          alert(responseJSON.message); // [!code ++]
        } // [!code ++]
      }); // [!code ++]
    }); // [!code ++]
  }); // [!code ++]
  function deleteUser(id) {
    if (confirm('确定要删除该用户吗？')) {
      $.ajax({
        url: '/admin/user/' + id,
        type: 'DELETE',
        success: function (res) {
          if (res.success) {
            window.location.reload()
          }
        }
      })
    }
  }
</script>
```

### admin-exception.filter

```ts
import { ExceptionFilter, Catch, ArgumentsHost, HttpException, BadRequestException } from '@nestjs/common';
// 导入 express 的 Response 对象，用于构建 HTTP 响应
import { Response } from 'express'; // [!code ++]
// 使用 @Catch 装饰器捕获所有 HttpException 异常
@Catch(HttpException)
export class AdminExceptionFilter implements ExceptionFilter {
  // 实现 catch 方法，用于处理捕获的异常
  catch(exception: HttpException, host: ArgumentsHost) {
    // 获取当前 HTTP 请求上下文
    const ctx = host.switchToHttp();
    // 获取当前 HTTP 请求对象
    const request = ctx.getRequest<Request>(); // [!code ++]
    // 获取 HTTP 响应对象
    const response = ctx.getResponse<Response>();
    // 获取异常的 HTTP 状态码
    const status = exception.getStatus();
    // 初始化错误信息，默认为异常的消息
    let errorMessage = exception.message;
    // 如果异常是 BadRequestException 类型，进一步处理错误信息
    if (exception instanceof BadRequestException) {
      // 获取异常的响应体
      const responseBody: any = exception.getResponse();
      // 检查响应体是否是对象并且包含 message 属性
      if (typeof responseBody === 'object' && responseBody.message) {
        // 如果 message 是数组，则将其拼接成字符串，否则直接使用 message
        errorMessage = Array.isArray(responseBody.message)
          ? responseBody.message.join(', ')
          : responseBody.message;
      }
    }
    // 如果请求头中包含 'application/json'，则返回 JSON 响应
    if (request.headers['accept'] === 'application/json') { // [!code ++]
      response.status(status).json({ // [!code ++]
        statusCode: status, // [!code ++]
        message: errorMessage // [!code ++]
      }); // [!code ++]
    } else { // [!code ++]
      // 使用响应对象构建并发送错误页面，包含错误信息和重定向 URL
      response.status(status).render('error', {
        message: errorMessage,
        redirectUrl: ctx.getRequest().url,
      });
    } // [!code ++]
  }
}
```

### user-controller

```ts
import { Body, Controller, Delete, Get, NotFoundException, Param, ParseIntPipe, Headers, Post, Put, Redirect, Render, Res, UseFilters } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';
import { AdminExceptionFilter } from '../filters/admin-exception.filter'; // [!code ++]
import type { Response } from 'express'; // [!code ++]

@ApiTags('admin/user')
@UseFilters(AdminExceptionFilter)<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>排序</th>
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>
        <span class="sort-text" data-id="{{this.id}}">{{this.sort}}</span>
        <input type="number" class="form-control sort-input d-none" style="width:80px" data-id="{{this.id}}"
          value="{{this.sort}}">
      </td>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}">
          {{#if this.status}}
          <i class="bi bi-check-circle-fill text-success"></i>
          {{else}}
          <i class="bi bi-x-circle-fill text-danger"></i>
          {{/if}}
        </span>
      </td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<script>
  $(function () {
    $('.sort-text').on('dblclick', function () {
      const userId = $(this).data('id');
      $(this).addClass('d-none');
      $(`.sort-input[data-id="${userId}"]`).removeClass('d-none').focus();
    });

    $('.sort-input').on('blur', function () {
      const userId = $(this).data('id');
      const newSort = $(this).val();
      $(this).addClass('d-none');
      $(`.sort-text[data-id="${userId}"]`).removeClass('d-none').text(newSort);
      $.ajax({
        url: `/admin/user/${userId}`,
        type: 'PUT',
        contentType: 'application/json',
        headers: {
          'accept': 'application/json'
        },
        data: JSON.stringify({ sort: newSort }),
        success: function (response) {
          if (response.success) {
            $(`.sort-text[data-id="${userId}"]`).text(newSort);
          }
        }
      });
    });

    $('.sort-input').on('keypress', function (e) {
      if (e.which == 13) {
        $(this).blur();
      }
    });
    $('.status-toggle').on('click', function () {
      const $this = $(this);
      const userId = $this.data('id');
      const currentStatus = $this.data('status');
      const newStatus = currentStatus === 1 ? 0 : 1;
      $.ajax({
        url: `/admin/user/${userId}`,
        type: 'PUT',
        contentType: 'application/json',
        headers: {
          'accept': 'application/json'
        },
        data: JSON.stringify({ status: newStatus }),
        success: function (response) {
          if (response.success) {
            $this.data('status', newStatus);
            $this.html(`<i class="bi ${newStatus ? "bi-check-circle-fill" : "bi-x-circle-fill"} ${newStatus ? "text-success" : "text-danger"}"></i>`);
          }
        },
        error: function (error) {
          const { responseJSON } = error;
          alert(responseJSON.message);
        }
      });
    });
  });
  function deleteUser(id) {
    if (confirm('确定要删除该用户吗？')) {
      $.ajax({
        url: '/admin/user/' + id,
        type: 'DELETE',
        success: function (res) {
          if (res.success) {
            window.location.reload()
          }
        }
      })
    }
  }
</script>

@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }

  @Get('create')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  @Render('user/user-form')
  async create() {
    return { user: {} };
  }

  @Post()
  @Redirect('/admin/user')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  async createUser(@Body() createUserDto: CreateUserDto) {
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  @Render('user/user-form')
  async edit(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Put(':id')
  @Redirect('admin/user')  // [!code --]
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  async updateUser(
    @Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto, 
    @Res() res: Response, @Headers('accept') accept: string
  ) {
    if (updateUserDto.password) {
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password);
    } else {
      delete updateUserDto.password;
    }
    await this.userService.update(id, updateUserDto);
    return { url: '/admin/user', success: true, message: '用户更新成功' }; // [!code --]
    if (accept.includes('application/json')) { // [!code ++]
      return res.json({ success: true, message: '用户更新成功' }); // [!code ++]
    } else { // [!code ++]
      return res.redirect('/admin/user'); // [!code ++]
    } // [!code ++]
  }

  @Get(':id')
  @ApiOperation({ summary: '获取用户详情(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户详情' })
  @Render('user/user-detail')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Delete(':id')
  @ApiOperation({ summary: '删除用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回删除用户页面' })
  async deleteUser(@Param('id', ParseIntPipe) id: number) {
    await this.userService.delete(id);
    return { success: true, message: '用户删除成功' };
  }
}
```

## 用户排序

### user-list

```handlebars
<h1>用户列表</h1>
<table class="table">
  <thead>
    <tr>
      <th>序号</th> // [!code ++]
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td> // [!code ++]
        <span class="sort-text" data-id="{{this.id}}">{{this.sort}}</span> // [!code ++]
        <input type="number" class="form-control sort-input d-none" style="width:80px" data-id="{{this.id}}"
          value="{{this.sort}}"> // [!code ++]
      </td> // [!code ++]
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}">
          {{#if this.status}}
          <i class="bi bi-check-circle-fill text-success"></i>
          {{else}}
          <i class="bi bi-x-circle-fill text-danger"></i>
          {{/if}}
        </span>
      </td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<script>
  $(function () {
    $('.sort-text').on('dblclick', function () { // [!code ++]
      const userId = $(this).data('id'); // [!code ++]
      $(this).addClass('d-none'); // [!code ++]
      $(`.sort-input[data-id="${userId}"]`).removeClass('d-none').focus(); // [!code ++]
    }); // [!code ++]

    $('.sort-input').on('blur', function () { // [!code ++]
      const userId = $(this).data('id'); // [!code ++]
      const newSort = $(this).val(); // [!code ++]
      $(this).addClass('d-none'); // [!code ++]
      $(`.sort-text[data-id="${userId}"]`).removeClass('d-none').text(newSort); // [!code ++]
      $.ajax({ // [!code ++]
        url: `/admin/user/${userId}`, // [!code ++]
        type: 'PUT', // [!code ++]
        contentType: 'application/json', // [!code ++]
        headers: { // [!code ++]
          'accept': 'application/json' // [!code ++]
        }, // [!code ++]
        data: JSON.stringify({ sort: newSort }), // [!code ++]
        success: function (response) { // [!code ++]
          if (response.success) { // [!code ++]
            $(`.sort-text[data-id="${userId}"]`).text(newSort); // [!code ++]
          } // [!code ++]
        } // [!code ++]
      }); // [!code ++]
    }); // [!code ++]

    $('.sort-input').on('keypress', function (e) { // [!code ++]
      if (e.which == 13) { // [!code ++]
        $(this).blur(); // [!code ++]
      } // [!code ++]
    }); // [!code ++]
    
    $('.status-toggle').on('click', function () {
      const $this = $(this);
      const userId = $this.data('id');
      const currentStatus = $this.data('status');
      const newStatus = currentStatus === 1 ? 0 : 1;
      $.ajax({
        url: `/admin/user/${userId}`,
        type: 'PUT',
        contentType: 'application/json',
        headers: {
          'accept': 'application/json'
        },
        data: JSON.stringify({ status: newStatus }),
        success: function (response) {
          if (response.success) {
            $this.data('status', newStatus);
            $this.html(`<i class="bi ${newStatus ? "bi-check-circle-fill" : "bi-x-circle-fill"} ${newStatus ? "text-success" : "text-danger"}"></i>`);
          }
        },
        error: function (error) {
          const { responseJSON } = error;
          alert(responseJSON.message);
        }
      });
    });
  });
  function deleteUser(id) {
    if (confirm('确定要删除该用户吗？')) {
      $.ajax({
        url: '/admin/user/' + id,
        type: 'DELETE',
        success: function (res) {
          if (res.success) {
            window.location.reload()
          }
        }
      })
    }
  }
</script>
```

## 搜索

### user.controller

```ts
import { Body, Controller, Delete, Get, NotFoundException, Query, Param, ParseIntPipe, Headers, Post, Put, Redirect, Render, Res, UseFilters } from '@nestjs/common'; // [!code ++]
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import type { Response } from 'express';

@ApiTags('admin/user')
@UseFilters(AdminExceptionFilter)
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll(@Query('search') search: string = '') { // [!code ++]
    const users = await this.userService.findAll(search); // [!code ++]
    return { users };
  }
}
```

### user.service

```ts
import { Injectable } from '@nestjs/common';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Like, Repository } from 'typeorm'; // [!code ++]

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>
  ) {
    super(userRepository);
  }

  async findAll(search: string = ''): Promise<User[]> { // [!code ++]
    const where = search ? [ // [!code ++]
      { username: Like(`%${search}%`) }, // [!code ++]
      { email: Like(`%${search}%`) } // [!code ++]
    ] : {}; // [!code ++]

    const users = await this.userRepository.find({ // [!code ++]
      where // [!code ++]
    }); // [!code ++]
    return users; // [!code ++]
  } // [!code ++]
}

```

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<form method="GET" action="/admin/user" class="mb-3"> // [!code ++]
  <div class="input-group"> // [!code ++]
    <input type="text" name="search" class="form-control" placeholder="搜索用户名或邮箱" value="{{search}}"> // [!code ++]
    <button class="btn btn-outline-secondary" type="submit">搜索</button> // [!code ++]
  </div> // [!code ++]
</form> // [!code ++]
<table class="table">
  <thead>
    <tr>
      <th>序号</th>
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>
        <span class="sort-text" data-id="{{this.id}}">{{this.sort}}</span>
        <input type="number" class="form-control sort-input d-none" style="width:80px" data-id="{{this.id}}"
          value="{{this.sort}}">
      </td>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}">
          {{#if this.status}}
          <i class="bi bi-check-circle-fill text-success"></i>
          {{else}}
          <i class="bi bi-x-circle-fill text-danger"></i>
          {{/if}}
        </span>
      </td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
```

## 分页

### helpers

`eq`

```ts
export function eq(a: any, b: any) {
  return a === b;
}
```

`dec`

```ts
export function dec(value: number | string) {
  return Number(value) - 1;
}
```

`inc`

```ts
export function dec(value: number | string) {
  return Number(value) + 1;
}
```

`range`

```ts
export function range(start: number, end: number) {
  let result: number[] = [];
  for (let i = start; i <= end; i++) {
      result.push(i);
  }
  return result;
}
```

`index`

```ts
export * from './eq';
export * from './inc';
export * from './dec';
export * from './range';
```

`src/main.ts`

```ts
import { NestFactory } from '@nestjs/core';
import session from 'express-session';
import cookieParser from 'cookie-parser';
import { join } from 'node:path';
import { engine } from 'express-handlebars';
import { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { ValidationPipe } from '@nestjs/common';
import { useContainer } from 'class-validator';
import * as helpers from './share/helpers'; // [!code ++]

async function bootstrap() {
  // 使用 NestFactory 创建一个 NestExpressApplication 实例
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  // 使用 useContainer 配置依赖注入容器 让自定义校验器可以支持依赖注入
  useContainer(app.select(AppModule), { fallbackOnErrors: true });
  // 配置静态资源目录
  app.useStaticAssets(join(__dirname, '..', 'public'));
  // 设置视图文件的基本目录
  app.setBaseViewsDir(join(__dirname, '..', 'views'));
  // 设置视图引擎为 hbs（Handlebars）
  app.set('view engine', 'hbs');
  // 配置 Handlebars 引擎
  app.engine('hbs', engine({
    // 设置文件扩展名为 .hbs
    extname: '.hbs',
    helpers, // [!code ++]
    // 配置运行时选项
    runtimeOptions: {
      // 允许默认情况下访问原型属性
      allowProtoPropertiesByDefault: true,
      // 允许默认情况下访问原型方法
      allowProtoMethodsByDefault: true,
    },
  }));
  // 配置 cookie 解析器
  app.use(cookieParser());
  // 配置 session
  app.use(
    session({
      secret: 'secret-key',
      resave: true, // 是否每次都重新保存
      saveUninitialized: true, // 是否保存未初始化的会话
      cookie: {
        maxAge: 1000 * 60 * 60 * 24 * 7, // 7天
      },
    }),
  );
  // 配置全局管道
  app.useGlobalPipes(new ValidationPipe({ transform: true }));
  // 配置 Swagger
  const config = new DocumentBuilder()
    // 设置标题
    .setTitle('CMS API')
    // 设置描述
    .setDescription('CMS API Description')
    // 设置版本
    .setVersion('1.0')
    // 设置标签
    .addTag('CMS')
    // 设置Cookie认证
    .addCookieAuth('connect.sid')
    // 设置Bearer认证
    .addBearerAuth({ type: 'http', scheme: 'bearer' })
    // 构建配置
    .build();
  // 使用配置对象创建 Swagger 文档
  const document = SwaggerModule.createDocument(app, config);
  // 设置 Swagger 模块的路径和文档对象，将 Swagger UI 绑定到 '/api-doc' 路径上
  SwaggerModule.setup('api-doc', app, document);
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

### parse-optional-int.pipe

```ts
import { Injectable, PipeTransform, ArgumentMetadata, BadRequestException } from '@nestjs/common';

/**
 * 解析可选的整数参数
 * 如果参数为空（undefined、null 或 ''），返回默认值
 * 如果参数不是有效整数，则抛出 400 错误
 * 否则返回解析后的整数
 */
@Injectable()
export class ParseOptionalIntPipe implements PipeTransform<string, number> {
  constructor(private readonly defaultValue: number) { }

  transform(value: string, metadata: ArgumentMetadata): number {
    // 1. 如果参数为空（undefined、null 或 ''），返回默认值
    if (!value) {
      return this.defaultValue;
    }

    // 2. 尝试解析为整数
    const parsedValue = parseInt(value, 10);

    // 3. 如果不是有效整数，则抛出 400 错误
    if (isNaN(parsedValue)) {
      throw new BadRequestException(`Validation failed. "${value}" is not an integer.`);
    }

    // 4. 否则返回解析后的整数
    return parsedValue;
  }
}
```

### user.controller

```ts
import { Body, Controller, Delete, Get, NotFoundException, Query, Param, ParseIntPipe, Headers, Post, Put, Redirect, Render, Res, UseFilters } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import type { Response } from 'express';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe'; // [!code ++]

@ApiTags('admin/user')
@UseFilters(AdminExceptionFilter)
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) { }

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number) { // [!code ++]
    const { users, total } = await this.userService.findAllWithPagination(page, limit, search); // [!code ++]
    const pageCount = Math.ceil(total / limit); // [!code ++]
    return { users, search, page, limit, pageCount }; // [!code ++]
  } // [!code ++]
}
```

### user.service

```ts
import { Injectable } from '@nestjs/common';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Like, Repository } from 'typeorm';

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>
  ) {
    super(userRepository);
  }

  async findAll(search: string = ''): Promise<User[]> {
    const where = search ? [
      { username: Like(`%${search}%`) },
      { email: Like(`%${search}%`) }
    ] : {};

    const users = await this.userRepository.find({
      where
    });
    return users;
  }

  async findAllWithPagination(page: number = 1, limit: number = 10, search: string = ''): Promise<{ users: User[], total: number }> { // [!code ++]
    const where = search ? [ // [!code ++]
      { username: Like(`%${search}%`) }, // [!code ++]
      { email: Like(`%${search}%`) } // [!code ++]
    ] : {}; // [!code ++]

    const [users, total] = await this.userRepository.findAndCount({ // [!code ++]
      where, // [!code ++]
      skip: (page - 1) * limit, // [!code ++]
      take: limit, // [!code ++]
    }); // [!code ++]
    return { users, total }; // [!code ++]
  } // [!code ++]
}
```

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<form method="GET" action="/admin/user" class="mb-3">
  <div class="input-group">
    <input type="text" name="search" class="form-control" placeholder="搜索用户名或邮箱" value="{{search}}">
    <button class="btn btn-outline-secondary" type="submit">搜索</button>
  </div>
</form>
<table class="table">
  <thead>
    <tr>
      <th>序号</th>
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>
        <span class="sort-text" data-id="{{this.id}}">{{this.sort}}</span>
        <input type="number" class="form-control sort-input d-none" style="width:80px" data-id="{{this.id}}"
          value="{{this.sort}}">
      </td>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}">
          {{#if this.status}}
          <i class="bi bi-check-circle-fill text-success"></i>
          {{else}}
          <i class="bi bi-x-circle-fill text-danger"></i>
          {{/if}}
        </span>
      </td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<nav> // [!code ++]
  <ul class="pagination"> // [!code ++]
    <li class="page-item {{#if (eq page 1)}}disabled{{/if}}"> // [!code ++]
      <a class="page-link" href="?page={{dec page}}&search={{search}}&limit={{limit}}">上一页</a> // [!code ++]
    </li> // [!code ++]
    {{#each (range 1 pageCount)}} // [!code ++]
    <li class="page-item {{#if (eq this ../page)}}active{{/if}}"> // [!code ++]
      <a class="page-link" href="?page={{this}}&search={{../search}}&limit={{../limit}}">{{this}}</a> // [!code ++]
    </li> // [!code ++]
    {{/each}} // [!code ++]
    <li class="page-item {{#if (eq page pageCount)}}disabled{{/if}}"> // [!code ++]
      <a class="page-link" href="?page={{inc page}}&search={{search}}&limit={{limit}}">下一页</a> // [!code ++]
    </li> // [!code ++]
    <li class="page-item"> // [!code ++]
      <form method="GET" action="/admin/user" class="d-inline-block ms-3"> // [!code ++]
        <input type="hidden" name="search" value="{{search}}"> // [!code ++]
        <input type="hidden" name="page" value="{{page}}"> // [!code ++]
        <div class="input-group"> // [!code ++]
          <input type="number" name="limit" class="form-control" placeholder="每页条数" value="{{limit}}" min="1"> // [!code ++]
          <button class="btn btn-outline-secondary" type="submit">设置</button> // [!code ++]
        </div> // [!code ++]
      </form> // [!code ++]
    </li> // [!code ++]
  </ul> // [!code ++]
</nav> // [!code ++]
```

## 角色管理页面

用户管理和角色管理页面几乎是一样的，就会有很多重复代码，可以可以使用代码生成器生成，自己也可以实现代码生成器

将项目资源下载到本地

```bash
npm i cms-resource
```

生成角色管理的代码

```bash
nest g cms-resource role 角色 --collection=./node_modules/cms-resource
```

## 资源管理页面

可以自己根据项目实现一个生成器，比如 code 中的 `nest/cms-generator` 项目来生成资源管理的页面

进入`nest/cms-generator`先运行下build

```bash
npm run build
```

在 cms 目录下执行命令进行生成，即可生成页面

```bash
nest g generateList access 资源 --collection=../cms-generator
```

## 给用户分配角色

### user.entity

```ts
import { ApiHideProperty, ApiProperty } from '@nestjs/swagger';
import { Exclude, Transform } from 'class-transformer';
import { Entity, Column, PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn, ManyToMany, JoinTable } from 'typeorm';
import { Role } from './role.entity'; // [!code ++]

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  @ApiProperty({ description: '用户ID', example: 1 })
  id: number;

  @Column({ length: 50, unique: true })
  @ApiProperty({ description: '用户名', example: 'admin' })
  username: string;

  @Column()
  @Exclude() // 在序列化时排除密码字段，不返回给前端
  @ApiHideProperty() // 隐藏密码字段，不在Swagger文档中显示
  password: string;

  @Column({ length: 15, nullable: true })
  @ApiProperty({ description: '手机号', example: '13124567890', format: '手机号码会被部分隐藏' })
  @Transform(({ value }) => value ? value.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') : value)
  mobile: string;

  @Column({ length: 100, nullable: true })
  @ApiProperty({ description: '邮箱', example: 'admin@example.com' })
  email: string;

  @Column({ default: 1 })
  @ApiProperty({ description: '状态', example: 1, enum: [1, 2] })
  status: number;

  @ManyToMany(() => Role) // [!code ++]
  @JoinTable() // [!code ++]
  roles: Role[]; // [!code ++]

  @Column({ default: false })
  @ApiProperty({ description: '是否超级管理员', example: false })
  is_super: boolean;

  @Column({ default: 100 })
  @ApiProperty({ description: '排序', example: 100 })
  sort: number;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '创建时间', example: '2021-01-01 00:00:00' })
  @CreateDateColumn()
  createdAt: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP', onUpdate: 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '更新时间', example: '2021-01-01 00:00:00' })
  @UpdateDateColumn()
  updatedAt: Date;
}
```

### user.dto

```ts
import { IsOptional, IsString, Validate } from "class-validator";
import { StartsWithConstraint, IsUsernameUniqueConstraint } from "../validators/user-validators";
import { ApiProperty, ApiPropertyOptional, PartialType } from "@nestjs/swagger"
import { IsOptionalString, IsOptionalEmail, IsOptionalNumber, IsOptionalBoolean } from "../decorators/alidation-and-transform.decorators";

export class CreateUserDto {
  @ApiProperty({ description: '用户名，必须唯一且以指定前缀开头', example: 'user_john_doe' })
  @IsString()
  @Validate(StartsWithConstraint, ['user_'], {
    message: `用户名必须以 "user_" 开头`,
  })
  @Validate(IsUsernameUniqueConstraint, { message: '用户名已存在' })
  // @StartsWith('user_', { message: '用户名必须以 "user_" 开头' })
  // @IsUsernameUnique({ message: '用户名已存在' })
  username: string;

  @ApiProperty({ description: '密码', example: 'securePassword123' })
  @IsString()
  password: string;

  @ApiPropertyOptional({ description: '手机号', example: '13124567890' })
  @IsOptionalString()
  mobile?: string;

  @ApiPropertyOptional({ description: '邮箱地址', example: 'john.doe@example.com' })
  @IsOptionalEmail()
  email?: string;

  @ApiPropertyOptional({ description: '用户状态', example: 1 })
  @IsOptionalNumber()
  status?: number;

  @ApiPropertyOptional({ description: '是否为超级管理员', example: true })
  @IsOptionalBoolean()
  is_super?: boolean;
}

export class UpdateUserDto extends PartialType(CreateUserDto) {
  @ApiProperty({ description: '用户ID', example: 1 })
  @IsOptionalNumber()
  id: number;
  @IsString() // [!code ++]
  @IsOptional() // [!code ++]
  @ApiProperty({ description: '用户名', example: 'nick' }) // [!code ++]
  username: string; // [!code ++]
  @ApiProperty({ description: '密码', example: '666666' }) // [!code ++]
  @IsOptional() // [!code ++]
  password?: string; // [!code ++]
}

export class UpdateUserRolesDto { // [!code ++]
  readonly roleIds: number[]; // [!code ++]
} // [!code ++]
```

### user.service

```ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { In, Like, Repository } from 'typeorm';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { Role } from '../entities/role.entity'; // [!code ++]
import { UpdateUserRolesDto } from '../dtos/user.dto'; // [!code ++]

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>,
    @InjectRepository(Role) // [!code ++]
    protected roleRepository: Repository<Role> // [!code ++]
  ) {
    super(userRepository);
  }

  async findAll(search: string = ''): Promise<User[]> {
    const where = search ? [
      { username: Like(`%${search}%`) },
      { email: Like(`%${search}%`) }
    ] : {};

    const users = await this.userRepository.find({
      where
    });
    return users;
  }

  async findAllWithPagination(page: number = 1, limit: number = 10, search: string = ''): Promise<{ users: User[], total: number }> {
    const where = search ? [
      { username: Like(`%${search}%`) },
      { email: Like(`%${search}%`) }
    ] : {};

    const [users, total] = await this.userRepository.findAndCount({
      where,
      skip: (page - 1) * limit,
      take: limit,
    });
    return { users, total };
  }

  async updateRoles(id: number, updateUserRolesDto: UpdateUserRolesDto) { // [!code ++]
    const user = await this.repository.findOneBy({ id }); // [!code ++]
    if (!user) throw new Error('User not found'); // [!code ++]
    user.roles = await this.roleRepository.findBy({ id: In(updateUserRolesDto.roleIds) }); // [!code ++]
    await this.repository.update(id, user); // [!code ++]
  } // [!code ++]
}
```

### user.controller

```ts
import { Body, Controller, Delete, Get, NotFoundException, Query, Param, ParseIntPipe, Headers, Post, Put, Redirect, Render, Res, UseFilters } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto, UpdateUserRolesDto } from 'src/share/dtos/user.dto'; // [!code ++]
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import type { Response } from 'express';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import { RoleService } from 'src/share/services/role.service'; // [!code ++]

@ApiTags('admin/user')
@UseFilters(AdminExceptionFilter)
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService,
    private readonly roleService: RoleService // [!code ++]
  ) { }

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number) {
    const { users, total } = await this.userService.findAllWithPagination(page, limit, search);
    const pageCount = Math.ceil(total / limit);
    const roles = await this.roleService.findAll(); // [!code ++]
    return { users, search, page, limit, pageCount, roles }; // [!code ++]
  }
    
  // ...
   
  @Get(':id')
  @ApiOperation({ summary: '获取用户详情(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户详情' })
  async findOne(@Param('id', ParseIntPipe) id: number, @Res() res: Response, @Headers('accept') accept: string) { // [!code ++]
    const user = await this.userService.findOne({ where: { id }, relations: ['roles'] }); // [!code ++]
    if (!user) throw new HttpException('User not Found', 404) // [!code ++]
    if (accept === 'application/json') { // [!code ++]
      return res.json(user); // [!code ++]
    } else { // [!code ++]
      res.render('user/user-detail', { user }); // [!code ++]
    } // [!code ++]
  }
  @Put(':id/roles') // [!code ++]
  @ApiOperation({ summary: '更新用户角色(管理后台)' }) // [!code ++]
  @ApiResponse({ status: 200, description: '成功返回更新用户角色页面' }) // [!code ++]
  async updateRoles(@Param('id', ParseIntPipe) id: number, @Body() updateUserRolesDto: UpdateUserRolesDto) { // [!code ++]
    await this.userService.updateRoles(id, updateUserRolesDto); // [!code ++]
    return { success: true }; // [!code ++]
  } // [!code ++]
}
```

### user-list.hbs

```handlebars
<h1>用户列表</h1>
<form method="GET" action="/admin/user" class="mb-3">
  <div class="input-group">
    <input type="text" name="search" class="form-control" placeholder="搜索用户名或邮箱" value="{{search}}">
    <button class="btn btn-outline-secondary" type="submit">搜索</button>
  </div>
</form>
<table class="table">
  <thead>
    <tr>
      <th>序号</th>
      <th>用户名</th>
      <th>邮箱</th>
      <th>状态</th>
      <th>操作</th>
    </tr>
  </thead>
  <tbody>
    {{#each users}}
    <tr>
      <td>
        <span class="sort-text" data-id="{{this.id}}">{{this.sort}}</span>
        <input type="number" class="form-control sort-input d-none" style="width:80px" data-id="{{this.id}}"
          value="{{this.sort}}">
      </td>
      <td>{{this.username}}</td>
      <td>{{this.email}}</td>
      <td>
        <span class="status-toggle" data-id="{{this.id}}" data-status="{{this.status}}">
          {{#if this.status}}
          <i class="bi bi-check-circle-fill text-success"></i>
          {{else}}
          <i class="bi bi-x-circle-fill text-danger"></i>
          {{/if}}
        </span>
      </td>
      <td>
        <a href="/admin/user/{{this.id}}">查看</a>
        <a href="/admin/user/edit/{{this.id}}">编辑</a>
        <a href="" class="delete-user" onclick="deleteUser({{this.id}})">删除</a>
        <button class="btn btn-info btn-sm" onclick="assignRoles({{this.id}})">分配角色</button>
      </td>
    </tr>
    {{/each}}
  </tbody>
</table>
<nav>
  <ul class="pagination">
    <li class="page-item {{#if (eq page 1)}}disabled{{/if}}">
      <a class="page-link" href="?page={{dec page}}&search={{search}}&limit={{limit}}">上一页</a>
    </li>
    {{#each (range 1 pageCount)}}
    <li class="page-item {{#if (eq this ../page)}}active{{/if}}">
      <a class="page-link" href="?page={{this}}&search={{../search}}&limit={{../limit}}">{{this}}</a>
    </li>
    {{/each}}
    <li class="page-item {{#if (eq page pageCount)}}disabled{{/if}}">
      <a class="page-link" href="?page={{inc page}}&search={{search}}&limit={{limit}}">下一页</a>
    </li>
    <li class="page-item">
      <form method="GET" action="/admin/user" class="d-inline-block ms-3">
        <input type="hidden" name="search" value="{{search}}">
        <input type="hidden" name="page" value="{{page}}">
        <div class="input-group">
          <input type="number" name="limit" class="form-control" placeholder="每页条数" value="{{limit}}" min="1">
          <button class="btn btn-outline-secondary" type="submit">设置</button>
        </div>
      </form>
    </li>
  </ul>
</nav>
<div class="modal fade" id="roleModal" tabindex="-1">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="roleModalLabel">分配角色</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
      </div>
      <div class="modal-body">
        <form id="roleForm">
          {{#each roles}}
          <div class="form-check">
            <input class="form-check-input" type="checkbox" value="{{this.id}}" id="role{{this.id}}">
            <label class="form-check-label" for="role{{this.id}}">
              {{this.name}}
            </label>
          </div>
          {{/each}}
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">关闭</button>
        <button type="button" class="btn btn-primary" id="saveRoles">保存</button>
      </div>
    </div>
  </div>
</div>
<script>
  let selectedUserId;
  function assignRoles(userId) {
    selectedUserId = userId;
    $.ajax({
      url: `/admin/user/${userId}`,
      type: 'GET',
      headers: {
        'accept': 'application/json'
      },
      success: function (user) {
        const roles = user.roles.map(role => role.id);
        $('#roleForm input[type="checkbox"]').each(function () {
          $(this).prop('checked', roles.includes(parseInt($(this).val())));
        });
        $('#roleModal').modal('show');
      }
    });

  }
  $('#saveRoles').on('click', function () {
    const roleIds = $('#roleForm input[type="checkbox"]:checked').map(function () {
      return $(this).val();
    }).get();
    $.ajax({
      url: `/admin/user/${selectedUserId}/roles`,
      type: 'PUT',
      headers: {
        'accept': 'application/json'
      },
      contentType: 'application/json',
      data: JSON.stringify({ roleIds }),
      success: function (response) {
        $('#roleModal').modal('hide');
        location.reload();
      },
      error: function (error) {
        const { responseJSON } = error;
        alert(responseJSON.message);
      }
    });
  });
  $(function () {
    $('.sort-text').on('dblclick', function () {
      const userId = $(this).data('id');
      $(this).addClass('d-none');
      $(`.sort-input[data-id="${userId}"]`).removeClass('d-none').focus();
    });

    $('.sort-input').on('blur', function () {
      const userId = $(this).data('id');
      const newSort = $(this).val();
      $(this).addClass('d-none');
      $(`.sort-text[data-id="${userId}"]`).removeClass('d-none').text(newSort);
      $.ajax({
        url: `/admin/user/${userId}`,
        type: 'PUT',
        contentType: 'application/json',
        headers: {
          'accept': 'application/json'
        },
        data: JSON.stringify({ sort: newSort }),
        success: function (response) {
          if (response.success) {
            $(`.sort-text[data-id="${userId}"]`).text(newSort);
          }
        }
      });
    });

    $('.sort-input').on('keypress', function (e) {
      if (e.which == 13) {
        $(this).blur();
      }
    });

    $('.status-toggle').on('click', function () {
      const $this = $(this);
      const userId = $this.data('id');
      const currentStatus = $this.data('status');
      const newStatus = currentStatus === 1 ? 0 : 1;
      $.ajax({
        url: `/admin/user/${userId}`,
        type: 'PUT',
        contentType: 'application/json',
        headers: {
          'accept': 'application/json'
        },
        data: JSON.stringify({ status: newStatus }),
        success: function (response) {
          if (response.success) {
            $this.data('status', newStatus);
            $this.html(`<i class="bi ${newStatus ? "bi-check-circle-fill" : "bi-x-circle-fill"} ${newStatus ? "text-success" : "text-danger"}"></i>`);
          }
        },
        error: function (error) {
          const { responseJSON } = error;
          alert(responseJSON.message);
        }
      });
    });
  });
  function deleteUser(id) {
    if (confirm('确定要删除该用户吗？')) {
      $.ajax({
        url: '/admin/user/' + id,
        type: 'DELETE',
        success: function (res) {
          if (res.success) {
            const params = new URLSearchParams(window.location.search);
            params.delete('page');
            params.append('page', 1)
            const query = params.toString();
            window.location = window.location.pathname + '?' + query;
          }
        }
      })
    }
  }
</script>
```

## 给角色分配权限、文章管理、分类管理、标签管理

内容类似，可以直接查看code中的相关代码。

## 富文本编辑器

导入ckeditor5的css

`main.hbs`

```handlebars
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>CMS后台管理页面</title>
  <link href="/css/bootstrap.min.css" rel="stylesheet" />
  <link href="/css/bootstrap-icons.min.css" rel="stylesheet">
  <link href="/css/ckeditor5.css" rel="stylesheet" /> // [!code ++]
  <script src="/js/jquery.min.js"></script>
  <script src="/js/bootstrap.bundle.min.js"></script>
</head>

<body>
  {{> header}}
  <div class="container-fluid">
    <div class="row">
      {{> sidebar}}
      <!-- 右侧管理页面 -->
      <div class="col-md-9 col-lg-10">
        <div class="container mt-4">
          {{{body}}}
        </div>
      </div>
    </div>
  </div>
</body>

</html>
```

`article-form.hbs`

```handlebars
<h1>{{#if article.id}}编辑文章{{else}}添加文章{{/if}}</h1>
<form action="/admin/articles{{#if article.id}}/{{article.id}}{{/if}}" method="POST" id="articleForm"> // [!code ++]
  {{#if article.id}}<input type="hidden" name="_method" value="PUT">{{/if}}
  <div class="mb-3">
    <label for="title" class="form-label">标题</label>
    <input type="text" class="form-control" id="title" name="title" value="{{article.title}}">
  </div>
  <div class="mb-3">
    <label for="content" class="form-label">内容</label>
    <textarea class="form-control" id="content" name="content" rows="10">{{article.content}}</textarea> // [!code --]
    <div id="editor"> // [!code ++]
      {{{article.content}}} // [!code ++]
    </div> // [!code ++]
    <input type="hidden" name="content" id="contentInput">
  </div>
  <div class="mb-3">
    <label class="form-label">分类</label>
    <div id="categoryTree" class="border rounded p-3"></div>
  </div>
  <div class="mb-3">
    <label for="tags" class="form-label">标签</label>
    <div class="d-flex flex-wrap">
      {{#each tags}}
      <div class="form-check me-3 mb-2">
        <input class="form-check-input" type="checkbox" name="tagIds" value="{{this.id}}" {{#if (contains (mapToIds
          ../article.tags) this.id )}}checked{{/if}}>
        <label class="form-check-label">{{this.name}}</label>
      </div>
      {{/each}}
    </div>
  </div>
  <div class="mb-3">
    <label for="status" class="form-label">状态</label>
    <select class="form-control" id="status" name="status">
      <option value="1" {{#if article.status}}selected{{/if}}>激活</option>
      <option value="0" {{#unless article.status}}selected{{/unless}}>未激活</option>
    </select>
  </div>
  <button type="submit" class="btn btn-primary">保存</button>
</form>
<script type="importmap"> // [!code ++]
  { // [!code ++]
    "imports": { // [!code ++]
      "ckeditor5": "/js/ckeditor5.js" // [!code ++]
    } // [!code ++]
  } // [!code ++]
</script> // [!code ++]
<script type="module"> // [!code ++]
  import { // [!code ++]
    ClassicEditor, // [!code ++]
    Essentials, // [!code ++]
    Bold, // [!code ++]
    Italic, // [!code ++]
    Font, // [!code ++]
    Paragraph, // [!code ++]
    Image, // [!code ++]
    ImageToolbar, // [!code ++]
    ImageUpload, // [!code ++]
    ImageResize, // [!code ++]
    ImageStyle, // [!code ++]
    Plugin // [!code ++]
  } from 'ckeditor5'; // [!code ++]
  ClassicEditor // [!code ++]
    .create(document.querySelector('#editor'), { // [!code ++]
      plugins: [ // [!code ++]
        Essentials, // [!code ++]
        Bold, // [!code ++]
        Italic, // [!code ++]
        Font, // [!code ++]
        Paragraph, // [!code ++]
        Image, // [!code ++]
        ImageToolbar, // [!code ++]
        ImageStyle, // [!code ++]
        ImageResize, // [!code ++]
        ImageUpload // [!code ++]
      ], // [!code ++]
      image: { // [!code ++]
        toolbar: ['imageTextAlternative', 'imageStyle:side', 'resizeImage:50', 'resizeImage:75', 'resizeImage:original'] // [!code ++]
      }, // [!code ++]
      toolbar: { // [!code ++]
        items: [ // [!code ++]
          'undo', 'redo', '|', 'bold', 'italic', '|', // [!code ++]
          'fontSize', 'fontFamily', 'fontColor', 'fontBackgroundColor', '|', // [!code ++]
          'insertImage' // [!code ++]
        ] // [!code ++]
      } // [!code ++]
    }) // [!code ++]
    .then(editor => { // [!code ++]
      const form = document.getElementById('articleForm'); // [!code ++]
      const contentInput = document.getElementById('contentInput'); // [!code ++]
      form.addEventListener('submit', () => { // [!code ++]
        contentInput.value = editor.getData(); // [!code ++]
      }); // [!code ++]
    }) // [!code ++]
    .catch(error => { // [!code ++]
      console.error(error.stack); // [!code ++]
    }); // [!code ++]
</script> // [!code ++]
<script>
  const categoryTree = {{{ json categoryTree }}};
  const selectedCategoryIds = {{{ mapToIds article.categories }}};
  function renderCategoryTree(categoryes) {
    let html = '<ul class="list-unstyled">';
    categoryes.forEach(function (category) {
      html += `
           <li class="mb-2">
               <div class="d-flex align-items-center">
                   ${category.children?.length > 0 ? '<span class="toggle me-2 cursor-pointer"><i class="bi bi-folder-minus"></i></span>' : '<span class="me-4"></span>'}
                   <label class="form-check-label">
                       <input type="checkbox" class="form-check-input" name="categoryIds" value="${category.id}" ${selectedCategoryIds.includes(category.id) ? 'checked' : ''}>
                       ${category.name}
                   </label>
               </div>
               ${category.children?.length > 0 ? `<div class="children ms-4" >${renderCategoryTree(category.children)}</div>` : ''}
           </li>`;
    });
    html += '</ul>';
    return html;
  }
  $(function () {
    $('#categoryTree').html(renderCategoryTree(categoryTree));
    $('body').on('click', '.toggle', function () {
      const childrenContainer = $(this).parent().siblings('.children');
      if (childrenContainer.is(':visible')) {
        childrenContainer.hide();
        $(this).html('<i class="bi bi-folder-plus"></i>');
      } else {
        childrenContainer.show();
        $(this).html('<i class="bi bi-folder-minus"></i>');
      }
    });
  });
</script>
```

## 文件上传

```bash
npm i @nestjs/serve-static multer uuid
npm i @types/multer --save-dev
```

`src/global.d.ts`

定义 multer 类型

```ts
declare namespace Express {
  interface Multer {
    File: Express.Multer.File;
  }
}
```

设置静态资源目录

`app.module`

```ts
import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module';
import { ApiModule } from './api/api.module';
import { ShareModule } from './share/share.module';
import methodOverride from './share/middlewares/method-override';
import { ServeStaticModule } from '@nestjs/serve-static'; // [!code ++]
import * as path from 'path'; // [!code ++]

@Module({
  imports: [
    ServeStaticModule.forRoot({ // [!code ++]
      rootPath: path.join(__dirname, '..', 'uploads'), // [!code ++]
      serveRoot: '/uploads', // [!code ++]
    }), // [!code ++]
    ShareModule, AdminModule, ApiModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(methodOverride).forRoutes('*');
  }
}
```

### upload.controller

```ts
import { Controller, Get, Post, Query, UploadedFile, UseInterceptors } from '@nestjs/common';
// 导入文件上传拦截器
import { FileInterceptor } from '@nestjs/platform-express';
// 导入multer的磁盘存储配置
import { diskStorage } from 'multer';
// 使用Node内置的randomUUID生成唯一文件名，避免ESM/CJS兼容问题
import { randomUUID } from 'crypto';
// 导入Node.js路径处理模块
import * as path from 'path';

/**
 * 文件上传控制器
 * 负责处理管理后台的文件上传功能
 * 支持图片文件上传，包括jpg、jpeg、png、gif格式
 */
@Controller('admin')
export class UploadController {
  
  /**
   * 文件上传接口
   * POST /admin/upload
   * 
   * 功能说明：
   * 1. 接收客户端上传的文件
   * 2. 验证文件类型（仅支持图片格式）
   * 3. 生成唯一文件名避免冲突
   * 4. 将文件保存到服务器磁盘
   * 5. 返回文件访问URL
   * 
   * @param file 上传的文件对象，包含文件信息和元数据
   * @returns 返回包含文件访问URL的响应对象
   */
  @Post('upload')
  @UseInterceptors(FileInterceptor('upload', {
    // 配置文件存储方式为磁盘存储
    storage: diskStorage({
      // 设置文件保存目录为项目根目录下的uploads文件夹
      destination: './uploads',
      // 自定义文件名生成规则
      filename: (_req, file, callback) => {
        // 使用Node内置的randomUUID生成唯一标识符，保留原文件扩展名
        // 这样可以避免文件名冲突，同时保持文件类型信息
        const filename: string = randomUUID() + path.extname(file.originalname);
        callback(null, filename);
      }
    }),
    // 文件类型过滤器，只允许特定格式的图片文件
    fileFilter: (req, file, callback) => {
      // 使用正则表达式验证MIME类型
      // 只允许jpg、jpeg、png、gif格式的图片文件
      if (!file.mimetype.match(/\/(jpg|jpeg|png|gif)$/)) {
        // 如果文件类型不支持，返回错误信息
        return callback(new Error('不支持的文件类型'), false);
      }
      // 文件类型验证通过，允许上传
      callback(null, true);
    }
  }))
  async uploadFile(@UploadedFile() file: Express.Multer.File) {
    // 返回文件访问URL，客户端可以通过此URL访问上传的文件
    // URL格式：/uploads/生成的唯一文件名
    return { url: `/uploads/${file.filename}` };
  }
}
```

### admin.module

```ts
import { Module } from '@nestjs/common';
import { DashboardController } from './controllers/dashboard.controller';
import { UserController } from './controllers/user.controller';
import { AdminExceptionFilter } from './filters/admin-exception.filter';
import { RoleController } from "./controllers/role.controller";
import { AccessController } from "./controllers/access.controller";
import { ArticleController } from './controllers/article.controller';
import { CategoryController } from './controllers/category.controller';
import { TagController } from './controllers/tag.controller';
import { UploadController } from './controllers/upload.controller'; // [!code ++]

@Module({
  controllers: [
    DashboardController,
    UserController,
    RoleController,
    AccessController,
    ArticleController,
    CategoryController,
    TagController,
    UploadController // [!code ++]
  ],
  providers: [{
    provide: 'APP_FILTER',
    useClass: AdminExceptionFilter,
  }],
})
export class AdminModule { }

```

### article-detail.hbs

```handlebars
<h1>
  文章详情
</h1>
<div class="mb-3">
  <label class="form-label">标题:</label>
  <p class="form-control-plaintext">{{article.title}}</p>
</div>
<div class="mb-3"> // [!code ++]
  <label class="form-label">内容:</label> // [!code ++]
  <div class="form-control-plaintext"> // [!code ++]
    {{{article.content}}} // [!code ++]
  </div> // [!code ++]
</div> // [!code ++]
<div class="mb-3">
  <label class="form-label">分类:</label>
  <p class="form-control-plaintext">
    {{#each article.categories}}
    <span class="badge bg-secondary">{{this.name}}</span>
    {{/each}}
  </p>
</div>
<div class="mb-3">
  <label class="form-label">标签:</label>
  <p class="form-control-plaintext">
    {{#each article.tags}}
    <span class="badge bg-info text-dark">{{this.name}}</span>
    {{/each}}
  </p>
</div>
<a href="/admin/articles/{{article.id}}/edit" class="btn btn-warning btn-sm">修改</a>
<a href="/admin/articles" class="btn btn-secondary btn-sm">返回列表</a>
</div>
```

### article-form.hbs

```handlebars
<h1>{{#if article.id}}编辑文章{{else}}添加文章{{/if}}</h1>
<form action="/admin/articles{{#if article.id}}/{{article.id}}{{/if}}" method="POST" id="articleForm">
  {{#if article.id}}<input type="hidden" name="_method" value="PUT">{{/if}}
  <div class="mb-3">
    <label for="title" class="form-label">标题</label>
    <input type="text" class="form-control" id="title" name="title" value="{{article.title}}">
  </div>
  <div class="mb-3">
    <label for="content" class="form-label">内容</label>
    {{!-- <textarea class="form-control" id="content" name="content" rows="10">{{article.content}}</textarea> --}}
    <div id="editor">
      {{{article.content}}}
    </div>
    <input type="hidden" name="content" id="contentInput">
  </div>
  <div class="mb-3">
    <label class="form-label">分类</label>
    <div id="categoryTree" class="border rounded p-3"></div>
  </div>
  <div class="mb-3">
    <label for="tags" class="form-label">标签</label>
    <div class="d-flex flex-wrap">
      {{#each tags}}
      <div class="form-check me-3 mb-2">
        <input class="form-check-input" type="checkbox" name="tagIds" value="{{this.id}}" {{#if (contains (mapToIds
          ../article.tags) this.id )}}checked{{/if}}>
        <label class="form-check-label">{{this.name}}</label>
      </div>
      {{/each}}
    </div>
  </div>
  <div class="mb-3">
    <label for="status" class="form-label">状态</label>
    <select class="form-control" id="status" name="status">
      <option value="1" {{#if article.status}}selected{{/if}}>激活</option>
      <option value="0" {{#unless article.status}}selected{{/unless}}>未激活</option>
    </select>
  </div>
  <button type="submit" class="btn btn-primary">保存</button>
</form>
<script type="importmap">
  {
    "imports": {
      "ckeditor5": "/js/ckeditor5.js"
    }
  }
</script>
<script type="module">
  import {
    ClassicEditor,
    Essentials,
    Bold,
    Italic,
    Font,
    Paragraph,
    Image,
    ImageToolbar,
    ImageUpload,
    ImageResize,
    ImageStyle,
    Plugin,
    SimpleUploadAdapter // [!code ++]
  } from 'ckeditor5';
  ClassicEditor
    .create(document.querySelector('#editor'), {
      plugins: [
        Essentials,
        Bold,
        Italic,
        Font,
        Paragraph,
        Image,
        ImageToolbar,
        ImageStyle,
        ImageResize,
        ImageUpload,
        SimpleUploadAdapter // [!code ++]
      ],
      image: {
        toolbar: ['imageTextAlternative', 'imageStyle:side', 'resizeImage:50', 'resizeImage:75', 'resizeImage:original']
      },
      toolbar: {
        items: [
          'undo', 'redo', '|', 'bold', 'italic', '|',
          'fontSize', 'fontFamily', 'fontColor', 'fontBackgroundColor', '|',
          'insertImage'
        ]
      },
      simpleUpload: { // [!code ++]
        uploadUrl: '/admin/upload', // [!code ++]
        withCredentials: true, // [!code ++]
        headers: { // [!code ++]
          Authorization: 'Bearer <JSON Web Token>' // [!code ++]
        } // [!code ++]
      } // [!code ++]
    })
    .then(editor => {
      const form = document.getElementById('articleForm');
      const contentInput = document.getElementById('contentInput');
      form.addEventListener('submit', () => {
        contentInput.value = editor.getData();
      });
    })
    .catch(error => {
      console.error(error.stack);
    });
</script>
<script>
  const categoryTree = {{{ json categoryTree }}};
  const selectedCategoryIds = {{{ mapToIds article.categories }}};
  function renderCategoryTree(categoryes) {
    let html = '<ul class="list-unstyled">';
    categoryes.forEach(function (category) {
      html += `
           <li class="mb-2">
               <div class="d-flex align-items-center">
                   ${category.children?.length > 0 ? '<span class="toggle me-2 cursor-pointer"><i class="bi bi-folder-minus"></i></span>' : '<span class="me-4"></span>'}
                   <label class="form-check-label">
                       <input type="checkbox" class="form-check-input" name="categoryIds" value="${category.id}" ${selectedCategoryIds.includes(category.id) ? 'checked' : ''}>
                       ${category.name}
                   </label>
               </div>
               ${category.children?.length > 0 ? `<div class="children ms-4" >${renderCategoryTree(category.children)}</div>` : ''}
           </li>`;
    });
    html += '</ul>';
    return html;
  }
  $(function () {
    $('#categoryTree').html(renderCategoryTree(categoryTree));
    $('body').on('click', '.toggle', function () {
      const childrenContainer = $(this).parent().siblings('.children');
      if (childrenContainer.is(':visible')) {
        childrenContainer.hide();
        $(this).html('<i class="bi bi-folder-plus"></i>');
      } else {
        childrenContainer.show();
        $(this).html('<i class="bi bi-folder-minus"></i>');
      }
    });
  });
</script>
```

## 文件压缩

使用 sharp 进行图片压缩

```bash
npm i sharp
```

`upload.controller` 

```ts
async uploadFile(@UploadedFile() file: Express.Multer.File) {
    // 生成压缩后的文件名，扩展名为 .min.jpeg
    const filename = `${path.basename(file.filename, path.extname(file.filename))}.min.jpeg`;
    // 压缩后的文件路径
    const outputFilePath = path.resolve('./uploads', filename);
    // 先读入 buffer，避免 sharp 占用源文件句柄
    const buffer = await fs.promises.readFile(file.path);
    // 使用 sharp 压缩
    await sharp(buffer)
        .resize(800, 600, {
        fit: sharp.fit.inside,
        withoutEnlargement: true,
    })
        .toFormat('jpeg')
        .jpeg({ quality: 80 })
        .toFile(outputFilePath);
    // safe unlink（删除原始上传文件）
    try {
        await fs.promises.unlink(file.path);
    } catch (err) {
        console.warn(`⚠️ 删除原文件失败: ${file.path}`, err);
    }
    // 返回压缩后的 URL
    return { url: `/uploads/${filename}` };
}
```

## 对象存储COS

安装sdk

```bash
npm i cos-nodejs-sdk-v5 --save
```

配置环境变量

```
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=cms
MYSQL_USER=root
MYSQL_PASSWORD=password
COS_SECRET_ID=COS_SECRET_ID
COS_SECRET_KEY=COS_SECRET_KEY
COS_BUCKET=COS_BUCKET
COS_REGION=COS_REGION
```

### cos.service

```ts
// 导入 Injectable 装饰器，用于标记一个服务类
import { Injectable } from '@nestjs/common';
// 导入 ConfigService，用于获取配置文件中的配置信息
import { ConfigService } from '@nestjs/config';
// 导入 COS SDK
import COS from 'cos-nodejs-sdk-v5';
// 使用 Injectable 装饰器将 CosService 标记为可注入的服务
@Injectable()
export class CosService {
  // 定义一个私有变量，用于存储 COS 实例
  private cos: COS;
  // 构造函数，注入 ConfigService 以获取配置信息
  constructor(private readonly configService: ConfigService) {
    // 初始化 COS 实例，使用配置服务中的 SecretId 和 SecretKey
    this.cos = new COS({
      SecretId: this.configService.get('COS_SECRET_ID'),
      SecretKey: this.configService.get('COS_SECRET_KEY'),
    });
  }
  // 获取签名认证信息的方法，默认过期时间为 60 秒
  getAuth(key, expirationTime = 60) {
    // 从配置服务中获取 COS 存储桶名称和区域
    const bucket = this.configService.get('COS_BUCKET');
    const region = this.configService.get('COS_REGION');
    // 获取 COS 签名，用于 PUT 请求
    const sign = this.cos.getAuth({
      Method: 'put', // 请求方法为 PUT
      Key: key, // 文件的对象键（路径）
      Expires: expirationTime, // 签名的有效期
      Bucket: bucket, // 存储桶名称
      Region: region, // 存储桶所在区域
    });
    // 返回包含签名、键名、存储桶和区域的信息对象
    return {
      sign,
      key: key,
      bucket,
      region,
    };
  }
}
```

### share.module

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { RoleService } from './services/role.service';
import { AccessService } from "./services/access.service";
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';
import { Role } from './entities/role.entity';
import { Access } from "./entities/access.entity";
import { Article } from './entities/article.entity';
import { Category } from './entities/category.entity';
import { Tag } from './entities/tag.entity';
import { ArticleService } from './services/article.service';
import { CategoryService } from './services/category.service';
import { TagService } from './services/tag.service';
import { CosService } from './services/cos.service'; // [!code ++]

@Global()
@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] }), // [!code ++]
        TypeOrmModule.forFeature([User, Role, Access, Article, Category, Tag]),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            inject: [ConfigurationService],
            useFactory: (configService: ConfigurationService) => ({
                type: 'mysql',
                ...configService.mysqlConfig,
                entities: [User, Role, Access, Article, Category, Tag],
                synchronize: true,
                autoLoadEntities: true,
                logging: false
            }),
        }),
    ],
    providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService], // [!code ++]
    exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService], // [!code ++]
})
export class ShareModule {
}
```

### upload.controller

```ts
import fs from 'fs';
import { CosService } from '../../share/services/cos.service'; // [!code ++]

/**
 * 文件上传控制器
 * 负责处理管理后台的文件上传功能
 * 支持图片文件上传，包括jpg、jpeg、png、gif格式
 */
@Controller('admin')
export class UploadController {
  constructor(private readonly cosService: CosService) { } // [!code ++]

  @Get('cos-signature') // [!code ++]
  async getCosSignature(@Query('key') key: string) { // [!code ++]
    return this.cosService.getAuth(key, 60); // [!code ++]
  } // [!code ++]
}
```

### article-form

```handlebars
<script type="module">
  import {
    ClassicEditor,
    Essentials,
    Bold,
    Italic,
    Font,
    Paragraph,
    Image,
    ImageToolbar,
    ImageUpload,
    ImageResize,
    ImageStyle,
    Plugin,
    SimpleUploadAdapter
  } from 'ckeditor5';
  // 异步函数，用于获取文件上传的签名信息 // [!code ++]
  async function getSignature(key) { // [!code ++]
    // 发送请求到后台接口，获取 COS 上传的签名信息 // [!code ++]
    const response = await fetch(`/admin/cos-signature?key=${encodeURIComponent(key)}`); // [!code ++]
    // 返回签名信息的 JSON 数据 // [!code ++]
    return response.json(); // [!code ++]
  } // [!code ++]
  // 自定义的 COS 上传适配器类，用于将文件上传到腾讯云 COS // [!code ++]
  class COSUploadAdapter { // [!code ++]
    // 构造函数，接收一个文件加载器实例 // [!code ++]
    constructor(loader) { // [!code ++]
      this.loader = loader; // [!code ++]
    } // [!code ++]
    // 上传方法，负责将文件上传到 COS // [!code ++]
    async upload() { // [!code ++]
      // 等待加载器获取要上传的文件 // [!code ++]
      const file = await this.loader.file; // [!code ++]
      // 获取文件的上传签名信息 // [!code ++]
      const signature = await getSignature(file.name); // [!code ++]
      // 从签名信息中解构出存储桶、区域、文件键名和签名 // [!code ++]
      const { bucket, region, key, sign } = signature; // [!code ++]
      // 构造文件上传的 URL // [!code ++]
      const url = `https://${signature.bucket}.cos.${signature.region}.myqcloud.com/${signature.key}`; // [!code ++]
      // 发送 PUT 请求，将文件上传到 COS // [!code ++]
      return fetch(url, { // [!code ++]
        method: 'PUT', // 使用 PUT 方法上传文件 // [!code ++]
        headers: { Authorization: sign }, // 设置请求头，包含签名信息 // [!code ++]
        body: file // 请求体为文件本身 // [!code ++]
      }).then(response => { // [!code ++]
        // 上传成功后，返回包含文件 URL 的对象 // [!code ++]
        return { default: url }; // [!code ++]
      }); // [!code ++]
    } // [!code ++]
  } // [!code ++]
  // 插件类，用于将 COS 上传适配器集成到 CKEditor // [!code ++]
  class COSUploadAdapterPlugin extends Plugin { // [!code ++]
    // 静态方法，定义插件的依赖关系 // [!code ++]
    static get requires() { // [!code ++]
      return [ImageUpload]; // 依赖 ImageUpload 插件 // [!code ++]
    } // [!code ++]
    // 插件初始化方法 // [!code ++]
    init() { // [!code ++]
      // 获取文件库插件，并设置创建上传适配器的函数 // [!code ++]
      this.editor.plugins.get('FileRepository').createUploadAdapter = (loader) => new COSUploadAdapter(loader); // [!code ++]
    } // [!code ++]
  } // [!code ++]
  ClassicEditor
    .create(document.querySelector('#editor'), {
      plugins: [
        Essentials,
        Bold,
        Italic,
        Font,
        Paragraph,
        Image,
        ImageToolbar,
        ImageStyle,
        ImageResize,
        ImageUpload,
        SimpleUploadAdapter,
        COSUploadAdapterPlugin, // [!code ++]
      ],
      image: {
        toolbar: ['imageTextAlternative', 'imageStyle:side', 'resizeImage:50', 'resizeImage:75', 'resizeImage:original']
      },
      toolbar: {
        items: [
          'undo', 'redo', '|', 'bold', 'italic', '|',
          'fontSize', 'fontFamily', 'fontColor', 'fontBackgroundColor', '|',
          'insertImage'
        ]
      },
      simpleUpload: {
        uploadUrl: '/admin/upload',
        withCredentials: true,
        headers: {
          Authorization: 'Bearer <JSON Web Token>'
        }
      }
    })
    .then(editor => {
      const form = document.getElementById('articleForm');
      const contentInput = document.getElementById('contentInput');
      form.addEventListener('submit', () => {
        contentInput.value = editor.getData();
      });
    })
    .catch(error => {
      console.error(error.stack);
    });
</script>
```

## 发送审核通知

```bash
npm install @nestjs/event-emitter eventemitter2
```

### notification.service

```ts
import { Injectable } from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import { ArticleService } from './article.service';
import { UserService } from './user.service';

@Injectable()
export class NotificationService {
  constructor(
    private readonly articleService: ArticleService,
    private readonly userService: UserService,
  ) { }

  @OnEvent('article.submitted')
  async handleArticleSubmittedEvent(payload: { articleId: number }) {
    const article = await this.articleService.findOne({ where: { id: payload.articleId }, relations: ['categories', 'tags'] });
    const admin = await this.userService.findOne({ where: { is_super: true } });
    if (admin) {
      const subject = `文章审核请求: ${article?.title}`;
      const body = `有一篇新的文章需要审核，点击链接查看详情: http://localhost:3000/admin/articles/${payload.articleId}`;
      console.log(admin.email, subject, body);
    }
  }
}
```

### shard.module

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { RoleService } from './services/role.service';
import { AccessService } from "./services/access.service";
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';
import { Role } from './entities/role.entity';
import { Access } from "./entities/access.entity";
import { Article } from './entities/article.entity';
import { Category } from './entities/category.entity';
import { Tag } from './entities/tag.entity';
import { ArticleService } from './services/article.service';
import { CategoryService } from './services/category.service';
import { TagService } from './services/tag.service';
import { CosService } from './services/cos.service';
import { NotificationService } from './services/notification.service'; // [!code ++]

@Global()
@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] }),
        TypeOrmModule.forFeature([User, Role, Access, Article, Category, Tag]),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            inject: [ConfigurationService],
            useFactory: (configService: ConfigurationService) => ({
                type: 'mysql',
                ...configService.mysqlConfig,
                entities: [User, Role, Access, Article, Category, Tag],
                synchronize: true,
                autoLoadEntities: true,
                logging: false
            }),
        }),
    ],
    providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService], // [!code ++]
    exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService], // [!code ++]
})
export class ShareModule {
}
```

### app.module

```ts
import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module';
import { ApiModule } from './api/api.module';
import { ShareModule } from './share/share.module';
import methodOverride from './share/middlewares/method-override';
import { ServeStaticModule } from '@nestjs/serve-static';
import * as path from 'path';
import { EventEmitterModule } from '@nestjs/event-emitter'; // [!code ++]

@Module({
  imports: [
    // 配置 EventEmitterModule 模块 // [!code ++]
    EventEmitterModule.forRoot({ // [!code ++]
      // 启用通配符功能，允许使用通配符来订阅事件 // [!code ++]
      wildcard: true, // [!code ++]
      // 设置事件名的分隔符，这里使用 '.' 作为分隔符 // [!code ++]
      delimiter: '.', // [!code ++]
      // 将事件发射器设置为全局模块，所有模块都可以共享同一个事件发射器实例 // [!code ++]
      global: true // [!code ++]
    }), // [!code ++]
    ServeStaticModule.forRoot({
      rootPath: path.join(__dirname, '..', 'uploads'),
      serveRoot: '/uploads',
    }),
    ShareModule, AdminModule, ApiModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(methodOverride).forRoutes('*');
  }
}
```

### article.controller

```ts
import { EventEmitter2 } from '@nestjs/event-emitter'; // [!code ++]

@UseFilters(AdminExceptionFilter)
@Controller('admin/articles')
export class ArticleController {
  constructor(
    private readonly articleService: ArticleService,
    private readonly categoryService: CategoryService,
    private readonly tagService: TagService,
    private readonly eventEmitter: EventEmitter2,
  ) { }
    
  @Put(':id/submit')
  async submitForReview(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.PENDING } as UpdateArticleDto);
    this.eventEmitter.emit('article.submitted', { articleId: id }); // [!code ++]
    return { success: true };
  }
}
```

## 发送邮件

```bash
npm install nodemailer
```

使用qq邮箱进行发送，配置环境变量，设置正确的user和授权码

```
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=cms
MYSQL_USER=root
MYSQL_PASSWORD=password
COS_SECRET_ID=COS_SECRET_ID
COS_SECRET_KEY=COS_SECRET_KEY
COS_BUCKET=COS_BUCKET
COS_REGION=COS_REGION
SMTP_HOST=smtp.qq.com // [!code ++]
SMTP_PORT=465 // [!code ++]
SMTP_USER=xxx@qq.com // [!code ++]
SMTP_PASS=code // [!code ++]
```

### mail.service

```ts
import { Injectable } from '@nestjs/common';
import * as nodemailer from 'nodemailer';
import { ConfigService } from '@nestjs/config';
@Injectable()
export class MailService {
  private transporter;
  constructor(private readonly configService: ConfigService) {
    this.transporter = nodemailer.createTransport({
      host: configService.get('SMTP_HOST'),
      port: configService.get('SMTP_PORT'),
      secure: true,
      auth: {
        user: configService.get('SMTP_USER'),
        pass: configService.get('SMTP_PASS'),
      },
    });
  }

  async sendEmail(to: string, subject: string, body: string) {
    const mailOptions = {
      from: this.configService.get('SMTP_USER'), // 发件人
      to, // 收件人
      subject, // 主题
      text: body, // 邮件正文
    };
    try {
      const info = await this.transporter.sendMail(mailOptions);
      console.log(`邮件已发送: ${info.messageId}`);
    } catch (error) {
      console.error(`发送邮件失败: ${error.message}`);
    }
  }
}
```

### share.module

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { RoleService } from './services/role.service';
import { AccessService } from "./services/access.service";
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';
import { Role } from './entities/role.entity';
import { Access } from "./entities/access.entity";
import { Article } from './entities/article.entity';
import { Category } from './entities/category.entity';
import { Tag } from './entities/tag.entity';
import { ArticleService } from './services/article.service';
import { CategoryService } from './services/category.service';
import { TagService } from './services/tag.service';
import { CosService } from './services/cos.service';
import { NotificationService } from './services/notification.service';
import { MailService } from './services/mail.service'; // [!code ++]

@Global()
@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] }),
        TypeOrmModule.forFeature([User, Role, Access, Article, Category, Tag]),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            inject: [ConfigurationService],
            useFactory: (configService: ConfigurationService) => ({
                type: 'mysql',
                ...configService.mysqlConfig,
                entities: [User, Role, Access, Article, Category, Tag],
                synchronize: true,
                autoLoadEntities: true,
                logging: false
            }),
        }),
    ],
    providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService], // [!code ++]
    exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService], // [!code ++]
})
export class ShareModule {
}
```

### notification.service

```ts
import { Injectable } from '@nestjs/common';
import { OnEvent } from '@nestjs/event-emitter';
import { ArticleService } from './article.service';
import { UserService } from './user.service';
import { MailService } from './mail.service'; // [!code ++]

@Injectable()
export class NotificationService {
  constructor(
    private readonly articleService: ArticleService,
    private readonly userService: UserService,
    private readonly mailService: MailService, // [!code ++]
  ) { }

  @OnEvent('article.submitted')
  async handleArticleSubmittedEvent(payload: { articleId: number }) {
    const article = await this.articleService.findOne({ where: { id: payload.articleId }, relations: ['categories', 'tags'] });
    const admin = await this.userService.findOne({ where: { is_super: true } });
    if (admin) {
      const subject = `文章审核请求: ${article?.title}`;
      const body = `有一篇新的文章需要审核，点击链接查看详情: http://localhost:3000/admin/articles/${payload.articleId}`;
      console.log(admin.email, subject, body);
      this.mailService.sendEmail(admin.email, subject, body); // [!code ++]
    }
  }
}
```

## 导出为word

```bash
npm install html-to-docx
```

### word-export.service

```ts
import { Injectable } from '@nestjs/common';
import htmlToDocx from 'html-to-docx';

@Injectable()
export class WordExportService {
  async exportToWord(htmlContent: string): Promise<Buffer> {
    return await htmlToDocx(htmlContent);
  }
}
```

在 share.module.ts中引入word-export.service

### article.detail

```handlebars
<a href="/admin/articles/{{article.id}}/edit" class="btn btn-warning btn-sm">修改</a>
<a href="/admin/articles" class="btn btn-secondary btn-sm">返回列表</a>
<a href="/admin/articles/{{article.id}}/export-word" class="btn btn-primary btn-sm">导出为 Word</a> // [!code ++]
```

### article.controller

```ts
import { Controller, Get, Render, Post, Redirect, Body, UseFilters, Param, ParseIntPipe, Put, Delete, Headers, Res, Query, NotFoundException, StreamableFile, Header } from '@nestjs/common';
import { CreateArticleDto, UpdateArticleDto } from 'src/share/dtos/article.dto';
import { ArticleService } from 'src/share/services/article.service';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import { CategoryService } from 'src/share/services/category.service';
import { TagService } from 'src/share/services/tag.service';
import type { Response } from 'express';
import { ArticleStateEnum } from 'src/share/enums/article.enum';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { WordExportService } from 'src/share/services/word-export.service'; // [!code ++]

@UseFilters(AdminExceptionFilter)
@Controller('admin/articles')
export class ArticleController {
  constructor(
    private readonly articleService: ArticleService,
    private readonly categoryService: CategoryService,
    private readonly tagService: TagService,
    private readonly eventEmitter: EventEmitter2,
    private readonly wordExportService: WordExportService // [!code ++]
  ) { }

  @Get()
  @Render('article/article-list')
  async findAll(@Query('keyword') keyword: string = '',
    @Query('page', new ParseOptionalIntPipe(1)) page: number,
    @Query('limit', new ParseOptionalIntPipe(10)) limit: number) {
    const { articles, total } = await this.articleService.findAllWithPagination(page, limit, keyword);
    const pageCount = Math.ceil(total / limit);
    return { articles, keyword, page, limit, pageCount };
  }

  @Get('create')
  @Render('article/article-form')
  async createForm() {
    const categoryTree = await this.categoryService.findAll();
    const tags = await this.tagService.findAll();
    return { article: { categories: [], tags: [] }, categoryTree, tags };
  }

  @Post()
  @Redirect('/admin/articles')
  async create(@Body() createArticleDto: CreateArticleDto) {
    await this.articleService.create(createArticleDto);
    return { success: true }
  }

  @Get(':id/edit')
  @Render('article/article-form')
  async editForm(@Param('id', ParseIntPipe) id: number) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not Found');
    const categoryTree = await this.categoryService.findAll();
    const tags = await this.tagService.findAll();
    return { article, categoryTree, tags };
  }

  @Put(':id')
  async update(@Param('id', ParseIntPipe) id: number, @Body() updateArticleDto: UpdateArticleDto, @Res({ passthrough: true }) res: Response, @Headers('accept') accept: string) {
    await this.articleService.update(id, updateArticleDto);
    if (accept === 'application/json') {
      return { success: true };
    } else {
      return res.redirect(`/admin/articles`);
    }
  }

  @Delete(":id")
  async delete(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.delete(id);
    return { success: true }
  }

  @Get(':id')
  @Render('article/article-detail')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not Found');
    return { article };
  }

  @Put(':id/submit')
  async submitForReview(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.PENDING } as UpdateArticleDto);
    this.eventEmitter.emit('article.submitted', { articleId: id });
    return { success: true };
  }

  @Put(':id/approve')
  async approveArticle(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.PUBLISHED, rejectionReason: undefined } as UpdateArticleDto);
    return { success: true };
  }

  @Put(':id/reject')
  async rejectArticle(
    @Param('id', ParseIntPipe) id: number,
    @Body('rejectionReason') rejectionReason: string
  ) {
    await this.articleService.update(id, { state: ArticleStateEnum.REJECTED, rejectionReason } as UpdateArticleDto);
    return { success: true };
  }

  @Put(':id/withdraw')
  async withdrawArticle(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.WITHDRAWN } as UpdateArticleDto);
    return { success: true };
  }

  @Get(':id/export-word') // [!code ++]
  @Header('Content-Type', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document') // [!code ++]
  async exportWord(@Param('id', ParseIntPipe) id: number, @Res({ passthrough: true }) res: Response) { // [!code ++]
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] }); // [!code ++]
    if (!article) throw new NotFoundException('Article not found'); // [!code ++]

    const htmlContent = ` // [!code ++]
           <h1>${article.title}</h1> // [!code ++]
           <p><strong>状态:</strong> ${article.state}</p> // [!code ++]
           <p><strong>分类:</strong> ${article.categories.map(c => c.name).join(', ')}</p> // [!code ++]
           <p><strong>标签:</strong> ${article.tags.map(t => t.name).join(', ')}</p> // [!code ++]
           <hr/> // [!code ++]
           ${article.content} // [!code ++]
       `; // [!code ++]

    const buffer = await this.wordExportService.exportToWord(htmlContent); // [!code ++]
    res.setHeader('Content-Disposition', `attachment;  filename="${encodeURIComponent(article.title)}.docx"`); // [!code ++]
    return new StreamableFile(buffer); // [!code ++]
  } // [!code ++]
}
```

## 导出为ppt

```bash
npm install html-pptxgenjs pptxgenjs 
```

### ppt-export.service

```ts
// 导入 Injectable 装饰器，用于标记一个服务类
import { Injectable } from '@nestjs/common';
// 导入 PptxGenJS 库，用于生成 PPTX 文件
import PptxGenJS from 'pptxgenjs';
// 导入 html-pptxgenjs 库，用于将 HTML 转换为 PPTX 内容
import * as html2ppt from 'html-pptxgenjs';
// 使用 Injectable 装饰器将 PptExportService 标记为可注入的服务
@Injectable()
export class PptExportService {
  // 异步方法，用于将文章列表导出为 PPTX 文件
  async exportToPpt(articles: any[]) {
    // 创建一个新的 PPTX 对象
    const pptx = new (PptxGenJS as any)();
    // 遍历每篇文章，将其内容添加到 PPTX 幻灯片中
    for (const article of articles) {
      // 添加一个新的幻灯片到 PPTX
      const slide = pptx.addSlide();
      // 构建 HTML 内容，包含文章标题、状态、分类、标签和正文内容
      const htmlContent = `
                <h1>${article.title}</h1>
                <p><strong>状态:</strong> ${article.state}</p>
                <p><strong>分类:</strong> ${article.categories.map(c => c.name).join(', ')}</p>
                <p><strong>标签:</strong> ${article.tags.map(t => t.name).join(', ')}</p>
                <hr/>
                ${article.content}
            `;
      // 使用 html-pptxgenjs 将 HTML 内容转换为 PPTX 可用的文本项
      const items = html2ppt.htmlToPptxText(htmlContent);
      // 将生成的文本项添加到幻灯片中，设置其位置和大小
      slide.addText(items, { x: 0.5, y: 0.5, w: 9.5, h: 6, valign: 'top' });
    }
    // 将生成的 PPTX 文件以 nodebuffer 的形式输出
    return await pptx.write({ outputType: 'nodebuffer' });
  }
}
```

在 share.module.ts中引入ppt-export.service

### article.controller

```ts
@Get('export-ppt')
@Header('Content-Type', 'application/vnd.openxmlformats-officedocument.presentationml.presentation')
async exportPpt(@Query('keyword') keyword: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number, @Res({ passthrough: true }) res: Response) {
    const { articles } = await this.articleService.findAllWithPagination(page, limit, keyword);
    const buffer = await this.pptExportService.exportToPpt(articles);
    res.setHeader('Content-Disposition', 'attachment; filename=articles.pptx');
    return new StreamableFile(buffer);
}
```

### article-list

```handlebars
<a href="/admin/articles/create" class="btn btn-success mb-3">添加文章</a>
<button id="exportPptBtn" class="btn btn-warning btn-sm mb-3">导出PPT</button> // [!code ++]
<script>
  $(function () {
    $('#exportPptBtn').click(function () { // [!code ++]
      window.location.href = `/admin/articles/export-ppt?page={{page}}&keyword={{keyword}}&limit={{limit}}`; // [!code ++]
    }); // [!code ++]
  }
</script>
```

## 导出为excel

```bash
npm i exceljs
```

### excel-export.service

```ts
// 导入 Injectable 装饰器，用于将服务类标记为可注入的依赖
import { Injectable } from '@nestjs/common';
// 导入 ExcelJS 库，用于创建和操作 Excel 文件
import * as ExcelJS from 'exceljs';
@Injectable()
export class ExcelExportService {
  // 异步方法，用于将数据导出为 Excel 文件
  async exportAsExcel(data: any[], columns: { header: string, key: string, width: number }[]) {
    // 创建一个新的 Excel 工作簿
    const workbook = new ExcelJS.Workbook();
    // 添加一个新的工作表，并命名为 'Data'
    const worksheet = workbook.addWorksheet('Data');
    // 设置工作表的列，根据传入的列定义数组
    worksheet.columns = columns;
    // 遍历数据数组，将每一项数据作为一行添加到工作表中
    data.forEach(item => {
      worksheet.addRow(item);
    });
    // 将工作簿内容写入缓冲区，并返回该缓冲区（用于进一步处理或保存）
    return workbook.xlsx.writeBuffer();
  }
}
```

在 share.module.ts中引入excel-export.service

### article.controller

```ts
@Get('export-excel')
@Header('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
async exportExcel(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number, @Res({ passthrough: true }) res: Response) {
    const { articles } = await this.articleService.findAllWithPagination(page, limit, search);
    const data = articles.map(article => ({
        title: article.title,
        categories: article.categories.map(c => c.name).join(', '),
        tags: article.tags.map(t => t.name).join(', '),
        state: article.state,
        createdAt: article.createdAt,
    }));
    const columns = [
        { header: '标题', key: 'title', width: 30 },
        { header: '分类', key: 'categories', width: 20 },
        { header: '标签', key: 'tags', width: 20 },
        { header: '状态', key: 'state', width: 15 },
        { header: '创建时间', key: 'createdAt', width: 20 },
    ];
    const buffer = await this.excelExportService.exportAsExcel(data, columns);
    res.setHeader('Content-Disposition', `attachment; filename="articles.xlsx"`);
    return new StreamableFile(new Uint8Array(buffer));
}
```

### article-list

```handlebars
<button id="frontExportBtn" class="btn btn-info btn-sm mb-3">前台导出Excel</button>
<button id="backendExportBtn" class="btn btn-info btn-sm mb-3">后台导出Excel</button>
<script>
  $(function () {
    // 当用户点击 ID 为 backendExportBtn 的按钮时，触发事件处理函数
    $('#backendExportBtn').click(function () {
      // 将浏览器的窗口位置重定向到指定的导出 Excel 文件的 URL
      // URL 中包含当前页面、搜索条件和限制参数
      window.location.href = `/admin/articles/export-excel?page={{page}}&search={{search}}&limit={{limit}}`;
    });
    // 当用户点击 ID 为 frontExportBtn 的按钮时，触发事件处理函数
    $('#frontExportBtn').click(function () {
      // 初始化一个空的字符串，用于存储 CSV 文件内容
      let csvContent = '';
      // 遍历表格的表头行（thead 中的所有 tr），生成 CSV 文件的表头内容
      $('#articleTable thead tr').each(function () {
        let rowContent = ''; // 初始化空字符串，用于存储当前行的内容
        let cells = $(this).find('th'); // 获取当前行中的所有 th 单元格
        cells.each(function (index) {
          if (index < cells.length - 1) { // 遍历时排除最后一列（假设不需要导出）
            rowContent += $(this).text().trim() + ','; // 将单元格文本内容加入行内容，并用逗号分隔
          }
        });
        csvContent += rowContent.slice(0, -1) + '\n'; // 移除最后一个多余的逗号并添加换行符
      });
      // 遍历表格的主体行（tbody 中的所有 tr），生成 CSV 文件的表格内容
      $('#articleTable tbody tr').each(function () {
        let rowContent = ''; // 初始化空字符串，用于存储当前行的内容
        let cells = $(this).find('td'); // 获取当前行中的所有 td 单元格
        cells.each(function (index) {
          if (index < cells.length - 1) { // 遍历时排除最后一列（假设不需要导出）
            rowContent += $(this).text().trim().replace(/,/g, '') + ','; // 将单元格文本内容加入行内容，并用逗号分隔，替换掉内容中的逗号以避免干扰 CSV 格式
          }
        });
        csvContent += rowContent.slice(0, -1) + '\n'; // 移除最后一个多余的逗号并添加换行符
      });
      // 创建一个 Blob 对象，包含生成的 CSV 内容，类型为 text/csv，并设置字符编码为 UTF-8
      let blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
      // 创建一个指向 Blob 对象的 URL
      let url = URL.createObjectURL(blob);
      // 创建一个隐藏的 <a> 元素，并设置其 href 属性为 Blob URL，设置下载文件名为 articles.csv
      let link = $('<a>').attr({
        href: url,
        download: 'articles.csv'
      }).appendTo('body'); // 将链接元素临时添加到文档的 body 中
      // 模拟点击 <a> 元素以触发下载
      link[0].click();
      // 移除临时添加的 <a> 元素
      link.remove();
    });
  }
</script>
```

## 首頁设置(使用MongoDB)

设置数据保存到mongodb中

安装 mongodb 所用的库

```bash
npm install mongoose @nestjs/mongoose
```

添加环境变量

```
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=cms
MONGO_USER=root
MONGO_PASSWORD=root
```

### configuration.service

```ts
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ConfigurationService {
  constructor(private configService: ConfigService) { }
  get mysqlHost(): string {
    return this.configService.get<string>('MYSQL_HOST')!;
  }
  get mysqlPort(): number {
    return this.configService.get<number>('MYSQL_PORT')!;
  }
  get mysqlDb(): string {
    return this.configService.get<string>('MYSQL_DB')!;
  }
  get mysqlUser(): string {
    return this.configService.get<string>('MYSQL_USER')!;
  }
  get mysqlPass(): string {
    return this.configService.get<string>('MYSQL_PASSWORD')!;
  }
  get mysqlConfig() {
    return {
      host: this.mysqlHost,
      port: this.mysqlPort,
      database: this.mysqlDb,
      username: this.mysqlUser,
      password: this.mysqlPass,
    };
  }
  get mongodbHost(): string { // [!code ++]
    return this.configService.get<string>('MONGO_HOST')!; // [!code ++]
  } // [!code ++]
  get mongodbPort(): number { // [!code ++]
    return this.configService.get<number>('MONGO_PORT')!; // [!code ++]
  } // [!code ++]
  get mongodbDB(): string { // [!code ++]
    return this.configService.get<string>('MONGO_DB')!; // [!code ++]
  } // [!code ++]
  get mongodbUser(): string { // [!code ++]
    return this.configService.get<string>('MONGO_USER')!; // [!code ++]
  } // [!code ++]
  get mongodbPassword(): string { // [!code ++]
    return this.configService.get<string>('MONGO_PASSWORD')!; // [!code ++]
  } // [!code ++]
  get mongodbConfig() { // [!code ++]
    return { // [!code ++]
      uri: `mongodb://${this.mongodbHost}:${this.mongodbPort}/${this.mongodbDB}` // [!code ++]
    } // [!code ++]
  } // [!code ++]
}
```

### mongodb-base.service

```ts
import { Model } from 'mongoose';

export abstract class MongoDBBaseService<T, C, U> {
  constructor(
    protected readonly model: Model<T>,
  ) { }

  async findAll() {
    return await this.model.find();
  }

  async findOne(id: string) {
    return await this.model.findById(id);
  }

  async create(createDto: C) {
    const createdEntity = new this.model(createDto);
    await createdEntity.save();
    return createdEntity;
  }

  async update(id: string, updateDto: U) {
    await this.model.findByIdAndUpdate(id, updateDto as any, { new: true });
  }

  async delete(id: string) {
    await this.model.findByIdAndDelete(id);
  }
}
```

### setting.service

```ts
import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from 'mongoose';
import { SettingDocument } from '../schemas/setting.schema';
import { CreateSettingDto, UpdateSettingDto } from '../dtos/setting.dto';
import { MongoDBBaseService } from './mongodb-base.service';

@Injectable()
export class SettingService extends MongoDBBaseService<SettingDocument, CreateSettingDto, UpdateSettingDto> {
  constructor(@InjectModel('Setting') settingModel: Model<SettingDocument>) {
    super(settingModel);
  }
  async findFirst(): Promise<SettingDocument | null> {
    return await this.model.findOne().exec();
  }
}
```

### share.module

```ts
import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { RoleService } from './services/role.service';
import { AccessService } from "./services/access.service";
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';
import { Role } from './entities/role.entity';
import { Access } from "./entities/access.entity";
import { Article } from './entities/article.entity';
import { Category } from './entities/category.entity';
import { Tag } from './entities/tag.entity';
import { ArticleService } from './services/article.service';
import { CategoryService } from './services/category.service';
import { TagService } from './services/tag.service';
import { CosService } from './services/cos.service';
import { NotificationService } from './services/notification.service';
import { MailService } from './services/mail.service';
import { WordExportService } from './services/word-export.service';
import { PptExportService } from './services/ppt-export.service';
import { ExcelExportService } from './services/excel-export.service'
import { MongooseModule } from '@nestjs/mongoose'; // [!code ++]
import { Setting, SettingSchema } from './schemas/setting.schema'; // [!code ++]
import { SettingService } from './services/setting.service'; // [!code ++]

@Global()
@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] }),
        MongooseModule.forRootAsync({ // [!code ++]
            inject: [ConfigurationService], // [!code ++]
            useFactory: (configurationService: ConfigurationService) => ({ // [!code ++]
                uri: configurationService.mongodbConfig.uri // [!code ++]
            }), // [!code ++]
        }), // [!code ++]
        MongooseModule.forFeature([ // [!code ++]
            { name: Setting.name, schema: SettingSchema }, // [!code ++]
        ]), // [!code ++]
        TypeOrmModule.forFeature([User, Role, Access, Article, Category, Tag]),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            inject: [ConfigurationService],
            useFactory: (configService: ConfigurationService) => ({
                type: 'mysql',
                ...configService.mysqlConfig,
                entities: [User, Role, Access, Article, Category, Tag],
                synchronize: true,
                autoLoadEntities: true,
                logging: false
            }),
        }),
    ],
    providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService, WordExportService, PptExportService, ExcelExportService, SettingService], // [!code ++]
    exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService, WordExportService, PptExportService, ExcelExportService, SettingService], // [!code ++]
})
export class ShareModule {
}
```



### setting.dto

```ts
import { ApiProperty } from '@nestjs/swagger';
import { PartialType } from '@nestjs/mapped-types';

export class CreateSettingDto {
  @ApiProperty({ description: '网站名称', example: '我的网站' })
  siteName: string;

  @ApiProperty({ description: '网站描述', example: '这是我的个人网站' })
  siteDescription: string;

  @ApiProperty({ description: '联系邮箱', example: 'contact@example.com' })
  contactEmail: string;
}

export class UpdateSettingDto extends PartialType(CreateSettingDto) {
  id: string
}
```

### setting.schema

```ts
import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { HydratedDocument } from 'mongoose';
export type SettingDocument = HydratedDocument<Setting>;
@Schema()
export class Setting {
  id: string;
  @Prop({ required: true })
  siteName: string;
  @Prop()
  siteDescription: string;
  @Prop()
  contactEmail: string;
}
export const SettingSchema = SchemaFactory.createForClass(Setting);
SettingSchema.virtual('id').get(function () {
  return this._id.toHexString();
});
SettingSchema.set('toJSON', { virtuals: true });
SettingSchema.set('toObject', { virtuals: true });
```

### setting.controller

```ts
import { Controller, Get, Post, Body, Render, Redirect } from '@nestjs/common';
import { SettingService } from '../../share/services/setting.service';
import { UpdateSettingDto } from '../../share/dtos/setting.dto';

@Controller('admin/settings')
export class SettingController {
  constructor(private readonly settingService: SettingService) { }

  @Get()
  @Render('settings')
  async getSettings() {
    let settings = await this.settingService.findFirst();
    if (!settings) {
      settings = await this.settingService.create({
        siteName: '默认网站名称',
        siteDescription: '默认网站描述',
        contactEmail: 'default@example.com',
      });
    }
    return { settings };
  }

  @Post()
  @Redirect('/admin/dashboard')
  async updateSettings(@Body() updateSettingDto: UpdateSettingDto) {
    await this.settingService.update(updateSettingDto.id, updateSettingDto);
    return { success: true };
  }
}
```

### admin.module

```ts
import { Module } from '@nestjs/common';
import { DashboardController } from './controllers/dashboard.controller';
import { UserController } from './controllers/user.controller';
import { AdminExceptionFilter } from './filters/admin-exception.filter';
import { RoleController } from "./controllers/role.controller";
import { AccessController } from "./controllers/access.controller";
import { ArticleController } from './controllers/article.controller';
import { CategoryController } from './controllers/category.controller';
import { TagController } from './controllers/tag.controller';
import { UploadController } from './controllers/upload.controller';
import { SettingController } from './controllers/setting.controller'; // [!code ++]

@Module({
  controllers: [
    DashboardController,
    UserController,
    RoleController,
    AccessController,
    ArticleController,
    CategoryController,
    TagController,
    UploadController,
    SettingController // [!code ++]
  ],
  providers: [{
    provide: 'APP_FILTER',
    useClass: AdminExceptionFilter,
  }],
})
export class AdminModule { }
```

### sidebar.hbs

```handlebars
<div class="col-md-3 col-lg-2 p-0">
  <div class="accordion" id="sidebarMenu">
    <div class="accordion-item">
      <h2 class="accordion-header" id="heading">
        <button class="accordion-button" type="button" data-bs-toggle="collapse"
          data-bs-target="#collapse1">权限管理</button>
      </h2>
      <div class="accordion-collapse collapse" id="collapse1">
        <div class="accordion-body">
          <ul class="list-group">
            <li class="list-group-item">
              <a href="/admin/users">用户管理</a>
            </li>
            <li class="list-group-item">
              <a href="/admin/roles">角色管理</a>
            </li>
            <li class="list-group-item">
              <a href="/admin/accesses">资源管理</a>
            </li>
          </ul>
        </div>
      </div>
      <h2 class="accordion-header" id="heading">
        <button class="accordion-button" type="button" data-bs-toggle="collapse"
          data-bs-target="#collapse2">内容管理</button>
      </h2>
      <div class="accordion-collapse collapse" id="collapse2">
        <div class="accordion-body">
          <ul class="list-group">
            <li class="list-group-item">
              <a href="/admin/tags">标签管理</a>
            </li>
            <li class="list-group-item">
              <a href="/admin/categories">分类管理</a>
            </li>
            <li class="list-group-item">
              <a href="/admin/articles">文章管理</a>
            </li>
          </ul>
        </div>
      </div>
      <h2 class="accordion-header" id="heading">
        <button class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#collapse3">设置</button>
      </h2>
      <div class="accordion-collapse collapse" id="collapse3">
        <div class="accordion-body">
          <ul class="list-group">
            <li class="list-group-item">
              <a href="/admin/settings">网站设置</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
```

### settings.hbs

```handlebars
<form action="/admin/settings" method="post">
  <input type="hidden" name="id" value="{{settings.id}}">
  <div class="mb-3">
    <label for="siteName" class="form-label">网站名称</label>
    <input type="text" class="form-control" id="siteName" name="siteName" value="{{settings.siteName}}" required>
  </div>
  <div class="mb-3">
    <label for="siteDescription" class="form-label">网站描述</label>
    <input type="text" class="form-control" id="siteDescription" name="siteDescription"
      value="{{settings.siteDescription}}" required>
  </div>
  <div class="mb-3">
    <label for="contactEmail" class="form-label">联系邮箱</label>
    <input type="email" class="form-control" id="contactEmail" name="contactEmail" value="{{settings.contactEmail}}"
      required>
  </div>
  <button type="submit" class="btn btn-primary">保存设置</button>
</form>
```

