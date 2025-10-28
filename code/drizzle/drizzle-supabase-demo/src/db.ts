import { drizzle } from 'drizzle-orm/node-postgres';
import { Client } from 'pg';
import * as schema from './schema';

// 这是一个临时的占位符，稍后我们会替换为 Supabase 的实际连接字符串。
const DATABASE_URL = 'postgres://postgres:postgres@localhost:5432/drizzle'; 

// 1. 创建 pg 客户端
const client = new Client({
  connectionString: DATABASE_URL,
});

// 2. 连接到数据库
// 这一步在实际应用中通常在应用启动时执行一次
client.connect();

// 3. 使用 drizzle(client) 创建 Drizzle 实例
// 传入 schema 使 Drizzle 具备类型感知能力
export const db = drizzle(client, { schema });

// 导出 client 以便后续关闭连接
export const pgClient = client;

// 提醒：在实际项目中，使用连接池 (pg.Pool) 替代 pg.Client 是更好的做法。
