import { HttpException, HttpStatus, Inject, Injectable, Logger } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { RedisService } from 'src/service/redis.service';
import { RegisterUserDto } from 'src/dtos/register-user.dto';
import { UtilityService } from 'src/service/utility.service';

@Injectable()
export class UserService {
  @Inject(PrismaService)
  private prismaService: PrismaService;

  @Inject(RedisService)
  private redisService: RedisService;

  @Inject(UtilityService)
  private utilityService: UtilityService;

  private logger = new Logger();

  async register(user: RegisterUserDto) {
    // const captcha = await this.redisService.get(`captcha_${user.email}`);

    // if (!captcha) {
    //   throw new HttpException('验证码已失效', HttpStatus.BAD_REQUEST);
    // }

    // if (user.captcha !== captcha) {
    //   throw new HttpException('验证码不正确', HttpStatus.BAD_REQUEST);
    // }

    const foundUser = await this.prismaService.user.findUnique({
      where: {
        username: user.username
      }
    });

    if (foundUser) {
      throw new HttpException('用户已存在', HttpStatus.BAD_REQUEST);
    }

    let hashedPassword = await this.utilityService.hashPassword(user.password);
    try {
      return await this.prismaService.user.create({
        data: {
          username: user.username,
          password: hashedPassword,
          nickName: user.nickName,
          email: user.email
        },
        select: {
          id: true,
          username: true,
          nickName: true,
          email: true,
          headPic: true,
          createTime: true
        }
      });
    } catch (e) {
      this.logger.error(e, UserService);
      return null;
    }
  }
}
