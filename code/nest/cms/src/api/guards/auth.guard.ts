// 导入所需的装饰器、模块和服务
import { Injectable, CanActivate, ExecutionContext, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { UserService } from '../../share/services/user.service';
import { ConfigurationService } from '../../share/services/configuration.service';
import { Request } from 'express';

// 使用 @Injectable() 装饰器将此类标记为可注入的服务
@Injectable()
export class AuthGuard implements CanActivate {
  // 构造函数，注入所需的服务
  constructor(
    private readonly userService: UserService,
    private readonly jwtService: JwtService,
    private readonly configurationService: ConfigurationService,
  ) { }
  // 实现 CanActivate 接口的 canActivate 方法，用于进行身份验证
  async canActivate(
    context: ExecutionContext,
  ): Promise<boolean> {
    // 获取 HTTP 请求对象
    const request = context.switchToHttp().getRequest<Request>();
    // 从请求头中提取令牌
    const token = this.extractTokenFromHeader(request);

    // 如果没有令牌，抛出未授权异常
    if (!token) {
      throw new UnauthorizedException('Token not provided');
    }
    try {
      // 验证令牌并获取解码后的数据
      const decoded = this.jwtService.verify(token, { secret: this.configurationService.jwtSecret });
      // 查找用户并获取其关联的角色和权限
      const user = await this.userService.findOne({ where: { id: decoded.id }, relations: ['roles', 'roles.accesses'] });
      // 如果找到用户
      if (user) {
        // 删除用户密码，防止暴露敏感信息
        if (user.password) {
          user.password = '******';
        }
        // 将用户信息附加到请求对象
        request.user = user;
        // 返回 true，表示通过身份验证
        return true;
      } else {
        // 如果用户不存在，抛出未授权异常
        throw new UnauthorizedException('User not found');
      }
    } catch (error) {
      // 捕获令牌验证错误并抛出未授权异常
      throw new UnauthorizedException('Invalid or expired token');
    }
  }
  // 从请求头中提取令牌的私有方法
  private extractTokenFromHeader(request: Request): string | undefined {
    // 从请求头中获取授权信息，并按照空格分隔为类型和令牌
    const [type, token] = request.headers.authorization?.split(' ') ?? [];
    // 如果类型为 'Bearer'，则返回令牌，否则返回 undefined
    return type === 'Bearer' ? token : undefined;
  }
}
