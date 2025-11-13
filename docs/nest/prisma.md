# Prisma

## 在nest.js中使用Prisma，统一返回接口格式

在 Nest.js 中，实现这一目标的最佳组合是 拦截器 (Interceptor) 和 异常过滤器 (Exception Filter)。

- 拦截器 (Interceptor)：用于处理成功的响应，将其包装成你的标准格式。

- 异常过滤器 (Exception Filter)：用于捕获所有抛出的异常（包括 Nest.js 的 HttpException、Prisma 的 PrismaClientKnownRequestError 等），并将其转换为你的标准错误格式。

### 1. 定义标准响应结构

首先，定义一个你的API响应将遵循的结构。通常会创建一个 DTO 或 Interface：

```ts
// src/interfaces/api-response.interface.ts
export interface ApiResponse<T = any> {
  success: boolean;
  code: number; // 通常是 HTTP 状态码或自定义业务码
  data: T | null;
  message: string;
}
```

### 2. 统一成功响应 (使用拦截器)

当你的 Controller/Service 成功处理请求并返回数据时（例如，`return this.prisma.user.findMany()`），拦截器会捕获这个返回值，并将其包装。

```ts
// src/interceptors/transform.interceptor.ts
import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiResponse } from '../interfaces/api-response.interface';

@Injectable()
export class TransformInterceptor<T> implements NestInterceptor<T, ApiResponse<T>> {
  intercept(context: ExecutionContext, next: CallHandler): Observable<ApiResponse<T>> {
    
    // 获取 HTTP 状态码
    const httpStatusCode = context.switchToHttp().getResponse().statusCode;

    return next.handle().pipe(
      map(data => ({
        success: true,
        code: httpStatusCode, // 使用 HTTP 状态码
        data: data, // data 是你 controller/service 返回的原始数据
        message: 'Request successful',
      })),
    );
  }
}
```

### 3. 统一异常响应 (使用异常过滤器)

这是处理 Prisma 错误的关键。你需要一个过滤器来捕获所有异常。

```ts
// src/filters/http-exception.filter.ts
import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus } from '@nestjs/common';
import { Request, Response } from 'express';
import { Prisma } from '@prisma/client';

@Catch() // 捕获所有类型的异常
export class AllExceptionsFilter implements ExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();

    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = 'Internal server error';
    let code = 500;

    // 1. 处理 Nest.js 的 HttpException
    if (exception instanceof HttpException) {
      status = exception.getStatus();
      message = exception.message;
      code = status;
    } 
    // 2. 处理 Prisma 特定的已知请求错误
    else if (exception instanceof Prisma.PrismaClientKnownRequestError) {
      message = `[Prisma Error ${exception.code}] ${exception.message.split('\n').pop()}`; // 简化 Prisma 错误信息

      switch (exception.code) {
        case 'P2002': // 唯一约束失败 (e.g., email already exists)
          status = HttpStatus.CONFLICT; // 409
          message = 'A record with this value already exists.';
          break;
        case 'P2025': // 记录未找到 (e.g., for update/delete)
          status = HttpStatus.NOT_FOUND; // 404
          message = 'Record to update/delete not found.';
          break;
        // ... 在这里添加更多你关心的 Prisma 错误码
        default:
          status = HttpStatus.BAD_REQUEST; // 400
          break;
      }
      code = status; // 也可以使用 Prisma code 作为业务码
    } 
    // 3. 处理其他未知错误
    else if (exception instanceof Error) {
      message = exception.message;
    }
    
    // ... 可以在这里添加日志记录 (Logger) ...
    // logger.error(message, exception.stack, `${request.method} ${request.url}`);

    response.status(status).json({
      success: false,
      code: code,
      data: null,
      message: message,
      timestamp: new Date().toISOString(), // 额外信息
      path: request.url,                // 额外信息
    });
  }
}
```

### 4. 全局注册

最后，在你的 `main.ts` 中全局应用它们

```ts
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { TransformInterceptor } from './interceptors/transform.interceptor'; // [!code ++]
import { AllExceptionsFilter } from './filters/http-exception.filter'; // [!code ++]

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  // 配置全局管道 transform: true 自动转换请求参数为DTO对象
  app.useGlobalPipes(new ValidationPipe({ transform: true }));
  // 1. 全局注册成功响应拦截器
  app.useGlobalInterceptors(new TransformInterceptor()); // [!code ++]
  // 2. 全局注册异常过滤器
  // 注意：HttpAdapter 必须在 useGlobalFilters 之前获取
  // const { httpAdapter } = app.get(HttpAdapterHost); // 如果 AllExceptionsFilter 需要
  app.useGlobalFilters(new AllExceptionsFilter()); // [!code ++]
  app.enableCors();
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
```

### 5.controller/service中如何使用

如何返回自己所需的code，比如在登录接口，当用户名或密码错误时，需要返回状态码为500，可以如下写：

```ts
@Post('login')
async login(@Body() body: LoginUserDto) {
    const { username, password } = body;
    const user = await this.validateUser(username, password);
    if (!user) {
        throw new HttpException('用户名或密码错误', HttpStatus.INTERNAL_SERVER_ERROR) // 500
    }
    return {
        user: username,
        token: this.createJwtTokens(user)
    };
}
```



