import { Module } from '@nestjs/common';
import { DashboardController } from './controllers/dashboard.controller';
import { UserController } from './controllers/user.controller';
import { AdminExceptionFilter } from './filters/admin-exception.filter';
import { RoleController } from "./controllers/role.controller";
@Module({
  controllers: [DashboardController, UserController, RoleController],
  providers: [{
    provide: 'APP_FILTER',
    useClass: AdminExceptionFilter,
  }],
})
export class AdminModule {}
