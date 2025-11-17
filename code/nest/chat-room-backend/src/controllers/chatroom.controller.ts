import { BadRequestException, Controller, Get, Param, Query } from '@nestjs/common';
import { ChatroomService } from 'src/service/chatroom.service';
import { UserInfo } from 'src/decorators/index';
import { ApiOperation, ApiParam, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('chatroom')
@Controller('chatroom')
export class ChatroomController {
  constructor(private readonly chatroomService: ChatroomService) { }

  @Get('create-one-to-one')
  @ApiOperation({ summary: '创建一对一聊天室' })
  @ApiQuery({ name: 'friendId', description: '聊天好友的 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回一对一聊天室' })
  async oneToOne(@Query('friendId') friendId: number, @UserInfo('id') userId: number) {
    if (!friendId) {
      throw new BadRequestException('聊天好友的 id 不能为空');
    }
    return this.chatroomService.createOneToOneChatroom(friendId, userId);
  }

  @Get('create-group')
  @ApiOperation({ summary: '创建群聊聊天室' })
  @ApiQuery({ name: 'name', description: '群聊名称', required: true })
  @ApiResponse({ status: 200, description: '成功返回群聊聊天室' })
  async group(@Query('name') name: string, @UserInfo('id') userId: number) {
    return this.chatroomService.createGroupChatroom(name, userId);
  }

  @Get('list')
  @ApiOperation({ summary: '获取聊天室列表' })
  @ApiQuery({ name: 'name', description: '聊天室名称', required: false })
  @ApiResponse({ status: 200, description: '成功返回聊天室列表' })
  async list(@UserInfo('id') userId: number, @Query('name') name: string) {
    if (!userId) {
      throw new BadRequestException('userId 不能为空')
    }
    return this.chatroomService.list(userId, name);
  }

  @Get('members')
  @ApiOperation({ summary: '获取聊天室成员列表' })
  @ApiQuery({ name: 'chatroomId', description: '聊天室 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回聊天室成员列表' })
  async members(@Query('chatroomId') chatroomId: number) {
    if (!chatroomId) {
      throw new BadRequestException('chatroomId 不能为空')
    }
    return this.chatroomService.members(chatroomId);
  }

  @Get('info/:id')
  @ApiOperation({ summary: '获取聊天室信息' })
  @ApiParam({ name: 'id', description: '聊天室 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回聊天室信息' })
  async info(@Param('id') id: number) {
    if (!id) {
      throw new BadRequestException('id 不能为空')
    }
    return this.chatroomService.info(id);
  }

  @Get('join/:id')
  @ApiOperation({ summary: '加入聊天室' })
  @ApiParam({ name: 'id', description: '聊天室 id', required: true })
  @ApiQuery({ name: 'joinUsername', description: '加入用户名', required: true })
  @ApiResponse({ status: 200, description: '成功返回加入聊天室' })
  async join(@Param('id') id: number, @Query('joinUsername') joinUsername: string) {
    if (!id) {
      throw new BadRequestException('id 不能为空')
    }
    if (!joinUsername) {
      throw new BadRequestException('joinUsername 不能为空')
    }
    return this.chatroomService.join(id, joinUsername);
  }

  @Get('quit/:id')
  @ApiOperation({ summary: '退出聊天室' })
  @ApiParam({ name: 'id', description: '聊天室 id', required: true })
  @ApiQuery({ name: 'quitUserId', description: '退出用户 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回退出聊天室' })
  async quit(@Param('id') id: number, @Query('quitUserId') quitUserId: number) {
    if (!id) {
      throw new BadRequestException('id 不能为空')
    }
    if (!quitUserId) {
      throw new BadRequestException('quitUserId 不能为空')
    }
    return this.chatroomService.quit(id, quitUserId);
  }

  @Get('findChatroom')
  @ApiOperation({ summary: '查找一对一聊天室' })
  @ApiQuery({ name: 'userId1', description: '用户 id 1', required: true })
  @ApiQuery({ name: 'userId2', description: '用户 id 2', required: true })
  @ApiResponse({ status: 200, description: '成功返回一对一聊天室' })
  async findChatroom(@Query('userId1') userId1: string, @Query('userId2') userId2: string) {
    if (!userId1 || !userId2) {
      throw new BadRequestException('用户 id 不能为空');
    }
    return this.chatroomService.queryOneToOneChatroom(+userId1, +userId2);
  }
}
