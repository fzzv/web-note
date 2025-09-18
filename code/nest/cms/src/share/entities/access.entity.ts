import { ApiProperty } from '@nestjs/swagger';
import { Tree, Entity, Column, PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn, TreeChildren, TreeParent } from 'typeorm';
import { AccessType } from '../dtos/access.dto';

@Entity()
@Tree("materialized-path")
export class Access {
  @PrimaryGeneratedColumn()
  @ApiProperty({ description: 'ID', example: 1 })
  id: number;

  @Column({ length: 50, unique: true })
  @ApiProperty({ description: '名称', example: 'name' })
  name: string;

  @Column({ type: 'enum', enum: AccessType })
  type: AccessType;

  @Column({ length: 200, nullable: true })
  url: string;

  @Column({ length: 200, nullable: true })
  description: string;

  @TreeChildren()
  children: Access[];

  @TreeParent()
  parent: Access;

  @Column({ default: 1 })
  @ApiProperty({ description: '生效状态', example: 1 })
  status: number;

  @Column({ default: 100 })
  @ApiProperty({ description: '排序号', example: 100 })
  sort: number;

  @CreateDateColumn()
  @ApiProperty({ description: '创建时间', example: '2020年1月11日16:49:22' })
  createdAt: Date;

  @UpdateDateColumn()
  @ApiProperty({ description: '更新时间', example: '2020年1月11日16:49:22' })
  updatedAt: Date;
}
