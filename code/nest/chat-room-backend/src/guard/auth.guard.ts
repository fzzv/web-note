import { CanActivate, ExecutionContext, Injectable, UnauthorizedException } from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { JwtService } from "@nestjs/jwt";
import { Request } from "express";
import { ConfigurationService } from "src/service/configuration.service";
import { PrismaService } from 'src/prisma/prisma.service';
import { IS_PUBLIC_KEY } from 'src/decorators/public.decorator';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(
    private readonly configurationService: ConfigurationService,
    private readonly jwtService: JwtService,
    private readonly prismaService: PrismaService,
    private readonly reflector: Reflector
  ) { }

  async canActivate(context: ExecutionContext): Promise<boolean> {

    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);

    // 如果标记为公开接口，直接放行
    if (isPublic) return true;

    const request = context.switchToHttp().getRequest<Request>();
    const token = this.extractTokenFromHeader(request);

    if (!token) {
      throw new UnauthorizedException('Token not provided');
    }

    try {
      const decoded = this.jwtService.verify(token, { secret: this.configurationService.jwtSecret });
      const user = await this.prismaService.user.findUnique({
        where: {
          id: decoded.id
        }
      });
      if (!user) {
        throw new UnauthorizedException('User not found');
      }
      // 删除用户密码，防止暴露敏感信息
      if (user.password) {
        user.password = '******';
      }
      console.log(user, 'user');
      request.user = user;
      return true;
    } catch (error) {
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
