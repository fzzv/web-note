import { Controller, Get, Render } from '@nestjs/common';

@Controller('admin')
export class DashboardController {
  @Get()
  @Render('dashboard')
  dashboard() {
    return { title: 'dashboard' }
  }
}
