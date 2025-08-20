import { Database } from './types'
import { createPool } from 'mysql2'
import { Kysely, MysqlDialect, MysqlPool } from 'kysely'

const dialect = new MysqlDialect({
  pool: createPool({
    database: 'kysely',
    host: 'localhost',
    user: 'root',
    password: 'Fan0124.',
    port: 3306,
    connectionLimit: 10,
  }) as unknown as MysqlPool
})

// Database interface is passed to Kysely's constructor, and from now on, Kysely 
// knows your database structure.
// Dialect is passed to Kysely's constructor, and from now on, Kysely knows how 
// to communicate with your database.
/**
 * 将 `Database` 接口传递给 Kysely 的构造函数，从现在开始，Kysely 就知道了你的数据库结构。
 * 将 `dialect` 传递给 Kysely 的构造函数，从现在开始，Kysely 就知道如何与你的数据库通信。
 */
export const db = new Kysely<Database>({
  dialect,
})
