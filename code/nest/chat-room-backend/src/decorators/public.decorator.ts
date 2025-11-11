import { SetMetadata } from '@nestjs/common';

export const IS_PUBLIC_KEY = 'isPublic';

// 公开接口装饰器 在不需要登录的接口使用 @Public() 装饰器
export const Public = () => SetMetadata(IS_PUBLIC_KEY, true);
