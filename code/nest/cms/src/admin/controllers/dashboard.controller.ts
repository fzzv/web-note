import { Controller, Get, Render } from '@nestjs/common';
import { ApiCookieAuth, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { DashboardService } from '../../share/services/dashboard.service';

@ApiTags('admin/dashboard')
@Controller('admin')
export class DashboardController {

  constructor(private readonly dashboardService: DashboardService) { }

  @Get('dashboard')
  @ApiCookieAuth()
  @ApiOperation({ summary: '管理后台仪表盘' })
  @ApiResponse({ status: 200, description: '成功返回仪表盘页面' })
  @Render('dashboard')
  async dashboard() {
    return await this.dashboardService.getDashboardData();  
  }
}
