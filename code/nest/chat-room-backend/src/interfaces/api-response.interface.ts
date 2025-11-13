export interface ApiResponse<T = any> {
  success: boolean;
  code: number; // 通常是 HTTP 状态码或自定义业务码
  data: T | null;
  message: string;
}
