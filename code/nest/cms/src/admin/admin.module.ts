import { Module } from '@nestjs/common';
import { DashboardController } from './controllers/dashboard.controller';
import { UserController } from './controllers/user.controller';
import { AdminExceptionFilter } from './filters/admin-exception.filter';
import { RoleController } from "./controllers/role.controller";
import { AccessController } from "./controllers/access.controller";

@Module({
  controllers: [DashboardController, UserController, RoleController, AccessController],
  providers: [{
    provide: 'APP_FILTER',
    useClass: AdminExceptionFilter,
  }],
})
export class AdminModule {}
