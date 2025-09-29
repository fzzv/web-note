// 导入 Injectable 装饰器，用于标记一个服务类
import { Injectable } from '@nestjs/common';
// 导入 ConfigService，用于获取配置文件中的配置信息
import { ConfigService } from '@nestjs/config';
// 导入 COS SDK
import COS from 'cos-nodejs-sdk-v5';
// 使用 Injectable 装饰器将 CosService 标记为可注入的服务
@Injectable()
export class CosService {
  // 定义一个私有变量，用于存储 COS 实例
  private cos: COS;
  // 构造函数，注入 ConfigService 以获取配置信息
  constructor(private readonly configService: ConfigService) {
    // 初始化 COS 实例，使用配置服务中的 SecretId 和 SecretKey
    this.cos = new COS({
      SecretId: this.configService.get('COS_SECRET_ID'),
      SecretKey: this.configService.get('COS_SECRET_KEY'),
    });
  }
  // 获取签名认证信息的方法，默认过期时间为 60 秒
  getAuth(key, expirationTime = 60) {
    // 从配置服务中获取 COS 存储桶名称和区域
    const bucket = this.configService.get('COS_BUCKET');
    const region = this.configService.get('COS_REGION');
    // 获取 COS 签名，用于 PUT 请求
    const sign = this.cos.getAuth({
      Method: 'put', // 请求方法为 PUT
      Key: key, // 文件的对象键（路径）
      Expires: expirationTime, // 签名的有效期
      Bucket: bucket, // 存储桶名称
      Region: region, // 存储桶所在区域
    });
    // 返回包含签名、键名、存储桶和区域的信息对象
    return {
      sign,
      key: key,
      bucket,
      region,
    };
  }
}
