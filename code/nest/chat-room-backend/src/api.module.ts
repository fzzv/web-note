import { Module } from "@nestjs/common";
import { UserController } from "src/controllers/user.controller";
import { EmailController } from "src/controllers/email.controller";

@Module({
  controllers: [
    UserController,
    EmailController,
  ],
})
export class ApiModule {}
