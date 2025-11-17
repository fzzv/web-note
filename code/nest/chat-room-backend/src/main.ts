import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';
import { TransformInterceptor } from './interceptors/transform.interceptor';
import { AllExceptionsFilter } from './filters/http-exception.filter';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  // 配置全局管道 transform: true 自动转换请求参数为DTO对象
  app.useGlobalPipes(new ValidationPipe({ transform: true }));
  // 配置 swagger 文档
  const config = new DocumentBuilder()
    .setTitle('Chat Room API')
    .setDescription('Chat Room API Description')
    .setVersion('1.0')
    .addTag('Chat Room')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api-doc', app, document);
  // 1. 全局注册成功响应拦截器
  app.useGlobalInterceptors(new TransformInterceptor());
  // 2. 全局注册异常过滤器
  // 注意：HttpAdapter 必须在 useGlobalFilters 之前获取
  // const { httpAdapter } = app.get(HttpAdapterHost); // 如果 AllExceptionsFilter 需要
  app.useGlobalFilters(new AllExceptionsFilter());
  app.enableCors();
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
