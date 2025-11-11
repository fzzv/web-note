import { Injectable } from '@nestjs/common';
// 导入 bcrypt 库，用于处理密码哈希和验证
import bcrypt from 'bcrypt';

@Injectable()
export class UtilityService {
  async hashPassword(password: string): Promise<string> {
    const salt = await bcrypt.genSalt();
    return bcrypt.hash(password, salt);
  }
  async comparePassword(password: string, hash: string): Promise<boolean> {
    return bcrypt.compare(password, hash);
  }
}
