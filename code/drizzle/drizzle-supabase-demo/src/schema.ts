import {
  pgTable, // 用于定义 PostgreSQL 表
  serial, // 自动递增的整数 (ID)
  text, // 变长文本/字符串
  varchar, // 限制长度的字符串
  timestamp, // 时间戳
  unique, // 唯一约束
} from 'drizzle-orm/pg-core';
import { relations } from 'drizzle-orm';

// --- 1. Users Table (用户表) ---
export const users = pgTable('users', {
  id: serial('id').primaryKey(), // 主键，自动递增
  
  // 使用 varchar(256) 限制长度
  name: varchar('name', { length: 256 }).notNull(),
  
  // email 必须唯一且非空，这是登录的关键字段
  email: varchar('email', { length: 256 }).notNull().unique(),
  
  // 密码存储为哈希值，使用 text 类型
  password: text('password').notNull(), 
  
  // 头像 URL，可以为空
  avatarUrl: varchar('avatar_url', { length: 512 }),
  
  // 记录创建时间
  createdAt: timestamp('created_at').defaultNow().notNull(),
});

// --- 2. Posts Table (帖子表) ---
export const posts = pgTable('posts', {
  id: serial('id').primaryKey(), // 主键
  title: varchar('title', { length: 256 }).notNull(),
  content: text('content').notNull(),
  
  // 外键：作者ID，关联到 users 表的 id 字段
  authorId: serial('author_id')
    .references(() => users.id, { onDelete: 'cascade' }) // 设定关联和级联删除
    .notNull(),

  published: timestamp('published_at'), // 发布时间，可以为空
  createdAt: timestamp('created_at').defaultNow().notNull(),
});


// --- 3. Relations (定义关系) ---
// 这对于 Drizzle ORM 的 join 和 Eager Loading 非常重要
export const usersRelations = relations(users, ({ many }) => ({
  // 一个用户可以有多个帖子 (一对多)
  posts: many(posts), 
}));

export const postsRelations = relations(posts, ({ one }) => ({
  // 一篇帖子只属于一个作者 (多对一)
  author: one(users, {
    fields: [posts.authorId], // posts.authorId 字段指向 users 表
    references: [users.id],
  }),
}));
