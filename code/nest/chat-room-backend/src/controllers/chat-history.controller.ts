import { Controller, Get, Query } from '@nestjs/common';
import { ChatHistoryService } from 'src/service/chat-history.service';
import { ApiOperation, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('chat-history')
@Controller('chat-history')
export class ChatHistoryController {
  constructor(private readonly chatHistoryService: ChatHistoryService) { }

  @Get('list')
  @ApiOperation({ summary: '获取聊天记录列表' })
  @ApiQuery({ name: 'chatroomId', description: '聊天室 id', required: true })
  @ApiResponse({ status: 200, description: '成功返回聊天记录列表' })
  async list(@Query('chatroomId') chatroomId: string) {
    return this.chatHistoryService.list(+chatroomId);
  }
}
