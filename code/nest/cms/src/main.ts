import { NestFactory } from '@nestjs/core';
import session from 'express-session';
import cookieParser from 'cookie-parser';
import { join } from 'node:path';
import { engine } from 'express-handlebars';
import { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';

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
  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
