import { ApiHideProperty, ApiProperty } from '@nestjs/swagger';
import { Exclude, Transform } from 'class-transformer';
import { Entity, Column, PrimaryGeneratedColumn, CreateDateColumn, UpdateDateColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  @ApiProperty({ description: '用户ID', example: 1 })
  id: number;

  @Column({ length: 50, unique: true })
  @ApiProperty({ description: '用户名', example: 'admin' })
  username: string;

  @Column()
  @Exclude() // 在序列化时排除密码字段，不返回给前端
  @ApiHideProperty() // 隐藏密码字段，不在Swagger文档中显示
  password: string;

  @Column({ length: 15, nullable: true })
  @ApiProperty({ description: '手机号', example: '13124567890', format: '手机号码会被部分隐藏' })
  @Transform(({ value }) => value ? value.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2') : value)
  mobile: string;

  @Column({ length: 100, nullable: true })
  @ApiProperty({ description: '邮箱', example: 'admin@example.com' })
  email: string;

  @Column({ default: 1 })
  @ApiProperty({ description: '状态', example: 1, enum: [1, 2] })
  status: number;

  @Column({ default: false })
  @ApiProperty({ description: '是否超级管理员', example: false })
  is_super: boolean;

  @Column({ default: 100 })
  @ApiProperty({ description: '排序', example: 100 })
  sort: number;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '创建时间', example: '2021-01-01 00:00:00' })
  @CreateDateColumn()
  createdAt: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP', onUpdate: 'CURRENT_TIMESTAMP' })
  @ApiProperty({ description: '更新时间', example: '2021-01-01 00:00:00' })
  @UpdateDateColumn()
  updatedAt: Date;
}
