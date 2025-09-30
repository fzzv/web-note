import { ApiProperty } from '@nestjs/swagger';
import { PartialType } from '@nestjs/mapped-types';

export class CreateSettingDto {
  @ApiProperty({ description: '网站名称', example: '我的网站' })
  siteName: string;

  @ApiProperty({ description: '网站描述', example: '这是我的个人网站' })
  siteDescription: string;

  @ApiProperty({ description: '联系邮箱', example: 'contact@example.com' })
  contactEmail: string;
}

export class UpdateSettingDto extends PartialType(CreateSettingDto) {
  id: string
}
