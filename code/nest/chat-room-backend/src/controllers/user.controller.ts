import { Controller, Post, Body, Get, Query } from '@nestjs/common';
import { UserService } from 'src/service/user.service';
import { RedisService } from 'src/service/redis.service';
import { EmailService } from 'src/service/email.service';
import { RegisterUserDto } from 'src/dtos/register-user.dto';

@Controller('user')
export class UserController {
  constructor(
    private readonly userService: UserService,
    private readonly redisService: RedisService,
    private readonly emailService: EmailService
  ) { }

  @Post('register')
  async createUser(@Body() user: RegisterUserDto) {
    return this.userService.register(user);
  }

  @Get('register-captcha')
  async captcha(@Query('address') address: string) {
    const code = Math.random().toString().slice(2, 8);

    await this.redisService.set(`captcha_${address}`, code, 5 * 60);

    await this.emailService.sendMail({
      to: address,
      subject: '注册验证码',
      html: `<p>你的注册验证码是 ${code}</p>`
    });
    return '发送成功';
  }
}
