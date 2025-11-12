import { BadRequestException, Body, Controller, Get, Param, Post } from '@nestjs/common';
import { FriendshipService } from 'src/service/friendship.service';
import { FriendAddDto } from 'src/dtos/friend-add.dto';
import { UserInfo } from 'src/decorators/index';

@Controller('friendship')
export class FriendshipController {
  constructor(private readonly friendshipService: FriendshipService) { }

  @Get('list')
  async friendship(@UserInfo("id") userId: number) {
    return this.friendshipService.getFriendship(userId);
  }

  @Post('add')
  async add(@Body() friendAddDto: FriendAddDto, @UserInfo("id") userId: number) {
    return this.friendshipService.add(friendAddDto, userId);
  }

  @Get('request_list')
  async list(@UserInfo("id") userId: number) {
    return this.friendshipService.list(userId);
  }

  @Get('agree/:id')
  async agree(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('添加的好友 id 不能为空');
    }
    return this.friendshipService.agree(friendId, userId);
  }

  @Get('reject/:id')
  async reject(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('添加的好友 id 不能为空');
    }
    return this.friendshipService.reject(friendId, userId);
  }

  @Get('remove/:id')
  async remove(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('删除的好友 id 不能为空');
    }
    return this.friendshipService.remove(friendId, userId);
  }
}
