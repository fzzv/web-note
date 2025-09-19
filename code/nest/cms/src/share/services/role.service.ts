import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Role } from '../entities/role.entity';
import { Repository, Like, In } from 'typeorm';
import { MysqlBaseService } from "./mysql-base.service";
import { Access } from "../entities/access.entity";
import { UpdateRoleAccessesDto } from "../dtos/role.dto";

@Injectable()
export class RoleService extends MysqlBaseService<Role> {
  constructor(
    @InjectRepository(Role)
    protected roleRepository: Repository<Role>,
    @InjectRepository(Access)
    private readonly accessRepository: Repository<Access>
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
  async findAllWithPagination(page: number = 1, limit: number = 10, search: string = ''): Promise<{ roles: Role[], total: number }> {
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

  async updateAccesses(id: number, updateRoleAccessesDto: UpdateRoleAccessesDto) {
    const role = await this.repository.findOneBy({ id });
    if (!role) throw new Error('Role not found');
    role.accesses = await this.accessRepository.findBy({ id: In(updateRoleAccessesDto.accessIds) });
    await this.repository.update(id, role);
  }
}
