import { Module } from "@nestjs/common";
import { JwtModule } from "@nestjs/jwt";
import { UserController } from "src/controllers/user.controller";
import { EmailController } from "src/controllers/email.controller";
import { FriendshipController } from "src/controllers/friendship.controller";

@Module({
  imports: [
    JwtModule.register({
      global: true,
      signOptions: { expiresIn: '7d' }
    }),
  ],
  controllers: [
    UserController,
    EmailController,
    FriendshipController,
  ],
})
export class ApiModule { }
