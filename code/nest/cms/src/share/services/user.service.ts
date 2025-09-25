import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { In, Like, Repository } from 'typeorm';
import { MysqlBaseService } from './mysql-base.service';
import { User } from '../entities/user.entity';
import { Role } from '../entities/role.entity';
import { UpdateUserRolesDto } from '../dtos/user.dto';

@Injectable()
export class UserService extends MysqlBaseService<User> {
  constructor(
    @InjectRepository(User)
    protected userRepository: Repository<User>,
    @InjectRepository(Role)
    protected roleRepository: Repository<Role>
  ) {
    super(userRepository);
  }

  async findAll(search: string = ''): Promise<User[]> {
    const where = search ? [
      { username: Like(`%${search}%`) },
      { email: Like(`%${search}%`) }
    ] : {};

    const users = await this.userRepository.find({
      where
    });
    return users;
  }

  async findAllWithPagination(page: number = 1, limit: number = 10, search: string = ''): Promise<{ users: User[], total: number }> {
    const where = search ? [
      { username: Like(`%${search}%`) },
      { email: Like(`%${search}%`) }
    ] : {};

    const [users, total] = await this.userRepository.findAndCount({
      where,
      skip: (page - 1) * limit,
      take: limit,
    });
    return { users, total };
  }

  async updateRoles(id: number, updateUserRolesDto: UpdateUserRolesDto) {
    const user = await this.repository.findOneBy({ id });
    if (!user) throw new Error('User not found');
    user.roles = await this.roleRepository.findBy({ id: In(updateUserRolesDto.roleIds) });
    await this.repository.save(user);
  }
}
