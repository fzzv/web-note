import { Controller, Get, Render } from '@nestjs/common';
import { ApiCookieAuth, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('admin/dashboard')
@Controller('admin')
export class DashboardController {
  @Get('dashboard')
  @ApiCookieAuth()
  @ApiOperation({ summary: '管理后台仪表盘' })
  @ApiResponse({ status: 200, description: '成功返回仪表盘页面' })
  @Render('dashboard')
  dashboard() {
    return { title: 'dashboard' }
  }
}
