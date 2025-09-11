import { Entity, Column, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ length: 50, unique: true })
  username: string;

  @Column()
  password: string;

  @Column({ length: 15, nullable: true })
  mobile: string;

  @Column({ length: 100, nullable: true })
  email: string;

  @Column({ default: 1 })
  status: number;

  @Column({ default: false })
  is_super: boolean;

  @Column({ default: 100 })
  sort: number;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP' })
  createdAt: Date;

  @Column({ type: 'timestamp', default: () => 'CURRENT_TIMESTAMP', onUpdate: 'CURRENT_TIMESTAMP' })
  updatedAt: Date;
}
