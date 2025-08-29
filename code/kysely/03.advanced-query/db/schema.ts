import { ColumnType, Generated } from 'kysely';

export interface UsersTable {
  id: Generated<number>;
  name: string;
  email: string;
  status: 'active' | 'inactive';
  post_count: number;
  created_at: ColumnType<Date, string | undefined, never>;
}

export interface PostsTable {
  id: Generated<number>;
  user_id: number;
  title: string;
  content: string;
  created_at: ColumnType<Date, string | undefined, never>;
}

export interface Database {
  users: UsersTable;
  posts: PostsTable;
}
