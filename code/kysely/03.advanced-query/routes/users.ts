import { Router } from 'express';
import { db } from '../db/db';

const router = Router();

/**
 * GET /users?page=1&pageSize=10
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

/**
 * GET /users/stats/summary
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
    .execute();

  res.json(stats);
});

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

/**
 * 查询文章数最多的前 2 个用户
 */
router.get('/posts/top-2', async (req, res) => {
  const result = await db
    .selectFrom('users')
    .leftJoin('posts', 'users.id', 'posts.user_id')
    .select((eb) => [eb.fn.count<number>('posts.id').as('post_count'), 'users.id', 'users.name'])
    .groupBy('users.id')
    .orderBy('post_count', 'desc')
    .limit(2)
    .execute();
  res.json(result);
});

/**
 * GET /users/:id
 * 获取用户详情
 */
router.get('/:id', async (req, res) => {
  const { id } = req.params;

  const user = await db
    .selectFrom('users')
    .selectAll()
    .where('id', '=', Number(id))
    .executeTakeFirst();

  if (!user) {
    return res.status(404).json({ message: 'User not found' });
  }

  res.json(user);
});

export default router;
