import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AdminModule } from './admin/admin.module';
import { ApiModule } from './api/api.module';
import { ShareModule } from './share/share.module';

@Module({
  imports: [AdminModule, ApiModule, ShareModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
