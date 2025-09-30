import { Injectable } from '@nestjs/common';
import { Repository, FindOneOptions, ObjectLiteral, DeepPartial } from 'typeorm';
import { QueryDeepPartialEntity } from 'typeorm/query-builder/QueryPartialEntity.js';

@Injectable()
export abstract class MysqlBaseService<T extends ObjectLiteral> {
  constructor(protected repository: Repository<T>) { }

  async findAll(): Promise<T[]> {
    return this.repository.find();
  }
  async findOne(options: FindOneOptions<T>): Promise<T | null> {
    return this.repository.findOne(options);
  }
  async create(createDto: DeepPartial<T>): Promise<T | T[]> {
    const entity = this.repository.create(createDto);
    return this.repository.save(entity);
  }
  async update(id: number, updateDto: QueryDeepPartialEntity<T>) {
    return await this.repository.update(id, updateDto);
  }
  async delete(id: number) {
    return await this.repository.delete(id);
  }
  async count(): Promise<number> {
    return this.repository.count();
  }
  async findLatest(limit: number) {
    const order: any = {
      id: 'DESC'
    }
    return this.repository.find({
      order,
      take: limit
    });
  }
  async getTrend(tableName): Promise<{ dates: string[]; counts: number[] }> {
    const result = await this.repository.query(`
             SELECT DATE_FORMAT(createdAt, '%Y-%m-%d') as date, COUNT(*) as count
             FROM ${tableName}
             GROUP BY date
             ORDER BY date ASC
         `);
    const dates = result.map(row => row.date);
    const counts = result.map(row => row.count);
    return { dates, counts };
  }
}
