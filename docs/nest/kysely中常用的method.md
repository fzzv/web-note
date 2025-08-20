# Kysely 常用方法总结

Kysely 是一个类型安全的 SQL 查询构建器，它通过 TypeScript 的类型系统在编译时捕捉错误，而不是在运行时。以下是 Kysely 中一些核心且常用的方法总结，可以帮助你快速上手和高效使用。

## 核心方法表格

这张表格总结了 Kysely 查询构建中最常用的方法。

| 方法名称 | 基本用法 | 作用 | 备注/适用数据库 |
| :--- | :--- | :--- | :--- |
| **查询构建** | | | |
| `selectFrom` | `db.selectFrom('person')` | 开始一个 `SELECT` 查询，指定查询的主表。 | 所有 |
| `insertInto` | `db.insertInto('person')` | 开始一个 `INSERT` 查询，指定要插入数据的表。 | 所有 |
| `updateTable` | `db.updateTable('person')` | 开始一个 `UPDATE` 查询，指定要更新的表。 | 所有 |
| `deleteFrom` | `db.deleteFrom('person')` | 开始一个 `DELETE` 查询，指定要删除数据的表。 | 所有 |
| **数据筛选与选择** | | | |
| `where` | `.where('id', '=', 1)` | 添加 `WHERE` 条件来筛选数据。可以链式调用添加多个 `AND` 条件。 | 所有 |
| `select` | `.select(['id', 'first_name'])` | 指定查询需要返回的列。 | 所有 |
| `selectAll` | `.selectAll()` 或 `.selectAll('person')` | 选择一个表或所有连接表中的所有列。 | 所有 |
| **数据操作** | | | |
| `values` | `.values({ first_name: 'John' })` | 在 `INSERT` 语句中指定要插入的一行或多行数据。 | 所有 |
| `set` | `.set({ first_name: 'Jane' })` | 在 `UPDATE` 语句中指定要更新的列和对应的值。 | 所有 |
| **关联与分组** | | | |
| `innerJoin`<br>`leftJoin`<br>`rightJoin` | `.innerJoin('pet', 'pet.owner_id', 'person.id')` | 将其他表与主表进行连接（JOIN）。 | 所有 |
| `orderBy` | `.orderBy('age', 'desc')` | 对查询结果进行排序。 | 所有 |
| `groupBy` | `.groupBy('gender')` | 对结果集进行分组，通常与聚合函数一起使用。 | 所有 |
| `having` | `.having(db.fn.count('id'), '>', 10)` | 在 `GROUP BY` 分组后对结果进行筛选。 | 所有 |
| **结果处理** | | | |
| `limit` | `.limit(10)` | 限制查询返回的行数。 | 所有 |
| `offset` | `.offset(20)` | 设置查询结果的偏移量，通常与 `limit` 配合用于分页。 | 所有 |
| `returning` | `.returning('id')` | 在 `INSERT`, `UPDATE`, `DELETE` 后返回指定的列。 | 主要适用于 PostgreSQL, SQLite (部分版本), MSSQL。MySQL 需要特定配置。 |
| **查询执行** | | | |
| `execute` | `.execute()` | 执行构建好的查询，并返回一个结果数组。 | 所有 |
| `executeTakeFirst` | `.executeTakeFirst()` | 执行查询并返回第一个结果。如果查询没有返回任何行，则返回 `undefined`。 | 所有 |
| `executeTakeFirstOrThrow` | `.executeTakeFirstOrThrow()` | 执行查询并返回第一个结果。如果查询没有返回任何行，则会抛出 `NoResultError` 异常。 | 所有 |
| **事务** | | | |
| `transaction` | `db.transaction().execute(async (trx) => { ... })` | 执行一个数据库事务。在回调函数 `execute` 中的所有查询都会在同一个事务中运行。如果回调函数抛出异常，事务将自动回滚。 | 所有支持事务的数据库。 |
| **原生SQL与函数** | | | |
| `sql` | `sql`...` ` (模板字符串) | 用于编写原生 SQL 片段，可以安全地将参数传递给查询，有效防止 SQL 注入。 | 所有 |
| `fn` | `db.fn.count('id').as('count')` | 调用 SQL 内置函数，如 `COUNT`, `AVG`, `SUM`, `MAX`, `MIN` 等聚合函数，或其他数据库特定函数。 | 所有 |
| `onConflict` | `.onConflict((oc) => oc.column('email').doUpdateSet({ ... }))` | 处理插入冲突（"upsert" 操作）。 | 主要适用于 PostgreSQL 和 SQLite。 |
| `onDuplicateKeyUpdate` | `.onDuplicateKeyUpdate({ ... })` | 在 MySQL 中处理插入时的主键或唯一键冲突。 | MySQL / MariaDB |

## 其他补充

### 1. Kysely 实例 (`db`)
`db` 对象是 `Kysely` 类的实例，是所有查询的入口点。它持有了数据库连接池和方言（dialect）配置。

### 2. 类型安全与数据库接口
Kysely 的核心优势在于其类型安全性。你需要为你的数据库结构定义一个 TypeScript 接口，Kysely 会利用这个接口来检查你的查询是否合法。

```typescript
import { Generated } from 'kysely';

interface Person {
  id: Generated<number>;
  first_name: string;
  last_name: string | null;
  age: number;
}

interface Database {
  person: Person;
  // ... other tables
}

// const db = new Kysely<Database>({ ... });
```
通过这种方式，如果你尝试查询一个不存在的列（例如 `db.selectFrom('person').select('non_existent_column')`），TypeScript 会在编译时报错。

### 3. 方言 (Dialects)
Kysely 通过不同的方言包来支持多种数据库，例如：
- `kysely-postgres` for PostgreSQL
- `mysql2` (Kysely 内置支持) for MySQL
- `better-sqlite3` (Kysely 内置支持) for SQLite

你需要在初始化 Kysely 实例时提供正确的方言配置。

### 4. 链式调用
Kysely 的 API 设计是高度可链式调用的。你可以从一个查询起点（如 `selectFrom`）开始，然后像流水线一样不断追加操作（如 `where`, `orderBy`, `limit`），最后使用执行方法（如 `execute`）来结束并运行查询。这种设计使得代码既清晰又易于组合。
