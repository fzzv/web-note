import { Controller, Get, Post, Body, Res, Session, Redirect } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { UtilityService } from '../../share/services/utility.service';
import type { Response } from 'express';

@Controller('admin')
export class AuthController {
  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService,
  ) { }

  @Get('login')
  showLogin(@Res() res: Response) {
    res.render('auth/login', { layout: false });
  }

  @Post('login')
  async login(@Body() body, @Res() res: Response, @Session() session) {
    const { username, password, captcha } = body;
    if (captcha?.toLowerCase() !== session.captcha?.toLowerCase()) {
      return res.render('auth/login', { message: '验证码错误', layout: false });
    }
    const user = await this.userService.findOne({ where: { username }, relations: ['roles', 'roles.accesses'] });
    if (user && await this.utilityService.comparePassword(password, user.password)) {
      session.user = user;
      return res.redirect('/admin/dashboard');
    } else {
      return res.render('auth/login', { message: '用户名或密码错误', layout: false });
    }
  }

  @Get('captcha')
  getCaptcha(@Res() res: Response, @Session() session) {
    const captcha = this.utilityService.generateCaptcha({ size: 1, ignoreChars: '0o1il' });
    session.captcha = captcha.text;
    res.type('svg');
    res.send(captcha.data);
  }

  @Get('logout')
  @Redirect('login')
  logout(@Session() session) {
    session.user = null;
    return { url: 'login' };
  }
}
