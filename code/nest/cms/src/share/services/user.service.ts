import { Injectable } from '@nestjs/common';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>
  ) {
    super(userRepository);
  }
}
