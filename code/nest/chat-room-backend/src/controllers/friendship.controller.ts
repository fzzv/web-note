import { BadRequestException, Body, Controller, Get, Param, Post, Query } from '@nestjs/common';
import { FriendshipService } from 'src/service/friendship.service';
import { FriendAddDto } from 'src/dtos/friend-add.dto';
import { UserInfo } from 'src/decorators/index';
import { ApiBody, ApiOperation, ApiParam, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('friendship')
@Controller('friendship')
export class FriendshipController {
  constructor(private readonly friendshipService: FriendshipService) { }

  @Get('list')
  @ApiOperation({ summary: '获取好友列表' })
  @ApiQuery({ name: 'name', description: '好友名称', required: true })
  @ApiResponse({ status: 200, description: '成功返回好友列表' })
  async friendship(@UserInfo("id") userId: number, @Query('name') name: string) {
    return this.friendshipService.getFriendship(userId, name);
  }

  @Post('add')
  @ApiOperation({ summary: '添加好友' })
  @ApiBody({ type: FriendAddDto })
  @ApiResponse({ status: 200, description: '成功返回添加好友' })
  async add(@Body() friendAddDto: FriendAddDto, @UserInfo("id") userId: number) {
    return this.friendshipService.add(friendAddDto, userId);
  }

  @Get('request_list')
  @ApiOperation({ summary: '获取好友请求列表' })
  @ApiResponse({ status: 200, description: '成功返回好友请求列表' })
  async list(@UserInfo("id") userId: number) {
    return this.friendshipService.list(userId);
  }

  @Get('agree/:id')
  @ApiOperation({ summary: '同意好友请求' })
  @ApiParam({ name: 'id', description: '好友请求 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回同意好友请求' })
  async agree(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('添加的好友 id 不能为空');
    }
    return this.friendshipService.agree(friendId, userId);
  }

  @Get('reject/:id')
  @ApiOperation({ summary: '拒绝好友请求' })
  @ApiParam({ name: 'id', description: '好友请求 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回拒绝好友请求' })
  async reject(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('添加的好友 id 不能为空');
    }
    return this.friendshipService.reject(friendId, userId);
  }

  @Get('remove/:id')
  @ApiOperation({ summary: '删除好友' })
  @ApiParam({ name: 'id', description: '好友 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回删除好友' })
  async remove(@Param('id') friendId: number, @UserInfo("id") userId: number) {
    if (!friendId) {
      throw new BadRequestException('删除的好友 id 不能为空');
    }
    return this.friendshipService.remove(friendId, userId);
  }
}
