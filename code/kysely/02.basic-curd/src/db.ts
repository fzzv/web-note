import { Database } from './schema'
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

export const db = new Kysely<Database>({
  dialect,
})
