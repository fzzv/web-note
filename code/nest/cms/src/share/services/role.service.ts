import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Role } from '../entities/role.entity';
import { Repository, Like } from 'typeorm';
import { MysqlBaseService } from "./mysql-base.service";

@Injectable()
export class RoleService extends MysqlBaseService<Role> {
  constructor(
    @InjectRepository(Role)
    protected roleRepository: Repository<Role>
  ) {
    super(roleRepository);
  }
  async findAll(search: string = ''): Promise<Role[]> {
    const where = search ? [
       { name: Like(`%${search}%`) },
    ] : {};

    const roles = await this.roleRepository.find({
      where
    });
     return roles;
  }
  async findAllWithPagination(page: number = 1, limit: number = 10, search: string = ''): Promise<{ roles:Role[], total: number }> {
    const where = search ? [
      { name: Like(`%${search}%`) },
    ] : {};

    const [roles, total] = await this.roleRepository.findAndCount({
      where,
      skip: (page - 1) * limit,
      take: limit,
    });
     return { roles, total };
  }
}
