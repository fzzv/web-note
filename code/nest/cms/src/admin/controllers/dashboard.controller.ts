import { Controller, Get, Render, Sse } from '@nestjs/common';
import { ApiCookieAuth, ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { DashboardService } from '../../share/services/dashboard.service';
import { WeatherService } from '../../share/services/weather.service';
import { interval, map, mergeMap } from 'rxjs';
import { SystemService } from '../../share/services/system.service';

@ApiTags('admin/dashboard')
@Controller('admin')
export class DashboardController {

  constructor(private readonly dashboardService: DashboardService, private readonly weatherService: WeatherService, private readonly systemService: SystemService) { }

  @Get('dashboard')
  @ApiCookieAuth()
  @ApiOperation({ summary: '管理后台仪表盘' })
  @ApiResponse({ status: 200, description: '成功返回仪表盘页面' })
  @Render('dashboard')
  async dashboard() {
    return await this.dashboardService.getDashboardData();
  }

  @Get('dashboard/weather')
  async getWeather() {
    const weather = await this.weatherService.getWeather();
    return { weather };
  }

  @Sse('dashboard/systemInfo')
  systemInfo() {
    return interval(3000).pipe(
      mergeMap(() => this.systemService.getSystemInfo()),
      map((systemInfo) => ({ data: systemInfo }))
    );
  }
}
