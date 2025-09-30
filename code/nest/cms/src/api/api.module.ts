import { Module } from '@nestjs/common';
import { UserController } from './controllers/user.controller';
import { AuthController } from './controllers/auth.controller';
import { CategoryController } from './controllers/category.controller';
import { ArticleController } from './controllers/article.controller';
import { JwtModule } from '@nestjs/jwt';

@Module({
  imports: [
    JwtModule.register({
      global: true,
      signOptions: { expiresIn: '7d' }
    }),
  ],
  controllers: [UserController, AuthController, CategoryController, ArticleController]
})
export class ApiModule { }
