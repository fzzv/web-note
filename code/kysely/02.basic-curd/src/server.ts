import express, { Request, Response } from 'express'
import { db } from './db'

const app = express()
app.use(express.json())

// 获取所有用户
app.get('/users', async (_req: Request, res: Response) => {
  const users = await db.selectFrom('users').selectAll().execute()
  res.json(users)
})

// 获取单个用户
app.get('/users/:id', async (req: Request, res: Response) => {
  const user = await db
    .selectFrom('users')
    .selectAll()
    .where('id', '=', Number(req.params.id))
    .executeTakeFirst()
  if (!user) return res.status(404).json({ message: 'User not found' })
  res.json(user)
})

// 新建用户
app.post('/users', async (req: Request, res: Response) => {
  const { name, email } = req.body
  if (!name || !email) {
    return res.status(400).json({ message: 'Name and email required' })
  }

  try {
    const { insertId} = await db
      .insertInto('users')
      .values({ name, email })
      .executeTakeFirstOrThrow()

    if (!insertId) {
      return res.status(500).json({ message: 'Insert failed, no ID returned' })
    }

    const insertedUser = await db
      .selectFrom('users')
      .selectAll()
      .where('id', '=', Number(insertId))
      .executeTakeFirstOrThrow()
      
    res.status(201).json(insertedUser)
  } catch (err) {
    res.status(500).json({ message: 'Insert failed', error: (err as Error).message })
  }
})

// 更新用户
app.put('/users/:id', async (req: Request, res: Response) => {
  const { name, email } = req.body
  const { id } = req.params

  const updated = await db
    .updateTable('users')
    .set({ name, email })
    .where('id', '=', Number(id))
    .executeTakeFirst()

  if (updated.numUpdatedRows === BigInt(0)) {
    return res.status(404).json({ message: 'User not found or not updated' })
  }

  res.json({ message: 'User updated' })
})

// 删除用户
app.delete('/users/:id', async (req: Request, res: Response) => {
  const deleted = await db
    .deleteFrom('users')
    .where('id', '=', Number(req.params.id))
    .executeTakeFirst()

  if (deleted.numDeletedRows === BigInt(0)) {
    return res.status(404).json({ message: 'User not found' })
  }

  res.json({ message: 'User deleted' })
})

// 启动服务
const PORT = 3000
app.listen(PORT, () => {
  console.log(`🚀 Server running at http://localhost:${PORT}`)
})
