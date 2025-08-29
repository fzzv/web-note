import { Router } from 'express';
import { db } from '../db/db';
import { sql } from 'kysely';

const router = Router();

/**
 * GET /posts?page=1&pageSize=10
 * 分页获取帖子
 */
router.get('/', async (req, res) => {
  const page = parseInt((req.query.page as string) || '1', 10);
  const pageSize = parseInt((req.query.pageSize as string) || '10', 10);

  const posts = await db
    .selectFrom('posts')
    .selectAll()
    .orderBy('id', 'desc')
    .limit(pageSize)
    .offset((page - 1) * pageSize)
    .execute();
  res.json(posts);
});

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

export default router;
