# kysely查询进阶

准备工作，创建一个users表和posts表，用于多表查询

```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) NOT NULL UNIQUE,
  status ENUM('active', 'inactive') NOT NULL DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE posts (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  content TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入几条数据
INSERT INTO users (name, email, status) VALUES
('Alice', 'alice@example.com', 'active'),
('Bob', 'bob@example.com', 'inactive'),
('Charlie', 'charlie@example.com', 'active');

INSERT INTO posts (user_id, title, content) VALUES
(1, 'Hello World', 'This is Alice\'s first post.'),
(1, 'Learning Kysely', 'Alice is learning how to use Kysely with MySQL.'),
(2, 'Bob\'s Thoughts', 'Bob is inactive but still writing posts.'),
(3, 'Charlie Post 1', 'Charlie shares something interesting.'),
(3, 'Charlie Post 2', 'Another post from Charlie.');
```

## 基础查询

```ts
/**
 * GET /users?page=1&pageSize=10&name=A
 * 分页获取用户
 */
router.get('/', async (req, res) => {
  const page = parseInt((req.query.page as string) || '1', 10);
  const pageSize = parseInt((req.query.pageSize as string) || '10', 10);

  const users = await db
    .selectFrom('users')
    .selectAll()
    .where('name', 'like', `%${req.query.name}%`)
    .orderBy('id', 'desc')
    .limit(pageSize)
    .offset((page - 1) * pageSize)
    .execute();

  res.json(users);
});
```

基本用法：`selectFrom().selectAll().execute()` 

条件查询：`.where('id', '=', 1)`

排序：`.orderBy('id', 'desc')`

分页：`.limit(10).offset(20)` 

## 条件 & 组合查询

```ts
/**
 * GET /users/ids?ids=1,2
 * 查询 id 在 [1,2] 的用户
 */
router.get('/ids', async (req, res) => {
  const { ids } = req.query;
  const idsArray = (ids as string).split(',').map(Number);
  const users = await db
    .selectFrom('users')
    .selectAll()
    .where('id', 'in', idsArray)
    .execute();
  res.json(users);
});
```

`where` + 逻辑运算（`and`, `or`）`in`, `like`, `between`

## 聚合查询

```ts
/**
 * GET /users/stats
 * 聚合统计
 */
router.get('/stats/summary', async (req, res) => {
  const stats = await db
    .selectFrom('users')
    .select((eb) => [
      eb.fn.count<number>('id').as('total'),
      eb.fn.count<number>(eb.case().when('status', '=', 'active').then(1).end()).as('activeCount'),
      eb.fn.count<number>(eb.case().when('status', '=', 'inactive').then(1).end()).as('inactiveCount'),
    ])
    .executeTakeFirst();

  res.json(stats);
});

/**
 * "total": 3,
 * "activeCount": 2,
 * "inactiveCount": 1
 */
// 也可以使用 .groupBy('status') 进行分组统计
/**
 * GET /users/stats/group-by
 * 分组统计
 */
router.get('/stats/group-by', async (req, res) => {
  const stats = await db
    .selectFrom('users')
    .select((eb) => [eb.fn.count<number>('id').as('total'), 'status'])
    .groupBy('status')
    .execute();
  res.json(stats);
});
```

指定字段查询：`.select(['id', 'name'])`

创建动态字段：`.select(expr => ...)` 

`.select(db.fn.count('id').as('count'))`

``sum`, `avg`, `max`, `min`

分组：`.groupBy('status')`

单条结果：`.executeTakeFirst()` / `.executeTakeFirstOrThrow()`

## 多表查询

```ts
/**
 * GET /users/with-posts
 * 用户 + 文章（左连接）
 */
router.get('/with-posts', async (_req, res) => {
  const result = await db
    .selectFrom('users')
    .leftJoin('posts', 'users.id', 'posts.user_id')
    .select([
      'users.id as user_id',
      'users.name as user_name',
      'posts.id as post_id',
      'posts.title as post_title',
    ])
    .orderBy('users.id')
    .execute();

  // 组装成 [{ id, name, posts: [...] }] 的形式
  const usersMap = new Map<number, { id: number; name: string; posts: any[] }>();

  for (const row of result) {
    if (!usersMap.has(row.user_id)) {
      usersMap.set(row.user_id, { id: row.user_id, name: row.user_name, posts: [] });
    }
    if (row.post_id) {
      usersMap.get(row.user_id)!.posts.push({ id: row.post_id, title: row.post_title });
    }
  }

  res.json(Array.from(usersMap.values()));
});
```

`leftJoin` 左连接

## 事务

```ts
/**
 * POST /posts
 * 创建帖子 使用事务同步更新用户帖子数量
 */
router.post('/', async (req, res) => {
  const { title, content, user_id } = req.body;
  await db.transaction().execute(async (tx) => {
    const post = await tx
      .insertInto('posts')
      .values({ title, content, user_id: user_id })
      .executeTakeFirst();
    if (!post) {
      throw new Error('帖子创建失败');
    }
    await tx
      .updateTable('users')
      .set({ post_count: sql`post_count + 1` })
      .where('id', '=', user_id)
      .execute();
  });
  res.json({ message: '帖子创建成功' });
});
```

`db.transaction().execute()`