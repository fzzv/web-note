import { MessageBody, SubscribeMessage, WebSocketGateway, WebSocketServer } from '@nestjs/websockets';
import { ChatService } from 'src/service/chat.service';
import { Server, Socket } from 'socket.io';
import { ChatHistoryService } from 'src/service/chat-history.service';
import { UserService } from 'src/service/user.service';

interface JoinRoomPayload {
  chatroomId: number
  userId: number
}

interface SendMessagePayload {
  sendUserId: number;
  chatroomId: number;
  message: {
    type: 'text' | 'image' | 'file',
    content: string
  }
}

@WebSocketGateway({ cors: { origin: '*' } })
export class ChatGateway {
  constructor(
    private readonly chatService: ChatService,
    private readonly chatHistoryService: ChatHistoryService,
    private readonly userService: UserService) { }

  @WebSocketServer() server: Server;

  @SubscribeMessage('joinRoom')
  joinRoom(client: Socket, payload: JoinRoomPayload): void {
    const roomName = payload.chatroomId.toString();

    client.join(roomName)

    this.server.to(roomName).emit('message', {
      type: 'joinRoom',
      userId: payload.userId
    });
  }

  @SubscribeMessage('sendMessage')
  async sendMessage(@MessageBody() payload: SendMessagePayload) {
    const roomName = payload.chatroomId.toString();

    const map = {
      text: 0,
      image: 1,
      file: 2
    }
    const history = await this.chatHistoryService.add(payload.chatroomId, {
      content: payload.message.content,
      type: map[payload.message.type],
      chatroomId: payload.chatroomId,
      senderId: payload?.sendUserId ?? 0
    });
    const sender = await this.userService.findUserDetailById(history.senderId);

    this.server.to(roomName).emit('message', {
      type: 'sendMessage',
      userId: payload.sendUserId,
      message: {
        ...history,
        sender
      }
    });
  }
}
