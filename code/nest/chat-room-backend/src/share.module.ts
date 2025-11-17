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
import { ChatService } from './service/chat.service';
import { ChatGateway } from './gateway/chat.gateway';
import { ChatHistoryService } from './service/chat-history.service';
import { FavoriteService } from './service/favorite.service';

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
          accessKey: 'lxy5LBZxZzRNxJouncsA',
          secretKey: 'EuSpsozFDIuZnIfaOLkiVAMEe4vGzeBFbfjgWn5A'
        })
        return client;
      }
    },
    ChatGateway,
    RedisService, EmailService, UserService, ConfigurationService, UtilityService, FriendshipService,
    ChatroomService, ChatService, ChatHistoryService, FavoriteService],
  exports: [
    'MINIO_CLIENT', RedisService, EmailService, UserService, ConfigurationService, UtilityService,
    FriendshipService, ChatroomService, ChatService, ChatHistoryService, FavoriteService],
})
export class ShareModule {
}
