import { Controller, Get, Query } from '@nestjs/common';
import { FavoriteService } from 'src/service/favorite.service';
import { UserInfo } from 'src/decorators';
import { ApiOperation, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('favorite')
@Controller('favorite')
export class FavoriteController {
  constructor(private readonly favoriteService: FavoriteService) { }

  @Get('list')
  @ApiOperation({ summary: '获取收藏列表' })
  @ApiResponse({ status: 200, description: '成功返回收藏列表' })
  async list(@UserInfo('id') userId: number) {
    return this.favoriteService.list(userId);
  }

  @Get('add')
  @ApiOperation({ summary: '添加收藏' })
  @ApiQuery({ name: 'chatHistoryId', description: '聊天记录 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回添加收藏' })
  async add(@UserInfo('id') userId: number, @Query('chatHistoryId') chatHistoryId: number) {
    return this.favoriteService.add(userId, chatHistoryId);
  }

  @Get('del')
  @ApiOperation({ summary: '删除收藏' })
  @ApiQuery({ name: 'id', description: '收藏 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回删除收藏' })
  async del(@Query('id') id: number) {
    return this.favoriteService.del(id);
  }
}
