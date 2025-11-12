import { FriendAddDto } from 'src/dtos/friend-add.dto';
import { Injectable } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class FriendshipService {
  constructor(private readonly prismaService: PrismaService) { }

  // 获取好友列表
  async getFriendship(userId: number) {
    const friendships = await this.prismaService.friendship.findMany({
      where: {
        OR: [
          {
            userId: userId
          },
          {
            friendId: userId
          }
        ]
      }
    });

    const set = new Set<number>();
    for (const friendship of friendships) {
      set.add(friendship.userId);
      set.add(friendship.friendId);
    }

    const friendIds = [...set].filter(id => id !== userId);

    const res: any[] = [];

    for (const friendId of friendIds) {
      const friend = await this.prismaService.user.findUnique({
        where: {
          id: friendId
        },
        select: {
          id: true,
          username: true,
          nickName: true,
          email: true
        }
      });
      res.push(friend);
    }
    return res;
  }

  // 添加好友
  async add(friendAddDto: FriendAddDto, userId: number) {
    return await this.prismaService.friendRequest.create({
      data: {
        fromUserId: userId,
        toUserId: friendAddDto.friendId,
        reason: friendAddDto.reason,
        status: 0
      }
    })
  }

  // 获取好友请求列表
  async list(userId: number) {
    return await this.prismaService.friendRequest.findMany({
      where: {
        fromUserId: userId
      }
    })
  }

  // 同意好友请求
  async agree(friendId: number, userId: number) {
    await this.prismaService.friendRequest.updateMany({
      where: {
        fromUserId: friendId,
        toUserId: userId,
        status: 0
      },
      data: {
        status: 1
      }
    })
    const res = await this.prismaService.friendship.findMany({
      where: {
        userId,
        friendId
      }
    })

    // 如果好友关系不存在，则创建好友关系
    if (!res.length) {
      await this.prismaService.friendship.create({
        data: {
          userId,
          friendId
        }
      })
    }
    return '添加成功'
  }

  // 拒绝好友请求
  async reject(friendId: number, userId: number) {
    await this.prismaService.friendRequest.updateMany({
      where: {
        fromUserId: friendId,
        toUserId: userId,
        status: 0
      },
      data: {
        status: 2
      }
    })
    return '已拒绝'
  }

  // 移除好友
  async remove(friendId: number, userId: number) {
    await this.prismaService.friendship.deleteMany({
      where: {
        userId,
        friendId
      }
    })
    return '删除成功'
  }
}
