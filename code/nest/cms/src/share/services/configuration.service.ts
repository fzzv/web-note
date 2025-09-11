import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ConfigurationService {
  constructor(private configService: ConfigService) { }
  get mysqlHost(): string {
    return this.configService.get<string>('MYSQL_HOST')!;
  }
  get mysqlPort(): number {
    return this.configService.get<number>('MYSQL_PORT')!;
  }
  get mysqlDb(): string {
    return this.configService.get<string>('MYSQL_DB')!;
  }
  get mysqlUser(): string {
    return this.configService.get<string>('MYSQL_USER')!;
  }
  get mysqlPass(): string {
    return this.configService.get<string>('MYSQL_PASSWORD')!;
  }
  get mysqlConfig() {
    return {
      host: this.mysqlHost,
      port: this.mysqlPort,
      database: this.mysqlDb,
      username: this.mysqlUser,
      password: this.mysqlPass,
    };
  }
}
