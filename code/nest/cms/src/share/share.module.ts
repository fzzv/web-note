import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';

@Global()
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    TypeOrmModule.forFeature([User]),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigurationService],
      useFactory: (configService: ConfigurationService) => ({
        type: 'mysql',
        ...configService.mysqlConfig,
        entities: [User],
        synchronize: true,
        autoLoadEntities: true,
        logging: false
      }),
    }),
  ],
  providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint],
  exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint],
})
export class ShareModule {}

