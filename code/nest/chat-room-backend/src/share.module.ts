import { Global, Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RedisService } from './service/redis.service';
import { EmailService } from './service/email.service';
import { UserService } from './service/user.service';
import { ConfigurationService } from './service/configuration.service';
import { UtilityService } from './service/utility.service';
import { FriendshipService } from './service/friendship.service';
import { ChatroomService } from './service/chatroom.service';

@Global()
@Module({
  imports: [ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] })],
  providers: [RedisService, EmailService, UserService, ConfigurationService, UtilityService, FriendshipService, ChatroomService],
  exports: [RedisService, EmailService, UserService, ConfigurationService, UtilityService, FriendshipService, ChatroomService],
})
export class ShareModule {
}
