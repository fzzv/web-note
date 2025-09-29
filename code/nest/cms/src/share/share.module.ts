import { Global, Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from './entities/user.entity';
import { ConfigModule } from '@nestjs/config';
import { ConfigurationService } from './services/configuration.service';
import { UserService } from './services/user.service';
import { RoleService } from './services/role.service';
import { AccessService } from "./services/access.service";
import { UtilityService } from './services/utility.service';
import { IsUsernameUniqueConstraint } from './validators/user-validators';
import { Role } from './entities/role.entity';
import { Access } from "./entities/access.entity";
import { Article } from './entities/article.entity';
import { Category } from './entities/category.entity';
import { Tag } from './entities/tag.entity';
import { ArticleService } from './services/article.service';
import { CategoryService } from './services/category.service';
import { TagService } from './services/tag.service';
import { CosService } from './services/cos.service';
import { NotificationService } from './services/notification.service';
import { MailService } from './services/mail.service';
import { WordExportService } from './services/word-export.service';
import { PptExportService } from './services/ppt-export.service';

@Global()
@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true, envFilePath: ['.env.local', '.env'] }),
        TypeOrmModule.forFeature([User, Role, Access, Article, Category, Tag]),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            inject: [ConfigurationService],
            useFactory: (configService: ConfigurationService) => ({
                type: 'mysql',
                ...configService.mysqlConfig,
                entities: [User, Role, Access, Article, Category, Tag],
                synchronize: true,
                autoLoadEntities: true,
                logging: false
            }),
        }),
    ],
    providers: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService, WordExportService, PptExportService],
    exports: [ConfigurationService, UserService, UtilityService, IsUsernameUniqueConstraint, RoleService, AccessService, ArticleService, CategoryService, TagService, CosService, NotificationService, MailService, WordExportService, PptExportService],
})
export class ShareModule {
}
