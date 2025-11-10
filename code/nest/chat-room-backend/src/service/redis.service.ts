import { Injectable, OnModuleDestroy } from '@nestjs/common';
import { createClient } from "redis"
import { ConfigurationService } from 'src/service/configuration.service';

@Injectable()
export class RedisService implements OnModuleDestroy {
  private redisClient
  constructor(private configurationService: ConfigurationService) {
    this.redisClient = createClient({
      socket: {
        port: this.configurationService.redisPort,
        host: this.configurationService.redisHost,
      }
    })
    this.redisClient.connect().catch(console.error)
  }
  onModuleDestroy() { // 当模块销毁的时候退出当前的客户端
    this.redisClient.quit();
  }
  getClient() {
    return this.redisClient;
  }
  async set(key: string, value: string, ttl?: number) {
    if (ttl) {
      await this.redisClient.set(key, value, 'EX', ttl)
    } else {
      await this.redisClient.set(key, value);
    }
  }
  async get(key: string) {
    return this.redisClient.get(key);
  }
  async del(key: string) {
    await this.redisClient.del(key)
  }
}
