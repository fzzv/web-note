import { Generated } from 'kysely'

export interface UserTable {
  id: Generated<number>
  name: string
  email: string
  created_at: Generated<Date>
}

export interface Database {
  users: UserTable
}
