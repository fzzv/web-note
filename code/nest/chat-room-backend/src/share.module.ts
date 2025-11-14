import { Global, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RedisService } from './service/redis.service';
import { EmailService } from './service/email.service';
import { UserService } from './service/user.service';
import { ConfigurationService } from './service/configuration.service';
import { UtilityService } from './service/utility.service';
import { FriendshipService } from './service/friendship.service';
import { ChatroomService } from './service/chatroom.service';
import * as Minio from 'minio';

@Global()
@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] })],
  providers: [
    {
      provide: 'MINIO_CLIENT',
      async useFactory() {
        const client = new Minio.Client({
          endPoint: 'localhost',
          port: 9000,
          useSSL: false,
          accessKey: 'NvAsQeTYG2zoLCrtbWJm',
          secretKey: 'dLoYZA0BePBvsP26kP05pMo7Ra3TRs1PswwTIo5p'
        })
        return client;
      }
    },
    RedisService, EmailService, UserService, ConfigurationService, UtilityService, FriendshipService, ChatroomService],
  exports: ['MINIO_CLIENT', RedisService, EmailService, UserService, ConfigurationService, UtilityService, FriendshipService, ChatroomService],
})
export class ShareModule {
}
