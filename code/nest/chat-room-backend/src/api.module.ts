import { Module } from "@nestjs/common";
import { UserController } from "src/controllers/user.controller";
import { EmailController } from "src/controllers/email.controller";
import { JwtModule } from "@nestjs/jwt";

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
  ],
})
export class ApiModule { }
