import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  username: string;

  @Column()
  password: string;

  @Column()
  nickName: string;

  @Column()
  email: string;

  @Column()
  avatar: string;

  @Column()
  createTime: Date;

  @Column()
  updateTime: Date;
}
