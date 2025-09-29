import { MiddlewareConsumer, Module, NestModule } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module';
import { ApiModule } from './api/api.module';
import { ShareModule } from './share/share.module';
import methodOverride from './share/middlewares/method-override';
import { ServeStaticModule } from '@nestjs/serve-static';
import * as path from 'path';
import { EventEmitterModule } from '@nestjs/event-emitter';

@Module({
  imports: [
    // 配置 EventEmitterModule 模块
    EventEmitterModule.forRoot({
      // 启用通配符功能，允许使用通配符来订阅事件
      wildcard: true,
      // 设置事件名的分隔符，这里使用 '.' 作为分隔符
      delimiter: '.',
      // 将事件发射器设置为全局模块，所有模块都可以共享同一个事件发射器实例
      global: true
    }),
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
