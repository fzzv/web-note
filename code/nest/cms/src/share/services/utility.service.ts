import { Injectable } from '@nestjs/common';
// 导入 bcrypt 库，用于处理密码哈希和验证
import bcrypt from 'bcrypt';
// 导入 svgCaptcha 库，用于生成验证码
import svgCaptcha from 'svg-captcha';

// 使用 Injectable 装饰器将类标记为可注入的服务
@Injectable()
export class UtilityService {
  // 定义一个异步方法，用于生成密码的哈希值
  async hashPassword(password: string): Promise<string> {
    // 生成一个盐值，用于增强哈希的安全性
    const salt = await bcrypt.genSalt();
    // 使用生成的盐值对密码进行哈希，并返回哈希结果
    return bcrypt.hash(password, salt);
  }
  // 定义一个异步方法，用于比较输入的密码和存储的哈希值是否匹配
  async comparePassword(password: string, hash: string): Promise<boolean> {
    // 使用 bcrypt 的 compare 方法比较密码和哈希值，返回比较结果（true 或 false）
    return bcrypt.compare(password, hash);
  }
  generateCaptcha(options) {
    return svgCaptcha.create(options);
  }
}
