import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { match } from 'path-to-regexp';
import { AccessService } from '../../share/services/access.service';
import { Access } from '../../share/entities/access.entity';
import { AccessType } from '../../share/dtos/access.dto';

@Injectable()
export class AuthMiddleware implements NestMiddleware {
  constructor(private readonly accessService: AccessService) { }
  async use(req: Request, res: Response, next: NextFunction) {
    const user = req.session.user;
    // const user = { is_super: true, roles: [] }
    if (!user) {
      return res.redirect('/admin/login');
    }

    res.locals.user = user;

    const accessTree = await this.accessService.findAll();
    const userAccessIds = this.getUserAccessIds(user);
    res.locals.menuTree = user.is_super ? accessTree : this.getMenuTree(accessTree, userAccessIds);
    if (user.is_super || req.originalUrl === '/admin/dashboard') {
      return next();
    }

    if (this.hasPermission(user, req.originalUrl)) {
      return next();
    } else {
      res.status(403).render('error', { message: '无权限访问此页面', layout: false });
    }
  }

  private getUserAccessIds(user): number[] {
    return user.roles.flatMap(role => role.accesses.map(access => access.id));
  }

  private getMenuTree(accessTree: Access[], userAccessIds: number[]): Access[] {
    return accessTree.filter(access => {
      if (access.type === AccessType.FEATURE || !userAccessIds.includes(access.id)) {
        return false;
      }
      if (access.children) {
        access.children = this.getMenuTree(access.children, userAccessIds);
      }
      return true;
    });
  }

  private hasPermission(user, url: string): boolean {
    const userAccessUrls = user.roles.flatMap(role => role.accesses.map(access => access.url));
    return userAccessUrls.some(urlPattern => match(urlPattern)(url));
  }
}
