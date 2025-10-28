import { db, pgClient } from './db';
import { users } from './schema';
import { eq } from 'drizzle-orm'; // 引入比较操作符

// 1. 新增用户
async function createUser(name: string, email: string, passwordHash: string) {
  console.log(`尝试创建用户: ${email}`);
  try {
    const result = await db
      .insert(users)
      .values({ name, email, password: passwordHash })
      .returning()
    console.log('用户创建成功:', result);
    return result[0];
  } catch (error) {
    console.error('用户创建失败:', error);
  }
}

// 2. 查询所有用户
async function getAllUsers() {
  console.log('查询所有用户...');
  const result = await db
    .select()
    .from(users)
    .execute()
  console.log('所有用户:', result);
  return result;
}

// 3. 根据 Email 查询单个用户
async function getUserByEmail(email: string) {
  console.log(`根据 Email 查询用户: ${email}`);
  const result = await db
    .select()
    .from(users)
    .where(eq(users.email, email))
    .execute()
  return result[0]; // 只需要返回第一条记录
}

// 运行示例
async function main() {
  // 注意：在实际运行前，你需要配置一个真实的 PostgreSQL 数据库连接字符串
  // 否则代码会因无法连接而报错。

  // 1. 创建用户
  const newUser = await createUser(
    'Alice',
    'alice@example.com',
    'hashed_password_123'
  );
  console.log('新用户:', newUser);

  // 2. 查询所有用户
  const allUsers = await getAllUsers();
  console.log('所有用户:', allUsers);

  // 3. 查询特定用户
  const specificUser = await getUserByEmail('alice@example.com');
  console.log('特定用户:', specificUser);

  await pgClient.end(); // 关闭数据库连接
}

main();
