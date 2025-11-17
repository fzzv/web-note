export interface User {
  id: number;
  email: string;
  headPic: string;
  nickName: string;
  username: string;
  createTime: Date;
}

export function getUserInfo(): User {
  return JSON.parse(localStorage.getItem('userInfo')!);
}
