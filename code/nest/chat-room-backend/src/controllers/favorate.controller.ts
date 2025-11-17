import { Controller, Get, Query } from '@nestjs/common';
import { FavoriteService } from 'src/service/favorite.service';
import { UserInfo } from 'src/decorators';

@Controller('favorite')
export class FavoriteController {
  constructor(private readonly favoriteService: FavoriteService) { }

  @Get('list')
  async list(@UserInfo('id') userId: number) {
    return this.favoriteService.list(userId);
  }

  @Get('add')
  async add(@UserInfo('id') userId: number, @Query('chatHistoryId') chatHistoryId: number) {
    return this.favoriteService.add(userId, chatHistoryId);
  }

  @Get('del')
  async del(@Query('id') id: number) {
    return this.favoriteService.del(id);
  }
}
