import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Access } from "../entities/access.entity";
import { Repository, Like } from 'typeorm';
import { MysqlBaseService } from './mysql-base.service';

@Injectable()
export class AccessService extends MysqlBaseService<Access> {
  constructor(
    @InjectRepository(Access)
    protected accessRepository: Repository<Access>
  ) {
    super(accessRepository);
  }

  async findAll(keyword?: string) {
    const where = keyword ? [
      { name: Like(`%${keyword}%`) }
    ] : {};
    return this.accessRepository.find({ where });
  }

  async findAllWithPagination(page: number, limit: number, keyword?: string) {
    const where = keyword ? [
      { name: Like(`%${keyword}%`) }
    ] : {};
    const [accesses, total] = await this.accessRepository.findAndCount({
      where,
      skip: (page - 1) * limit,
      take: limit
    });
    return { accesses, total };
  }
}
