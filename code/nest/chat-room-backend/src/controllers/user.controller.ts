import { Controller, Post, Body, Get, Query, HttpException, HttpStatus } from '@nestjs/common';
import { UserService } from 'src/service/user.service';
import { RedisService } from 'src/service/redis.service';
import { EmailService } from 'src/service/email.service';
import { RegisterUserDto } from 'src/dtos/register-user.dto';
import { LoginUserDto } from 'src/dtos/login-user.dto';
import { PrismaService } from 'src/prisma/prisma.service';
import { UtilityService } from 'src/service/utility.service';
import { ConfigurationService } from 'src/service/configuration.service';
import { JwtService } from '@nestjs/jwt';
import { Public } from 'src/decorators/public.decorator';

@Controller('user')
export class UserController {
  constructor(
    private readonly userService: UserService,
    private readonly prismaService: PrismaService,
    private readonly redisService: RedisService,
    private readonly emailService: EmailService,
    private readonly utilityService: UtilityService,
    private readonly jwtService: JwtService,
    private readonly configurationService: ConfigurationService
  ) { }

  @Public()
  @Post('register')
  async createUser(@Body() user: RegisterUserDto) {
    return this.userService.register(user);
  }

  @Public()
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

  @Public()
  @Post('login')
  async login(@Body() body: LoginUserDto) {
    const { username, password } = body;
    const user = await this.validateUser(username, password);
    if (user) {
      return {
        user: username,
        token: this.createJwtTokens(user)
      }
    }
    return null;
  }

  async validateUser(username: string, password: string) {
    const user = await this.prismaService.user.findUnique({
      where: {
        username
      }
    });
    if (user && await this.utilityService.comparePassword(password, user.password)) {
      return user;
    }
    return null;
  }

  private createJwtTokens(user: any) {
    // 创建访问令牌，设置过期时间为 30 分钟
    const access_token = this.jwtService.sign({ id: user.id, username: user.username }, {
      secret: this.configurationService.jwtSecret,
      expiresIn: '30m',
    });
    // 创建刷新令牌，设置过期时间为 7 天
    const refresh_token = this.jwtService.sign({ id: user.id, username: user.username }, {
      secret: this.configurationService.jwtSecret,
      expiresIn: '7d',
    });
    // 返回令牌信息
    return { access_token, refresh_token };
  }
}
