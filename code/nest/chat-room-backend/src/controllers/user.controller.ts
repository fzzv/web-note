import { Controller, Post, Body, Get, Query, HttpException, HttpStatus, BadRequestException, Res } from '@nestjs/common';
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
import { UserInfo } from 'src/decorators';
import { UpdateUserPasswordDto } from 'src/dtos/update-user-password.dto';
import { UpdateUserDto } from 'src/dtos/udpate-user.dto';

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
    if (!user) {
      throw new HttpException('用户名或密码错误', HttpStatus.INTERNAL_SERVER_ERROR)
    }
    return {
      user,
      token: this.createJwtTokens(user)
    };
  }

  // 刷新 token
  @Public()
  @Post('refresh-token')
  async refreshToken(@Body() body) {
    const { refresh_token } = body;
    try {
      const decoded = this.jwtService.verify(refresh_token, { secret: this.configurationService.jwtSecret });
      const tokens = this.createJwtTokens(decoded);
      return { ...tokens };
    } catch (error) {
      return { message: 'Refresh token无效或已过期' };
    }
  }

  @Get('info')
  async info(@UserInfo('id') userId: number) {
    return this.userService.findUserDetailById(userId);
  }

  @Post('update_password')
  async updatePassword(@Body() passwordDto: UpdateUserPasswordDto) {
    // 给密码加密
    passwordDto.password = await this.utilityService.hashPassword(passwordDto.password);
    return this.userService.updatePassword(passwordDto);
  }

  @Get('update_password/captcha')
  async updatePasswordCaptcha(@Query('address') address: string) {
    if (!address) {
      throw new BadRequestException('邮箱地址不能为空');
    }
    const code = Math.random().toString().slice(2, 8);

    await this.redisService.set(`update_password_captcha_${address}`, code, 10 * 60);

    await this.emailService.sendMail({
      to: address,
      subject: '更改密码验证码',
      html: `<p>你的更改密码验证码是 ${code}</p>`
    });
    return '发送成功';
  }

  @Post('update')
  async update(@UserInfo('id') userId: number, @Body() updateUserDto: UpdateUserDto) {
    return await this.userService.update(userId, updateUserDto);
  }

  @Get('update/captcha')
  async updateCaptcha(@Query('address') address: string) {
    if (!address) {
      throw new BadRequestException('邮箱地址不能为空');
    }
    const code = Math.random().toString().slice(2, 8);

    await this.redisService.set(`update_user_captcha_${address}`, code, 10 * 60);

    await this.emailService.sendMail({
      to: address,
      subject: '更改用户信息验证码',
      html: `<p>你的验证码是 ${code}</p>`
    });
    return '发送成功';
  }

  async validateUser(username: string, password: string) {
    const user = await this.prismaService.user.findUnique({
      where: {
        username
      },
      select: {
        id: true,
        username: true,
        password: true,
        nickName: true,
        email: true,
        headPic: true,
        createTime: true
      }
    });
    if (user && await this.utilityService.comparePassword(password, user.password)) {
      user.password = '********';
      return user;
    }
    return null;
  }

  private createJwtTokens(user: any) {
    // 创建访问令牌，设置过期时间为 30 分钟
    const access_token = this.jwtService.sign({ id: user.id, username: user.username }, {
      secret: this.configurationService.jwtSecret,
      expiresIn: '10m',
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
