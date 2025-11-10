import { Global, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RedisService } from './service/redis.service';
import { EmailService } from './service/email.service';
import { UserService } from './service/user.service';
import { ConfigurationService } from './service/configuration.service';

@Global()
@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] })],
  providers: [RedisService, EmailService, UserService, ConfigurationService],
  exports: [RedisService, EmailService, UserService, ConfigurationService],
})
export class ShareModule {
}
