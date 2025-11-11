import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ValidationPipe } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  // 配置全局管道 transform: true 自动转换请求参数为DTO对象
  app.useGlobalPipes(new ValidationPipe({ transform: true }));
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
