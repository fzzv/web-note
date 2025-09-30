import { Controller, Get, Post, Body, Render, Redirect } from '@nestjs/common';
import { SettingService } from '../../share/services/setting.service';
import { UpdateSettingDto } from '../../share/dtos/setting.dto';

@Controller('admin/settings')
export class SettingController {
  constructor(private readonly settingService: SettingService) { }

  @Get()
  @Render('settings')
  async getSettings() {
    let settings = await this.settingService.findFirst();
    if (!settings) {
      settings = await this.settingService.create({
        siteName: '默认网站名称',
        siteDescription: '默认网站描述',
        contactEmail: 'default@example.com',
      });
    }
    return { settings };
  }

  @Post()
  @Redirect('/admin/dashboard')
  async updateSettings(@Body() updateSettingDto: UpdateSettingDto) {
    await this.settingService.update(updateSettingDto.id, updateSettingDto);
    return { success: true };
  }
}
