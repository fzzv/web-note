import { HttpException, HttpStatus, Inject, Injectable, Logger } from '@nestjs/common';
import { PrismaService } from 'src/prisma/prisma.service';
import { RedisService } from 'src/service/redis.service';
import { RegisterUserDto } from 'src/dtos/register-user.dto';
import { UtilityService } from 'src/service/utility.service';
import { UpdateUserPasswordDto } from 'src/dtos/update-user-password.dto';
import { UpdateUserDto } from 'src/dtos/udpate-user.dto';

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

  async findUserDetailById(userId: number) {
    const user = await this.prismaService.user.findUnique({
      where: {
        id: userId
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
    return user;
  }

  async updatePassword(passwordDto: UpdateUserPasswordDto) {
    // const captcha = await this.redisService.get(`update_password_captcha_${passwordDto.email}`);

    // if (!captcha) {
    //   throw new HttpException('验证码已失效', HttpStatus.BAD_REQUEST);
    // }

    // if (passwordDto.captcha !== captcha) {
    //   throw new HttpException('验证码不正确', HttpStatus.BAD_REQUEST);
    // }

    const foundUser = await this.prismaService.user.findUnique({
      where: {
        username: passwordDto.username
      }
    });

    if (!foundUser) {
      throw new HttpException('用户不存在', HttpStatus.BAD_REQUEST);
    }

    foundUser.password = passwordDto.password;

    try {
      await this.prismaService.user.update({
        where: {
          id: foundUser.id
        },
        data: foundUser
      });
      return '密码修改成功';
    } catch (e) {
      this.logger.error(e, UserService);
      return '密码修改失败';
    }
  }

  async update(userId: number, updateUserDto: UpdateUserDto) {
    // const captcha = await this.redisService.get(`update_user_captcha_${updateUserDto.email}`);

    // if (!captcha) {
    //   throw new HttpException('验证码已失效', HttpStatus.BAD_REQUEST);
    // }

    // if (updateUserDto.captcha !== captcha) {
    //   throw new HttpException('验证码不正确', HttpStatus.BAD_REQUEST);
    // }

    const foundUser = await this.prismaService.user.findUnique({
      where: {
        id: userId
      }
    });

    if (!foundUser) {
      throw new HttpException('用户不存在', HttpStatus.BAD_REQUEST);
    }

    if (updateUserDto.nickName) {
      foundUser.nickName = updateUserDto.nickName;
    }
    if (updateUserDto.headPic) {
      foundUser.headPic = updateUserDto.headPic;
    }

    try {
      await this.prismaService.user.update({
        where: {
          id: userId
        },
        data: foundUser
      })
      return '用户信息修改成功';
    } catch (e) {
      this.logger.error(e, UserService);
      return '用户信息修改成功';
    }
  }
}
